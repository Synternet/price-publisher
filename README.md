#### Price Publisher
[![Latest release](https://img.shields.io/github/v/release/SyntropyNet/price-publisher)](https://github.com/SyntropyNet/price-publisher/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/SyntropyNet/price-publisher/github-ci.yml?label=github-ci)](https://github.com/SyntropyNet/price-publisher/actions/workflows/github-ci.yml)

Retrieves prices from CoinMarketCap and publishes them on interval basis.

## Prerequisites

- Go (v1.21)+

## Flags and Environment Variables

| Environment Variable      | Description                                                                                                                 |
| ------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| NATS_URLS                 | DL NATS broker hosts URLs. Default: `nats://europe-west3-gcp-dl-testnet-brokernode-frankfurt01.syntropynet.com`             |
| NATS_NKEY                 | DL NATS publisher access token.                                                                                             |
| NATS_SUBJECT              | DL publisher stream subject. `syntropy.price` results in: `syntropy.price.tokens`. Default `syntropy.price`.                |
| CMC_IDS                   | Comma separated list of CoinMarketCap tokens ids (e.g.: `825,3408,12220,3794,22861,21420,21686,7226,13678,7431,1027,3717`). |
| CMC_API_KEY               | CoinMarketCap API key.                                                                                                      |
| PUBLISH_INTERVAL_SEC      | Prices publish interval in seconds. Default: `5` seconds.                                                                   |

### CMC_IDS

CMC_IDS can be determined by running
```bash
curl -L 'https://pro-api.coinmarketcap.com/v1/cryptocurrency/map?symbol=USDT,USDC,OSMO,ATOM,TIA,AxlUSDC,STATOM,INJ,PICA,AKT,ETH,WBTC' -H 'X-CMC_PRO_API_KEY: {{API_KEY}}' -H 'Accept: */*' | jq .
```

`{{LIST_OF_SYMBOLS}}` - list of symbols, e.g.: USDT,USDC,OSMO,ATOM,TIA,AxlUSDC,STATOM,INJ,PICA,AKT,ETH,WBTC.

`{{API_KEY}}` - CMC API key. Retrieve from https://pro.coinmarketcap.com/account. Make sure you have appropriate CMC license.

Note: it is possible to map `jq` ids and join them, e.g.: `jq -r '.data | map(.id) | join(",")'`, but `/map?symbol=` can contain multiple entries for same symbol, so cherry-picking is required anyway.

# Makefile

Build from source
```bash
make build
```

Live reload
```bash
make watch
```

Format
```bash
make fmt
```

# Docker

### Build from source

1. Build image.
```
docker build -f ./docker/Dockerfile -t price-publisher .
```

2. Run container with passed environment variables.
```
docker run -it --rm --env-file=.env price-publisher
```
