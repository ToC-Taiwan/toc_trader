# TOC TRADER

[![pipeline status](https://gitlab.tocraw.com/root/toc_trader/badges/main/pipeline.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![coverage report](https://gitlab.tocraw.com/root/toc_trader/badges/main/coverage.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://gitlab.tocraw.com/root/toc_trader)
[![Go](https://img.shields.io/badge/Go-1.17.2-blue?logo=go&logoColor=blue)](https://golang.org)
[![OS](https://img.shields.io/badge/OS-Linux-orange?logo=linux&logoColor=orange)](https://www.linux.org/)
[![Container](https://img.shields.io/badge/Container-Docker-blue?logo=docker&logoColor=blue)](https://www.docker.com/)

## Features

### Initialize

- 確認 `global.db` 是否存在，不存在則依照 `default key` 創建，有些有預設值
- 初始化 database(postgresql)
- 初始化 task
- 預設開啟 EnableBuy, EnableSell, MeanTimeTradeStockNum = ?
- 2021 休市日匯入資料庫，並會自動跳過
- 找出程式啟用當下的前兩個交易日及當下的交易日

### Servers

- 優先開啟 API server(default: 6670)，確認 port is open 才繼續以下動作
- [API Docs](http://toc-trader.tocraw.com:6670/swagger/index.html)

### Flows

- TradeID, TradePass 須於初次啟動初始化，不存在僅開放更新此值的 API
- Sysparm API 可更新 reset 於下次啟動時，初始化資料庫，開啟完畢後，reset 會歸零
- EnableBuy, EnableSell, MeanTimeTradeStockNum 可於程式進行中修改
- Docker Container 模式下，System API 可重啟服務
- 匯入所有股票
- 更新前一交易日目標股價及成交量
- 透過 sysparm 中條件篩選出股票（目前條件不可更改）
- 擷取當前日期的前一交易日所有搓合交易，並與訂閱做相同分析
- 每日 PM 3:00 之後才開放擷取當日 tick
- 啟動時，會檢查委託，如有成交買單，會扣除對應額度
- 訂閱目標的 Streamtick, bid-ask
- 每 1.5 秒確認委託狀態
- AM 4:00 會清除所有事件(委託)

### Git

```sh
git fetch --prune --prune-tags origin
git check-ignore *
```

### Simulation

- Both 1 Day

```log
WARN[2021/10/21 21:47:22] 2021-10-21 Forward Balance: 939, Stock: 1308, Name: 亞聚, Total Time: 5672, 42.20, 43.25
WARN[2021/10/21 21:47:22] 2021-10-21 Forward Balance: 123, Stock: 1727, Name: 中華化, Total Time: 1719, 48.75, 49.00
WARN[2021/10/21 21:47:22] 2021-10-21 Forward Balance: -739, Stock: 1708, Name: 東鹼, Total Time: 4881, 35.20, 34.55
WARN[2021/10/21 21:47:22] 2021-10-21 Forward Balance: 311, Stock: 2371, Name: 大同, Total Time: 9513, 33.85, 34.25
WARN[2021/10/21 21:47:22] 2021-10-21 Forward Balance: 672, Stock: 5351, Name: 鈺創, Total Time: 441, 49.15, 49.95
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 742, Stock: 2399, Name: 映泰, Total Time: 12343, 21.70, 22.50
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 659, Stock: 2457, Name: 飛宏, Total Time: 8670, 34.70, 35.45
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 232, Stock: 2406, Name: 國碩, Total Time: 3426, 25.85, 26.15
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 231, Stock: 1711, Name: 永光, Total Time: 5253, 26.05, 26.35
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: -245, Stock: 2409, Name: 友達, Total Time: 3132, 17.15, 16.95
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: -21, Stock: 2027, Name: 大成鋼, Total Time: 4852, 46.40, 46.50
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 280, Stock: 2023, Name: 燁輝, Total Time: 2881, 26.45, 26.80
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 111, Stock: 2002, Name: 中鋼, Total Time: 7098, 33.95, 34.15
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 106, Stock: 2337, Name: 旺宏, Total Time: 8763, 35.90, 36.10
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 1027, Stock: 2498, Name: 宏達電, Total Time: 541, 46.80, 47.95
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 907, Stock: 2401, Name: 凌陽, Total Time: 11593, 35.65, 36.65
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 135, Stock: 2344, Name: 華邦電, Total Time: 6583, 25.10, 25.30
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: 61, Stock: 2340, Name: 光磊, Total Time: 4447, 34.00, 34.15
WARN[2021/10/21 21:47:22] 2021-10-21 Reverse Balance: -190, Stock: 3481, Name: 群創, Total Time: 2049, 15.55, 15.40
WARN[2021/10/21 21:47:22] Total Balance: 5341, TradeCount: 19
WARN[2021/10/21 21:47:22] Cond: {Model:{ID:5409 CreatedAt:2021-10-21 21:44:06.483372 +0800 CST UpdatedAt:2021-10-21 21:44:06.483372 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} HistoryCloseCount:2500 OutInRatio:70 ReverseOutInRatio:5 CloseDiff:0 CloseChangeRatioLow:-1 CloseChangeRatioHigh:8 OpenChangeRatio:4 RsiHigh:50.1 RsiLow:50 ReverseRsiHigh:50.1 ReverseRsiLow:50 TicksPeriodThreshold:5 TicksPeriodLimit:6.5 TicksPeriodCount:1 VolumePerSecond:5}

WARN[2021/10/26 01:36:39] 2021-10-25 Reverse Balance: 1065, Stock: 5351, Name: 鈺創, Total Time: 2072, 51.20, 52.40
WARN[2021/10/26 01:36:39] 2021-10-25 Reverse Balance: 105, Stock: 2340, Name: 光磊, Total Time: 6161, 36.80, 37.00
WARN[2021/10/26 01:36:39] 2021-10-25 Forward: 0, Reverse: 1170
WARN[2021/10/26 01:36:39] 2021-10-22 Reverse Balance: -519, Stock: 3508, Name: 位速, Total Time: 3917, 26.85, 26.40
WARN[2021/10/26 01:36:39] 2021-10-22 Reverse Balance: 882, Stock: 1710, Name: 東聯, Total Time: 7372, 25.70, 26.65
WARN[2021/10/26 01:36:39] 2021-10-22 Forward: 0, Reverse: 363
WARN[2021/10/26 01:36:39] 2021-10-21 Reverse Balance: 505, Stock: 2401, Name: 凌陽, Total Time: 6275, 36.10, 36.70
WARN[2021/10/26 01:36:39] 2021-10-21 Reverse Balance: 643, Stock: 2399, Name: 映泰, Total Time: 7644, 21.60, 22.30
WARN[2021/10/26 01:36:39] 2021-10-21 Reverse Balance: -677, Stock: 1727, Name: 中華化, Total Time: 2418, 49.30, 48.75
WARN[2021/10/26 01:36:39] 2021-10-21 Reverse Balance: 976, Stock: 2498, Name: 宏達電, Total Time: 716, 47.15, 48.25
WARN[2021/10/26 01:36:39] 2021-10-21 Forward: 0, Reverse: 1447
WARN[2021/10/26 01:36:39] Total Balance: 2980, TradeCount: 23, PositiveCount: 3
WARN[2021/10/26 01:36:39] Cond: {Model:{ID:41638 CreatedAt:2021-10-26 01:23:54.539557 +0800 CST UpdatedAt:2021-10-26 01:23:54.539557 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} HistoryCloseCount:2000 OutInRatio:80 ReverseOutInRatio:6 CloseDiff:0 CloseChangeRatioLow:0 CloseChangeRatioHigh:6 OpenChangeRatio:3 RsiHigh:50.1 RsiLow:50 ReverseRsiHigh:50.1 ReverseRsiLow:50 TicksPeriodThreshold:9 TicksPeriodLimit:11.700000000000001 TicksPeriodCount:2 VolumePerSecond:6}
```

### Trade Bot Service

![callvis](./assets/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)

```log
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -195, Stock: 2337, Name: 旺宏, Sell first at 2021-10-25 11:15:21, Buy later at 2021-10-25 11:15:58
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -189, Stock: 2014, Name: 中鴻, Sell first at 2021-10-25 10:41:22, Buy later at 2021-10-25 10:42:32
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: 102, Stock: 2401, Name: 凌陽, Sell first at 2021-10-25 09:33:20, Buy later at 2021-10-25 09:37:31
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -169, Stock: 1710, Name: 東聯, Sell first at 2021-10-25 10:21:07, Buy later at 2021-10-25 10:22:48
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -166, Stock: 1711, Name: 永光, Sell first at 2021-10-25 09:44:10, Buy later at 2021-10-25 09:44:12
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: 491, Stock: 1308, Name: 亞聚, Sell first at 2021-10-25 11:55:38, Buy later at 2021-10-25 13:10:06
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: 306, Stock: 2340, Name: 光磊, Sell first at 2021-10-25 09:29:21, Buy later at 2021-10-25 09:44:22
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -16, Stock: 2344, Name: 華邦電, Sell first at 2021-10-25 12:20:53, Buy later at 2021-10-25 13:10:00
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -207, Stock: 2328, Name: 廣宇, Sell first at 2021-10-25 09:51:44, Buy later at 2021-10-25 09:51:56
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -228, Stock: 2388, Name: 威盛, Sell first at 2021-10-25 09:57:35, Buy later at 2021-10-25 09:57:52
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -218, Stock: 1727, Name: 中華化, Sell first at 2021-10-25 09:51:54, Buy later at 2021-10-25 09:52:11
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: -165, Stock: 2023, Name: 燁輝, Sell first at 2021-10-25 11:45:01, Buy later at 2021-10-25 11:45:03
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: 143, Stock: 3062, Name: 建漢, Sell first at 2021-10-25 09:06:23, Buy later at 2021-10-25 09:06:23
WARN[2021/10/27 03:29:39] 2021-10-25 Reverse Balance: 164, Stock: 5351, Name: 鈺創, Sell first at 2021-10-25 09:25:02, Buy later at 2021-10-25 09:25:37
WARN[2021/10/27 03:29:39] 2021-10-25 Forward: 0, Reverse: -347
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: -166, Stock: 2426, Name: 鼎元, Sell first at 2021-10-26 11:41:54, Buy later at 2021-10-26 11:43:57
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: 81, Stock: 2344, Name: 華邦電, Sell first at 2021-10-26 09:35:40, Buy later at 2021-10-26 09:42:19
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: 201, Stock: 2340, Name: 光磊, Sell first at 2021-10-26 09:14:29, Buy later at 2021-10-26 09:16:29
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: -199, Stock: 2337, Name: 旺宏, Sell first at 2021-10-26 10:28:30, Buy later at 2021-10-26 10:29:08
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: 501, Stock: 2401, Name: 凌陽, Sell first at 2021-10-26 09:56:03, Buy later at 2021-10-26 10:11:55
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: -27, Stock: 1101, Name: 台泥, Sell first at 2021-10-26 11:00:44, Buy later at 2021-10-26 11:22:29
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: 35, Stock: 2104, Name: 國際中橡, Sell first at 2021-10-26 11:09:04, Buy later at 2021-10-26 11:21:05
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: 239, Stock: 3062, Name: 建漢, Sell first at 2021-10-26 09:42:43, Buy later at 2021-10-26 09:43:38
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: -174, Stock: 2605, Name: 新興, Sell first at 2021-10-26 12:12:22, Buy later at 2021-10-26 12:44:52
WARN[2021/10/27 03:29:39] 2021-10-26 Reverse Balance: 331, Stock: 3323, Name: 加百裕, Sell first at 2021-10-26 09:11:45, Buy later at 2021-10-26 09:15:38
WARN[2021/10/27 03:29:39] 2021-10-26 Forward: 0, Reverse: 822
WARN[2021/10/27 03:29:39] Total Balance: 475, TradeCount: 26, PositiveCount: 1
WARN[2021/10/27 03:29:39] Cond: {Model:{ID:43008 CreatedAt:2021-10-27 03:25:48.810576 +0800 CST UpdatedAt:2021-10-27 03:25:48.810576 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} TrimHistoryCloseCount:true HistoryCloseCount:1500 OutInRatio:95 ReverseOutInRatio:4 CloseDiff:0 CloseChangeRatioLow:-3 CloseChangeRatioHigh:7 OpenChangeRatio:7 RsiHigh:50.4 RsiLow:50 ReverseRsiHigh:50.4 ReverseRsiLow:50 TicksPeriodThreshold:6 TicksPeriodLimit:7.800000000000001 TicksPeriodCount:1 VolumePerSecond:8}
```
