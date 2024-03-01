package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/geth-exporter/config"
	"github.com/naturalselectionlabs/geth-exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

func probeHandler(c echo.Context, cfg config.Config) error {

	module := c.QueryParam("module")
	if module == "" {
		module = "default"
	}

	if _, ok := cfg.Modules[module]; !ok {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Unknown module %s", module))
	}

	registry := prometheus.NewPedanticRegistry()

	metric, err := exporter.CreateMetrics(cfg.Modules[module])
	if err != nil {
		zap.L().Error(
			"Failed to create metrics from module",
			zap.String("module", module),
			zap.Error(err))
	}

	gethMetricCollector := exporter.GethMetricsCollector{
		Metrics: metric,
		Logger:  zap.L(),
	}

	target := c.QueryParam("target")
	if target == "" {
		return c.String(http.StatusBadRequest, "Target parameter is missing")
	}

	paramKeys := cfg.Modules[module].Params

	var params []any

	for _, key := range paramKeys {
		value := c.QueryParam(key)
		switch key {
		case "address":
			params = append(params, common.HexToAddress(value))
		default:
			params = append(params, value)
		}
	}

	fetcher := exporter.NewGethFetcher(target, cfg.Modules[module], params...)

	result, err := fetcher.FetchData(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("Failed to fetch with error: %s", err))
	}

	gethMetricCollector.Result = result
	registry.MustRegister(gethMetricCollector)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(c.Response(), c.Request())
	return nil
}
