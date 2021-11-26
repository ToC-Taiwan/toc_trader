# TOC TRADER

[![pipeline status](https://gitlab.tocraw.com/root/toc_trader/badges/main/pipeline.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![coverage report](https://gitlab.tocraw.com/root/toc_trader/badges/main/coverage.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://gitlab.tocraw.com/root/toc_trader)
[![Go](https://img.shields.io/badge/Go-1.17.2-blue?logo=go&logoColor=blue)](https://golang.org)
[![OS](https://img.shields.io/badge/OS-Linux-orange?logo=linux&logoColor=orange)](https://www.linux.org/)
[![Container](https://img.shields.io/badge/Container-Docker-blue?logo=docker&logoColor=blue)](https://www.docker.com/)

## Features

[API Docs](http://toc-trader.tocraw.com:6670/swagger/index.html)

## Query Best Condition

```sql
SELECT * FROM simulate_result AS a LEFT JOIN simulate_cond AS b ON a.cond_id=b.id
WHERE positive_days=total_days AND trade_count!=total_days AND total_loss<balance AND forward_balance!=0
order by (balance-total_loss)/trade_count DESC, rsi_high DESC;

SELECT * FROM simulate_result AS a LEFT JOIN simulate_cond AS b ON a.cond_id=b.id
WHERE positive_days=total_days AND trade_count!=total_days AND total_loss<balance AND reverse_balance!=0
order by (balance-total_loss)/trade_count DESC, rsi_low ASC;
```

### Git

```sh
git fetch --prune --prune-tags origin
git check-ignore *
```

### Result and Conditions Sample

```json
```

### Balance Sample

```json
```

### Trade Bot Service

![callvis](./assets/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
