# VWAP Engine

**[VWAP](https://en.wikipedia.org/wiki/Volume-weighted_average_price) Calculator for Go.**

![](https://img.shields.io/github/license/reactivejson/vwap-engine.svg)

Go implementation of VWAP formula, vwap is calculated in real time utilizing
the Coinbase websocket stream "wss://ws-feed.pro.coinbase.com". For each trading pair, the calculated VWAP will be logged to Stdout.
The default trading pairs are BTC-USD, ETH-USD, and ETH-BTC, but you can define your owns via ENV variables

## License

Apache 2.0, see [LICENSE](LICENSE).
