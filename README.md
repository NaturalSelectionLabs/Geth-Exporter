Geth-Exporter
========================

A [prometheus](https://prometheus.io/) exporter which scrapes geth json rpc data.

## Example Usage

```console
## SETUP

$ go run main.go --config config.yaml &


## TEST with 'balance' module

$ curl "http://localhost:8000/probe?module=balance&target=<RPC_ENDPOINT>&address=<WALLET_ADDR>&block=<BLOCK>"
# HELP balance_wei 
# TYPE balance_wei gauge
balance_wei 3.7689696e+16

## TEST with 'gas' module

$ curl "http://localhost:8000/probe?module=gas&target=<RPC_ENDPOINT>"
# HELP geth_wei 
# TYPE geth_wei gauge
geth_wei 1e+09


## TEST through prometheus:

$ docker run --rm -it -p 9090:9090 -v $PWD/examples/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```
Then head over to http://localhost:9090/graph?g0.range_input=1h&g0.expr=example_value_active&g0.tab=1 or http://localhost:9090/targets to check the scraped metrics or the targets.

## Docker

```console
$ docker run -v $PWD/config.yaml:/config.yaml -p 8000:8000 rss3/geth-exporter --config=/config.yaml
```
