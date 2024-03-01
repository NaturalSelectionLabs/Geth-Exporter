package exporter

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/geth-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

func CreateMetrics(m config.Module) (GethMetric, error) {
	var variableLabels, variableLabelsValues []string
	for k, v := range m.Labels {
		variableLabels = append(variableLabels, k)
		variableLabelsValues = append(variableLabelsValues, v)
	}

	switch m.Type {
	case config.BoolResult, config.HexResult:
		return GethMetric{
			Desc: prometheus.NewDesc(
				m.Name,
				"",
				variableLabels,
				nil,
			),
			Type:           m.Type,
			LabelValues:    variableLabelsValues,
			ValueType:      prometheus.GaugeValue,
			EpochTimestamp: "",
		}, nil
	default:
		return GethMetric{}, fmt.Errorf("Unknown metric type '%s'", m.Type)
	}

}

type GethFetcher struct {
	Client *ethclient.Client
	Method string
	Params []any
}

func NewGethFetcher(endpoint string, module config.Module, params ...any) *GethFetcher {
	client, _ := ethclient.Dial(endpoint)
	return &GethFetcher{
		Client: client,
		Method: module.Method,
		Params: params,
	}
}

func (g *GethFetcher) FetchData(ctx context.Context) (interface{}, error) {
	var result interface{}

	err := g.Client.Client().CallContext(ctx, &result, g.Method, g.Params...)
	return result, err
}

func BoolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
