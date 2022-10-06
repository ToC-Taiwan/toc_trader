# TOC TRADER

[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://gitlab.tocraw.com/root/toc_trader)
[![Go](https://img.shields.io/badge/Go-1.17.3-blue?logo=go&logoColor=blue)](https://golang.org)
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
[
  {
    "trade_day": "2021-11-24T08:00:00+08:00",
    "trade_count": 4,
    "forward": 709,
    "reverse": 0,
    "original_balance": 709,
    "discount": 229,
    "total": 938
  },
  {
    "trade_day": "2021-11-25T08:00:00+08:00",
    "trade_count": 7,
    "forward": 0,
    "reverse": 783,
    "original_balance": 783,
    "discount": 424,
    "total": 1207
  },
  {
    "trade_day": "2021-11-26T08:00:00+08:00",
    "trade_count": 3,
    "forward": -450,
    "reverse": 819,
    "original_balance": 369,
    "discount": 139,
    "total": 508
  },
  {
    "trade_day": "2021-11-29T08:00:00+08:00",
    "trade_count": 3,
    "forward": 12,
    "reverse": -309,
    "original_balance": -297,
    "discount": 123,
    "total": -174
  },
  {
    "trade_day": "2021-11-30T08:00:00+08:00",
    "trade_count": 9,
    "forward": -1896,
    "reverse": 13,
    "original_balance": -1883,
    "discount": 474,
    "total": -1409
  },
  {
    "trade_day": "2021-12-01T08:00:00+08:00",
    "trade_count": 15,
    "forward": -3378,
    "reverse": -4961,
    "original_balance": -8339,
    "discount": 1007,
    "total": -7332
  },
  {
    "trade_day": "2021-12-02T08:00:00+08:00",
    "trade_count": 1,
    "forward": -491,
    "reverse": 0,
    "original_balance": -491,
    "discount": 59,
    "total": -432
  },
  {
    "trade_day": "2021-12-03T08:00:00+08:00",
    "trade_count": 10,
    "forward": -972,
    "reverse": -475,
    "original_balance": -1447,
    "discount": 542,
    "total": -905
  },
  {
    "trade_day": "2021-12-06T08:00:00+08:00",
    "trade_count": 3,
    "forward": -822,
    "reverse": -477,
    "original_balance": -1299,
    "discount": 147,
    "total": -1152
  },
  {
    "trade_day": "2021-12-07T08:00:00+08:00",
    "trade_count": 7,
    "forward": -488,
    "reverse": 716,
    "original_balance": 228,
    "discount": 386,
    "total": 614
  },
  {
    "trade_day": "2021-12-08T08:00:00+08:00",
    "trade_count": 8,
    "forward": -289,
    "reverse": -11,
    "original_balance": -300,
    "discount": 568,
    "total": 268
  },
  {
    "trade_day": "2021-12-09T08:00:00+08:00",
    "trade_count": 5,
    "forward": -1111,
    "reverse": -772,
    "original_balance": -1883,
    "discount": 393,
    "total": -1490
  },
  {
    "trade_day": "2021-12-10T08:00:00+08:00",
    "trade_count": 5,
    "forward": 1841,
    "reverse": -427,
    "original_balance": 1414,
    "discount": 455,
    "total": 1869
  },
  {
    "trade_day": "2021-12-13T08:00:00+08:00",
    "trade_count": 3,
    "forward": 176,
    "reverse": 705,
    "original_balance": 881,
    "discount": 177,
    "total": 1058
  },
  {
    "trade_day": "2021-12-14T08:00:00+08:00",
    "trade_count": 10,
    "forward": -2975,
    "reverse": -945,
    "original_balance": -3920,
    "discount": 725,
    "total": -3195
  },
  {
    "trade_day": "2021-12-15T08:00:00+08:00",
    "trade_count": 12,
    "forward": -1716,
    "reverse": -1983,
    "original_balance": -3699,
    "discount": 821,
    "total": -2878
  },
  {
    "trade_day": "2021-12-16T08:00:00+08:00",
    "trade_count": 9,
    "forward": -1649,
    "reverse": -543,
    "original_balance": -2192,
    "discount": 690,
    "total": -1502
  },
  {
    "trade_day": "2021-12-17T08:00:00+08:00",
    "trade_count": 6,
    "forward": -55,
    "reverse": 0,
    "original_balance": -55,
    "discount": 401,
    "total": 346
  }
]
```

### Trade Bot Service

![callvis](./assets/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
