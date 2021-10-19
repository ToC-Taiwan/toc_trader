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

- Forward

```log
WARN[2021/10/18 02:41:25] 2021-10-15 Forward Balance: -538, Stock: 2609, Name: 陽明, Total Time: 1168, 92.60, 92.30
WARN[2021/10/18 02:41:25] 2021-10-15 Forward Balance: 448, Stock: 2303, Name: 聯電, Total Time: 761, 58.20, 58.80
WARN[2021/10/18 02:41:25] 2021-10-15 Forward Balance: -743, Stock: 2603, Name: 長榮, Total Time: 808, 94.70, 94.20
WARN[2021/10/18 02:41:25] 2021-10-14 Forward Balance: -1739, Stock: 2603, Name: 長榮, Total Time: 230, 93.60, 92.10
WARN[2021/10/18 02:41:25] 2021-10-14 Forward Balance: 1768, Stock: 2609, Name: 陽明, Total Time: 1077, 88.00, 90.00
WARN[2021/10/18 02:41:25] 2021-10-14 Forward Balance: 353, Stock: 2303, Name: 聯電, Total Time: 1533, 56.70, 57.20
WARN[2021/10/18 02:41:25] 2021-10-13 Forward Balance: -450, Stock: 2303, Name: 聯電, Total Time: 369, 58.10, 57.80
WARN[2021/10/18 02:41:25] 2021-10-13 Forward Balance: -936, Stock: 2609, Name: 陽明, Total Time: 630, 92.00, 91.30
WARN[2021/10/18 02:41:25] 2021-10-13 Forward Balance: -1068, Stock: 2606, Name: 裕民, Total Time: 2684, 65.70, 64.80
WARN[2021/10/18 02:41:25] 2021-10-13 Forward Balance: -146, Stock: 2603, Name: 長榮, Total Time: 368, 95.20, 95.30
WARN[2021/10/18 02:41:25] Total Balance: -3051, TradeCount: 10
WARN[2021/10/18 02:41:25] Cond: {Model:{ID:43376 CreatedAt:2021-10-18 02:38:52.820999 +0800 CST UpdatedAt:2021-10-18 02:38:52.820999 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} HistoryCloseCount:2500 OutInRatio:55 ReverseOutInRatio:5 CloseDiff:0 CloseChangeRatioLow:-1 CloseChangeRatioHigh:8 OpenChangeRatio:4 RsiHigh:50.2 RsiLow:50 ReverseRsiHigh:50.2 ReverseRsiLow:50 TicksPeriodThreshold:7 TicksPeriodLimit:9.1 TicksPeriodCount:3 VolumePerSecond:3}
```

- Reverse

```log
WARN[2021/10/18 02:34:42] 2021-10-15 Reverse Balance: 654, Stock: 2603, Name: 長榮, Total Time: 183, 94.60, 95.50
WARN[2021/10/18 02:34:42] 2021-10-15 Reverse Balance: 359, Stock: 2609, Name: 陽明, Total Time: 298, 92.70, 93.30
WARN[2021/10/18 02:34:42] 2021-10-15 Reverse Balance: -51, Stock: 2303, Name: 聯電, Total Time: 448, 58.40, 58.50
WARN[2021/10/18 02:34:42] 2021-10-15 Reverse Balance: 543, Stock: 2606, Name: 裕民, Total Time: 2936, 60.10, 60.80
WARN[2021/10/18 02:34:42] 2021-10-14 Reverse Balance: 353, Stock: 2303, Name: 聯電, Total Time: 503, 56.80, 57.30
WARN[2021/10/18 02:34:42] 2021-10-14 Reverse Balance: 1155, Stock: 2603, Name: 長榮, Total Time: 260, 93.40, 94.80
WARN[2021/10/18 02:34:42] 2021-10-14 Reverse Balance: -1329, Stock: 2609, Name: 陽明, Total Time: 628, 89.50, 88.40
WARN[2021/10/18 02:34:42] 2021-10-14 Reverse Balance: -654, Stock: 2606, Name: 裕民, Total Time: 858, 59.90, 59.40
WARN[2021/10/18 02:34:42] 2021-10-14 Reverse Balance: 1565, Stock: 8121, Name: 越峰, Total Time: 988, 50.90, 52.60
WARN[2021/10/18 02:34:42] 2021-10-13 Reverse Balance: 955, Stock: 2603, Name: 長榮, Total Time: 239, 93.90, 95.10
WARN[2021/10/18 02:34:42] 2021-10-13 Reverse Balance: 349, Stock: 2303, Name: 聯電, Total Time: 484, 58.00, 58.50
WARN[2021/10/18 02:34:42] 2021-10-13 Reverse Balance: -1236, Stock: 2609, Name: 陽明, Total Time: 168, 92.30, 91.30
WARN[2021/10/18 02:34:42] 2021-10-13 Reverse Balance: 1224, Stock: 2606, Name: 裕民, Total Time: 601, 66.80, 68.20
WARN[2021/10/18 02:34:42] 2021-10-12 Reverse Balance: 2058, Stock: 2609, Name: 陽明, Total Time: 557, 92.20, 94.50
WARN[2021/10/18 02:34:42] 2021-10-12 Reverse Balance: -242, Stock: 3711, Name: 日月光投控, Total Time: 2069, 93.90, 93.90
WARN[2021/10/18 02:34:42] 2021-10-12 Reverse Balance: 2224, Stock: 2606, Name: 裕民, Total Time: 4259, 66.30, 68.70
WARN[2021/10/18 02:34:42] 2021-10-12 Reverse Balance: -352, Stock: 2303, Name: 聯電, Total Time: 644, 59.40, 59.20
WARN[2021/10/18 02:34:42] Total Balance: 7575, TradeCount: 17
WARN[2021/10/18 02:34:42] Cond: {Model:{ID:37659 CreatedAt:2021-10-18 02:29:41.251255 +0800 CST UpdatedAt:2021-10-18 02:29:41.251255 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} HistoryCloseCount:1500 OutInRatio:55 ReverseOutInRatio:5 CloseDiff:0 CloseChangeRatioLow:-1 CloseChangeRatioHigh:8 OpenChangeRatio:4 RsiHigh:50 RsiLow:50 ReverseRsiHigh:50 ReverseRsiLow:50 TicksPeriodThreshold:1 TicksPeriodLimit:1.3 TicksPeriodCount:2 VolumePerSecond:6}
```

- Both

```log
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -539, Stock: 2371, Name: 大同, Total Time: 4606, 34.95, 34.50
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -286, Stock: 2002, Name: 中鋼, Total Time: 4523, 33.65, 33.45
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: 119, Stock: 1718, Name: 中纖, Total Time: 6507, 11.30, 11.45
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -369, Stock: 1711, Name: 永光, Total Time: 1702, 26.95, 26.65
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -54, Stock: 2328, Name: 廣宇, Total Time: 2716, 40.35, 40.40
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: 636, Stock: 1710, Name: 東聯, Total Time: 2482, 24.45, 25.15
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -164, Stock: 2027, Name: 大成鋼, Total Time: 711, 44.55, 44.50
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: 1548, Stock: 2498, Name: 宏達電, Total Time: 5723, 38.10, 39.75
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -206, Stock: 2399, Name: 映泰, Total Time: 8169, 21.60, 21.45
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -208, Stock: 2201, Name: 裕隆, Total Time: 1679, 42.05, 41.95
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -97, Stock: 1455, Name: 集盛, Total Time: 3132, 18.25, 18.20
WARN[2021/10/20 01:44:33] 2021-10-19 Forward Balance: -92, Stock: 3481, Name: 群創, Total Time: 4700, 16.05, 16.00
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: 37, Stock: 2344, Name: 華邦電, Total Time: 5725, 24.50, 24.60
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: 51, Stock: 2618, Name: 長榮航, Total Time: 4621, 18.90, 19.00
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: -74, Stock: 1802, Name: 台玻, Total Time: 508, 28.45, 28.45
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: 686, Stock: 6477, Name: 安集, Total Time: 11803, 43.45, 44.25
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: -608, Stock: 8121, Name: 越峰, Total Time: 6297, 42.10, 41.60
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: 628, Stock: 1727, Name: 中華化, Total Time: 1083, 46.55, 47.30
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: -310, Stock: 1440, Name: 南紡, Total Time: 7363, 23.45, 23.20
WARN[2021/10/20 01:44:33] 2021-10-19 Reverse Balance: -36, Stock: 6116, Name: 彩晶, Total Time: 7469, 13.85, 13.85
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: 30, Stock: 1711, Name: 永光, Total Time: 1161, 26.65, 26.75
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -767, Stock: 2027, Name: 大成鋼, Total Time: 12169, 45.70, 45.05
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -213, Stock: 1308, Name: 亞聚, Total Time: 7672, 43.75, 43.65
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -203, Stock: 1305, Name: 華夏, Total Time: 9574, 39.80, 39.70
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -36, Stock: 2002, Name: 中鋼, Total Time: 8911, 33.50, 33.55
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -397, Stock: 1304, Name: 台聚, Total Time: 12591, 37.60, 37.30
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -145, Stock: 2409, Name: 友達, Total Time: 4962, 17.25, 17.15
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: -216, Stock: 1605, Name: 華新, Total Time: 3732, 25.75, 25.60
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: 308, Stock: 2371, Name: 大同, Total Time: 3271, 35.45, 35.85
WARN[2021/10/20 01:44:33] 2021-10-18 Forward Balance: 786, Stock: 6477, Name: 安集, Total Time: 3218, 43.45, 44.35
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: -606, Stock: 1440, Name: 南紡, Total Time: 4073, 22.45, 21.90
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 75, Stock: 2605, Name: 新興, Total Time: 1058, 29.15, 29.30
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 732, Stock: 2358, Name: 廷鑫, Total Time: 7220, 25.35, 26.15
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 235, Stock: 5608, Name: 四維航, Total Time: 6097, 44.50, 44.85
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 4231, Stock: 8121, Name: 越峰, Total Time: 7730, 42.80, 47.15
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 919, Stock: 6217, Name: 中探針, Total Time: 2868, 49.75, 50.80
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 1140, Stock: 2328, Name: 廣宇, Total Time: 4275, 41.45, 42.70
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: -911, Stock: 1727, Name: 中華化, Total Time: 12213, 43.55, 42.75
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: -363, Stock: 3701, Name: 大眾控, Total Time: 2542, 43.80, 43.55
WARN[2021/10/20 01:44:33] 2021-10-18 Reverse Balance: 50, Stock: 4503, Name: 金雨, Total Time: 2433, 38.65, 38.80
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 741, Stock: 2328, Name: 廣宇, Total Time: 2808, 41.65, 42.50
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 4, Stock: 2409, Name: 友達, Total Time: 3589, 17.35, 17.40
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 47, Stock: 1305, Name: 華夏, Total Time: 9585, 39.30, 39.45
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: -159, Stock: 4711, Name: 永純, Total Time: 8583, 22.70, 22.60
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: -192, Stock: 3481, Name: 群創, Total Time: 7883, 16.30, 16.15
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: -64, Stock: 5608, Name: 四維航, Total Time: 825, 44.15, 44.20
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 84, Stock: 1605, Name: 華新, Total Time: 11595, 25.20, 25.35
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 584, Stock: 1711, Name: 永光, Total Time: 1057, 25.10, 25.75
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: -16, Stock: 2344, Name: 華邦電, Total Time: 9769, 25.70, 25.75
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: -290, Stock: 2371, Name: 大同, Total Time: 4225, 34.90, 34.70
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 2880, Stock: 3701, Name: 大眾控, Total Time: 9014, 43.90, 46.90
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 1400, Stock: 1727, Name: 中華化, Total Time: 967, 37.20, 38.70
WARN[2021/10/20 01:44:33] 2021-10-15 Forward Balance: 115, Stock: 2002, Name: 中鋼, Total Time: 9237, 33.05, 33.25
WARN[2021/10/20 01:44:33] 2021-10-15 Reverse Balance: -1019, Stock: 8121, Name: 越峰, Total Time: 2404, 46.70, 45.80
WARN[2021/10/20 01:44:33] 2021-10-15 Reverse Balance: -1010, Stock: 1308, Name: 亞聚, Total Time: 9537, 43.35, 42.45
WARN[2021/10/20 01:44:33] 2021-10-15 Reverse Balance: -95, Stock: 1304, Name: 台聚, Total Time: 4097, 36.70, 36.70
WARN[2021/10/20 01:44:33] 2021-10-15 Reverse Balance: -549, Stock: 4503, Name: 金雨, Total Time: 5165, 39.00, 38.55
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: -245, Stock: 2610, Name: 華航, Total Time: 9021, 17.10, 16.90
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: 630, Stock: 1711, Name: 永光, Total Time: 926, 26.10, 26.80
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: -734, Stock: 2002, Name: 中鋼, Total Time: 2490, 33.30, 32.65
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: -64, Stock: 2027, Name: 大成鋼, Total Time: 864, 44.60, 44.65
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: 151, Stock: 1727, Name: 中華化, Total Time: 741, 38.35, 38.60
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: -145, Stock: 2409, Name: 友達, Total Time: 6435, 17.05, 16.95
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: 134, Stock: 2023, Name: 燁輝, Total Time: 3027, 25.55, 25.75
WARN[2021/10/20 01:44:33] 2021-10-13 Forward Balance: -141, Stock: 3481, Name: 群創, Total Time: 7609, 15.75, 15.65
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: 361, Stock: 2371, Name: 大同, Total Time: 2605, 33.90, 34.35
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: -48, Stock: 2618, Name: 長榮航, Total Time: 2188, 18.75, 18.75
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: 550, Stock: 1304, Name: 台聚, Total Time: 8813, 38.40, 39.05
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: -533, Stock: 6443, Name: 元晶, Total Time: 2699, 32.15, 31.70
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: -874, Stock: 8121, Name: 越峰, Total Time: 1075, 48.70, 47.95
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: 1029, Stock: 1308, Name: 亞聚, Total Time: 4052, 46.05, 47.20
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: 1339, Stock: 1305, Name: 華夏, Total Time: 10320, 41.50, 42.95
WARN[2021/10/20 01:44:33] 2021-10-13 Reverse Balance: 338, Stock: 2406, Name: 國碩, Total Time: 9817, 23.90, 24.30
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: 660, Stock: 2371, Name: 大同, Total Time: 1389, 34.15, 34.90
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -191, Stock: 3481, Name: 群創, Total Time: 4381, 16.50, 16.35
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -299, Stock: 1304, Name: 台聚, Total Time: 10705, 38.70, 38.50
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -114, Stock: 2027, Name: 大成鋼, Total Time: 1911, 44.50, 44.50
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -225, Stock: 2641, Name: 正德, Total Time: 686, 28.90, 28.75
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -194, Stock: 2605, Name: 新興, Total Time: 4249, 36.50, 36.40
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -495, Stock: 2409, Name: 友達, Total Time: 8504, 17.55, 17.10
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -316, Stock: 2344, Name: 華邦電, Total Time: 5359, 25.65, 25.40
WARN[2021/10/20 01:44:33] 2021-10-08 Forward Balance: -339, Stock: 2002, Name: 中鋼, Total Time: 1900, 34.35, 34.10
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: -114, Stock: 2201, Name: 裕隆, Total Time: 3713, 44.50, 44.50
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 385, Stock: 5351, Name: 鈺創, Total Time: 3592, 44.35, 44.85
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: -958, Stock: 2328, Name: 廣宇, Total Time: 1990, 42.25, 41.40
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 1423, Stock: 5608, Name: 四維航, Total Time: 11209, 47.95, 49.50
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 37, Stock: 1710, Name: 東聯, Total Time: 7013, 24.50, 24.60
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 38, Stock: 2406, Name: 國碩, Total Time: 8181, 24.15, 24.25
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 411, Stock: 1727, Name: 中華化, Total Time: 1488, 33.95, 34.45
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 562, Stock: 4503, Name: 金雨, Total Time: 4389, 33.50, 34.15
WARN[2021/10/20 01:44:33] 2021-10-08 Reverse Balance: 282, Stock: 1711, Name: 永光, Total Time: 5077, 25.75, 26.10
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: -306, Stock: 1584, Name: 精剛, Total Time: 9290, 21.90, 21.65
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: 214, Stock: 2371, Name: 大同, Total Time: 2334, 33.20, 33.50
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: 239, Stock: 1710, Name: 東聯, Total Time: 3672, 23.40, 23.70
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: -514, Stock: 1305, Name: 華夏, Total Time: 7684, 44.90, 44.50
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: -146, Stock: 2409, Name: 友達, Total Time: 4176, 17.70, 17.60
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: -348, Stock: 1304, Name: 台聚, Total Time: 8856, 37.85, 37.60
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: -746, Stock: 1455, Name: 集盛, Total Time: 9860, 18.30, 17.60
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: 1026, Stock: 1308, Name: 亞聚, Total Time: 821, 46.95, 48.10
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: 334, Stock: 1711, Name: 永光, Total Time: 1418, 25.25, 25.65
WARN[2021/10/20 01:44:33] 2021-10-06 Forward Balance: -390, Stock: 2002, Name: 中鋼, Total Time: 5111, 35.25, 34.95
WARN[2021/10/20 01:44:33] 2021-10-06 Reverse Balance: 1493, Stock: 5351, Name: 鈺創, Total Time: 8932, 40.35, 41.95
WARN[2021/10/20 01:44:33] 2021-10-06 Reverse Balance: 888, Stock: 1440, Name: 南紡, Total Time: 1886, 23.05, 24.00
WARN[2021/10/20 01:44:33] 2021-10-06 Reverse Balance: 348, Stock: 2328, Name: 廣宇, Total Time: 10156, 39.20, 39.65
WARN[2021/10/20 01:44:33] 2021-10-06 Reverse Balance: 157, Stock: 3481, Name: 群創, Total Time: 8297, 16.55, 16.75
WARN[2021/10/20 01:44:33] 2021-10-06 Reverse Balance: -830, Stock: 1727, Name: 中華化, Total Time: 13387, 31.55, 30.80
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: -129, Stock: 2486, Name: 一詮, Total Time: 112, 50.10, 50.10
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: 308, Stock: 3481, Name: 群創, Total Time: 7954, 15.65, 16.00
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: -511, Stock: 2406, Name: 國碩, Total Time: 7023, 24.25, 23.80
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: 634, Stock: 5608, Name: 四維航, Total Time: 8163, 44.45, 45.20
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: 659, Stock: 2371, Name: 大同, Total Time: 800, 34.65, 35.40
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: 627, Stock: 1711, Name: 永光, Total Time: 253, 27.95, 28.65
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: 5, Stock: 2409, Name: 友達, Total Time: 4176, 17.00, 17.05
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: -107, Stock: 5351, Name: 鈺創, Total Time: 6985, 40.90, 40.90
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: -184, Stock: 2002, Name: 中鋼, Total Time: 8794, 32.75, 32.65
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: -437, Stock: 4714, Name: 永捷, Total Time: 10757, 14.45, 14.05
WARN[2021/10/20 01:44:33] 2021-10-14 Forward Balance: -7, Stock: 2328, Name: 廣宇, Total Time: 251, 40.80, 40.90
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 768, Stock: 6443, Name: 元晶, Total Time: 6049, 30.85, 31.70
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 1641, Stock: 1709, Name: 和益, Total Time: 3552, 21.25, 22.95
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 318, Stock: 2605, Name: 新興, Total Time: 9458, 31.20, 31.60
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: -65, Stock: 2027, Name: 大成鋼, Total Time: 1384, 44.80, 44.85
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 212, Stock: 4513, Name: 福裕, Total Time: 1181, 14.90, 15.15
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 691, Stock: 1308, Name: 亞聚, Total Time: 7515, 41.55, 42.35
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: -165, Stock: 2344, Name: 華邦電, Total Time: 9478, 25.15, 25.05
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 3545, Stock: 4503, Name: 金雨, Total Time: 3631, 37.25, 40.90
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: -451, Stock: 1305, Name: 華夏, Total Time: 4658, 39.55, 39.20
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 3774, Stock: 3701, Name: 大眾控, Total Time: 3617, 45.45, 49.35
WARN[2021/10/20 01:44:33] 2021-10-14 Reverse Balance: 2438, Stock: 1727, Name: 中華化, Total Time: 3091, 41.30, 43.85
WARN[2021/10/20 01:44:33] 2021-10-12 Forward Balance: -191, Stock: 2337, Name: 旺宏, Total Time: 597, 35.55, 35.45
WARN[2021/10/20 01:44:33] 2021-10-12 Forward Balance: -314, Stock: 5351, Name: 鈺創, Total Time: 4107, 44.20, 44.00
WARN[2021/10/20 01:44:33] 2021-10-12 Forward Balance: -145, Stock: 2409, Name: 友達, Total Time: 6252, 17.10, 17.00
WARN[2021/10/20 01:44:33] 2021-10-12 Forward Balance: -849, Stock: 1304, Name: 台聚, Total Time: 8407, 39.00, 38.25
WARN[2021/10/20 01:44:33] 2021-10-12 Forward Balance: -197, Stock: 2618, Name: 長榮航, Total Time: 1521, 18.35, 18.20
WARN[2021/10/20 01:44:33] 2021-10-12 Forward Balance: -975, Stock: 1308, Name: 亞聚, Total Time: 2162, 49.05, 48.20
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: -88, Stock: 2371, Name: 大同, Total Time: 5992, 33.95, 33.95
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: 1034, Stock: 1305, Name: 華夏, Total Time: 1715, 44.00, 45.15
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: -112, Stock: 3701, Name: 大眾控, Total Time: 6189, 43.25, 43.25
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: 134, Stock: 1711, Name: 永光, Total Time: 6124, 25.60, 25.80
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: 690, Stock: 2328, Name: 廣宇, Total Time: 2978, 41.90, 42.70
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: -967, Stock: 8121, Name: 越峰, Total Time: 13003, 45.85, 45.00
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: 214, Stock: 2002, Name: 中鋼, Total Time: 2918, 33.25, 33.55
WARN[2021/10/20 01:44:33] 2021-10-12 Reverse Balance: 259, Stock: 1727, Name: 中華化, Total Time: 1562, 35.25, 35.60
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: 560, Stock: 2371, Name: 大同, Total Time: 815, 34.10, 34.75
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -61, Stock: 2406, Name: 國碩, Total Time: 7997, 23.45, 23.45
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -410, Stock: 1710, Name: 東聯, Total Time: 6184, 23.45, 23.10
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -146, Stock: 2409, Name: 友達, Total Time: 6335, 17.60, 17.50
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -242, Stock: 3481, Name: 群創, Total Time: 6946, 16.60, 16.40
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -119, Stock: 2356, Name: 英業達, Total Time: 3161, 26.40, 26.35
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -516, Stock: 1711, Name: 永光, Total Time: 776, 25.80, 25.35
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: 327, Stock: 1308, Name: 亞聚, Total Time: 6556, 47.20, 47.65
WARN[2021/10/20 01:44:33] 2021-10-07 Forward Balance: -458, Stock: 3232, Name: 昱捷, Total Time: 9121, 42.10, 41.75
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -110, Stock: 2201, Name: 裕隆, Total Time: 4496, 42.80, 42.80
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: 389, Stock: 1440, Name: 南紡, Total Time: 9535, 23.40, 23.85
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -913, Stock: 1305, Name: 華夏, Total Time: 6142, 44.30, 43.50
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -108, Stock: 5351, Name: 鈺創, Total Time: 898, 41.90, 41.90
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: 279, Stock: 2023, Name: 燁輝, Total Time: 3155, 27.40, 27.75
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: 34, Stock: 2344, Name: 華邦電, Total Time: 1714, 25.45, 25.55
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: 205, Stock: 2605, Name: 新興, Total Time: 1651, 36.50, 36.80
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -153, Stock: 2328, Name: 廣宇, Total Time: 274, 40.00, 39.95
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -299, Stock: 2014, Name: 中鴻, Total Time: 8400, 38.20, 38.00
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -834, Stock: 1727, Name: 中華化, Total Time: 1308, 33.40, 32.65
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: 323, Stock: 5608, Name: 四維航, Total Time: 4377, 48.55, 49.00
WARN[2021/10/20 01:44:33] 2021-10-07 Reverse Balance: -898, Stock: 1304, Name: 台聚, Total Time: 7740, 38.65, 37.85
WARN[2021/10/20 01:44:33] 2021-10-05 Forward Balance: 516, Stock: 2371, Name: 大同, Total Time: 6263, 32.10, 32.70
WARN[2021/10/20 01:44:33] 2021-10-05 Forward Balance: 460, Stock: 2002, Name: 中鋼, Total Time: 4523, 34.65, 35.20
WARN[2021/10/20 01:44:33] 2021-10-05 Forward Balance: 954, Stock: 1304, Name: 台聚, Total Time: 10588, 36.20, 37.25
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -38, Stock: 6116, Name: 彩晶, Total Time: 6804, 14.70, 14.70
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -655, Stock: 1584, Name: 精剛, Total Time: 2464, 21.30, 20.70
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -520, Stock: 2023, Name: 燁輝, Total Time: 9441, 27.35, 26.90
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -345, Stock: 3481, Name: 群創, Total Time: 3410, 17.15, 16.85
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -96, Stock: 1455, Name: 集盛, Total Time: 90, 17.75, 17.70
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -216, Stock: 2344, Name: 華邦電, Total Time: 7974, 25.60, 25.45
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -300, Stock: 2617, Name: 台航, Total Time: 1540, 39.25, 39.05
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -484, Stock: 2614, Name: 東森, Total Time: 7289, 32.65, 32.25
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -241, Stock: 2605, Name: 新興, Total Time: 9957, 35.70, 35.55
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -112, Stock: 2363, Name: 矽統, Total Time: 3068, 24.05, 24.00
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: 5, Stock: 2610, Name: 華航, Total Time: 2911, 16.90, 16.95
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -446, Stock: 2409, Name: 友達, Total Time: 2273, 17.75, 17.35
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -414, Stock: 1711, Name: 永光, Total Time: 2280, 24.95, 24.60
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -655, Stock: 5351, Name: 鈺創, Total Time: 10418, 41.10, 40.55
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: 343, Stock: 1440, Name: 南紡, Total Time: 7773, 22.00, 22.40
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -599, Stock: 2014, Name: 中鴻, Total Time: 8516, 38.70, 38.20
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -71, Stock: 1802, Name: 台玻, Total Time: 3134, 27.90, 27.90
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: 684, Stock: 1308, Name: 亞聚, Total Time: 2156, 44.30, 45.10
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: 786, Stock: 1305, Name: 華夏, Total Time: 5242, 43.35, 44.25
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -623, Stock: 5608, Name: 四維航, Total Time: 9857, 48.30, 47.80
WARN[2021/10/20 01:44:33] 2021-10-05 Reverse Balance: -97, Stock: 2618, Name: 長榮航, Total Time: 6884, 18.65, 18.60
WARN[2021/10/20 01:44:33] Total Balance: 18235, TradeCount: 187
WARN[2021/10/20 01:44:33] Cond: {Model:{ID:34430 CreatedAt:2021-10-20 01:31:26.794875 +0800 CST UpdatedAt:2021-10-20 01:31:26.794875 +0800 CST DeletedAt:{Time:0001-01-01 00:00:00 +0000 UTC Valid:false}} HistoryCloseCount:2500 OutInRatio:70 ReverseOutInRatio:5 CloseDiff:0 CloseChangeRatioLow:-1 CloseChangeRatioHigh:8 OpenChangeRatio:4 RsiHigh:50.4 RsiLow:50 ReverseRsiHigh:50.4 ReverseRsiLow:50 TicksPeriodThreshold:1 TicksPeriodLimit:1.3 TicksPeriodCount:1 VolumePerSecond:6}
```

### Trade Bot Service

![callvis](./assets/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
