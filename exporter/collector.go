package exporter

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/naturalselectionlabs/geth-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type GethMetricsCollector struct {
	Metrics GethMetric
	Result  any
	Logger  *zap.Logger
}

func (g GethMetricsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- g.Metrics.Desc
}

func (g GethMetricsCollector) Collect(metrics chan<- prometheus.Metric) {
	switch g.Metrics.Type {
	case config.BoolResult:
		value, ok := g.Result.(bool)
		if ok {
			metrics <- prometheus.MustNewConstMetric(
				g.Metrics.Desc,
				g.Metrics.ValueType,
				float64(BoolToInt(value)),
				g.Metrics.LabelValues...,
			)
		}
	case config.HexResult:
		value, ok := g.Result.(string)
		if ok {
			i, err := hexutil.DecodeBig(value)
			if err != nil {
				zap.L().Error("Error convert", zap.Error(err))
				return
			}
			f, _ := i.Float64()
			metrics <- prometheus.MustNewConstMetric(
				g.Metrics.Desc,
				g.Metrics.ValueType,
				f,
				g.Metrics.LabelValues...,
			)

		}
	}
}

type GethMetric struct {
	Desc           *prometheus.Desc
	Type           config.ResultType
	LabelValues    []string
	ValueType      prometheus.ValueType
	EpochTimestamp string
}
