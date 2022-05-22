# VWAP Engine

**[VWAP](https://en.wikipedia.org/wiki/Volume-weighted_average_price) Calculator for Go.**

![Maintainer](https://img.shields.io/badge/maintainer-Mohamed-Aly-blue)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/reactivejson/vwap-engine.svg)](https://github.com/reactivejson/vwap-engine)
[![Go Reference](https://pkg.go.dev/badge/github.com/reactivejson/vwap-engine)](https://pkg.go.dev/badge/github.com/reactivejson/vwap-engine)
![Build Status](https://github.com/reactivejson/vwap-engine/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/reactivejson/vwap-engine)](https://goreportcard.com/report/github.com/reactivejson/vwap-engine)
![](https://img.shields.io/github/license/reactivejson/vwap-engine.svg)

Go implementation of Volume-weighted average price. vwap is calculated in real time utilizing
the Coinbase websocket stream "wss://ws-feed.pro.coinbase.com". For each trading pair, the calculated VWAP will be logged to Stdout.
The default trading pairs are BTC-USD, ETH-USD, and ETH-BTC, but you can define your owns via ENV variables


## Project layout

This layout is following pattern:

```text
vwap-engine
└───
    ├── .env // optional
    ├── .github
    │   └── workflows
    │     └── go.yml
    ├── api
    │   └── models
    │     └── coinbase.go
    ├── cmd
    │   └── main.go
    ├── internal
    │   └── app
    │     └── app.go
    │     └── setup.go
    │     └── context.go
    │   └── tunnel
    │     └── receiver.go
    │     └── tunnel.go
    │   └── storage
    │     └── vwap.go
    │     └── vwap-queue.go
    ├── build
    │   └── Dockerfile
    ├── helm
    │   └── <helm chart files>
    ├── Makefile
    ├── Jenkinsfile
    ├── README.md
    ├── VERSION
    └── <source packages>
```

## Architecture
This application main components comprises:

### Tunnel interface
A coinbase websocket stream client to receive data from coinbase websocket server.
1) It subscribes(Tunnel.Subscribe) to the coinbase channel's websocket using trading pairs (productIDs).
2) Read (Tunnel.Read) real-time data points from the coinbase and passes them to the receiver channel.

### VWAP interface:
- represents a queue of DataPoints and their VWAPs.
- DataPoints  is fast circular fifo data structure (aka., queue) with a specific limit.
- For every new coinbase entry, it pushes an item onto the queue and calculates the new VWAP.
- When Limit is reached, will delete  the first one.
Every time a new data point is added to the queue and saved for each trading pair, the VWAP computation is updated accordingly.
For performance, and to avoid exponential complexity, the computation is cached for VWAP, CumulativeQuantity,
and CumulativePriceQuantity for existing data points and updated with new entries.


### Main
The core entry point into the app. will setup the config,
run the App context, It is resilient tolerant. It will gracefully shutdown and can receive an interrupt signal and safely close the connexio.

### App
Setup the config, run the App context, subscribe to the ws, and initiate the vwap storage and calculation for the trading pairs. It is resilient tolerant.

### Helm & K8S
Helm charts to deploy this micro-service in a Kubernetes platform
We generate the container image and reference it in a Helm chart

## Setup
The app is configurable via the ENV variables or Helm values for cloud-native deployment
Config parameters:
- TRADING_PAIRS: a list of coinbase product IDS. Example: BTC-USD,ETH-USD,ETH-BTC
- WEBSOCKET_URL: coinbase websocket server. Example: wss://ws-feed.pro.coinbase.com
- WINDOW_SIZE: Data points sliding window for VWAP computation.



### Getting started
Vwap engine is available in github
[Vwap-engine](https://github.com/reactivejson/vwap-engine)

```shell
go get github.com/reactivejson/storage-engine
```

#### Run
```shell
go run cmd/main.go
```

#### Build
```shell
make build
```
#### Testing
```shell
make test
```

### Build docker image:

```bash
make docker-build
```
This will build this application docker image so-called vwap-engine

### Deploy with Helm chart in a Kubernetes environment
```bash
 helm upgrade --namespace neo --install vwap-engine chart/vwap-engine -f your-custom-values.yml
```
## Test coverage

Test coverage is checked as a part of test execution with the gotestsum tool.

Test coverage is checked for unit tests and integration tests.

Coverage report files are available and stored as `*coverage.txt` and are also imported in the SonarQube for easier browsing.


## golangci-lint

In the effort of reducing errors and improving the overall quality of code, golangci-lint is run as a part of the pipeline. Linting is run for the services and packages that have changes since the previous green build (in master) or previous commit (in local or review).

Any issues found by golangci-lint for the changed code will lead to a failed build.

golangci-lint rules are configured in `.golangci.yml`.

## Sonar scan

A Sonar scan is run as a part of the pipeline.

### Requirements

- Go 1.15 or newer [https://golang.org/doc/install](https://golang.org/doc/install)
- Docker 18.09.6 or newer
- docker-compose 1.24.0 on newer

### Variable names
Commonly used one letter variable names:

- i for index
- r for reader
- w for writer
- c for client
- vwap for Volume Weighted Average Price


## License

Apache 2.0, see [LICENSE](LICENSE).
