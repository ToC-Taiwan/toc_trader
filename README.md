# TOC TRADER

[![pipeline status](https://gitlab.tocraw.com/root/toc_trader/badges/main/pipeline.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![coverage report](https://gitlab.tocraw.com/root/toc_trader/badges/main/coverage.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://gitlab.tocraw.com/root/toc_trader)
[![Go](https://img.shields.io/badge/Go-1.17.2-blue?logo=go&logoColor=blue)](https://golang.org)
[![OS](https://img.shields.io/badge/OS-Linux-orange?logo=linux&logoColor=orange)](https://www.linux.org/)
[![Container](https://img.shields.io/badge/Container-Docker-blue?logo=docker&logoColor=blue)](https://www.docker.com/)

## Features

- [API Docs](http://toc-trader.tocraw.com:6670/swagger/index.html)

## Query Best Condition

```sql
select * from simulate_result as a LEFT JOIN simulate_cond as b ON a.cond_id=b.id WHERE a.positive_days=a.total_days order by balance/trade_count DESC, rsi_high-rsi_low DESC, rsi_low ASC;
```

### Git

```sh
git fetch --prune --prune-tags origin
git check-ignore *
```

### Trade Bot Service

![callvis](./assets/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
