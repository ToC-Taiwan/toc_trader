# TOC TRADER

[![pipeline status](https://gitlab.tocraw.com/root/toc_trader/badges/main/pipeline.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![coverage report](https://gitlab.tocraw.com/root/toc_trader/badges/main/coverage.svg)](https://gitlab.tocraw.com/root/toc_trader/-/commits/main)
[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://gitlab.tocraw.com/root/toc_trader)
[![Go](https://img.shields.io/badge/Go-1.17.1-blue?logo=go&logoColor=blue)](https://golang.org)
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

```log
WARN[2021/09/30 16:00:51] Balance: -475, Stock: 6182, Name: 合晶, Total Time: 1720
WARN[2021/09/30 16:00:51] Balance: -274, Stock: 2303, Name: 聯電, Total Time: 467
WARN[2021/09/30 16:00:51] Balance: 611, Stock: 2436, Name: 偉詮電, Total Time: 529
WARN[2021/09/30 16:00:51] Balance: 418, Stock: 2606, Name: 裕民, Total Time: 1029
WARN[2021/09/30 16:00:51] Balance: 472, Stock: 6217, Name: 中探針, Total Time: 230
WARN[2021/09/30 16:00:51] Balance: 459, Stock: 2486, Name: 一詮, Total Time: 811
WARN[2021/09/30 16:00:51] Total Balance: 1211, TradeCount: 6, Cond: {500 150 55 0 -3 5 5 50 50}
```

```log
WARN[2021/10/01 03:10:33] Balance: -361, Stock: 1513, Name: 中興電, Total Time: 745
WARN[2021/10/01 03:10:33] Balance: 35, Stock: 2601, Name: 益航, Total Time: 3992
WARN[2021/10/01 03:10:33] Balance: -117, Stock: 9103, Name: 美德醫療-DR, Total Time: 476
WARN[2021/10/01 03:10:33] Balance: 135, Stock: 2344, Name: 華邦電, Total Time: 645
WARN[2021/10/01 03:10:33] Balance: 243, Stock: 2504, Name: 國產, Total Time: 567
WARN[2021/10/01 03:10:33] Balance: 564, Stock: 2641, Name: 正德, Total Time: 412
WARN[2021/10/01 03:10:33] Balance: 140, Stock: 1440, Name: 南紡, Total Time: 77
WARN[2021/10/01 03:10:33] Balance: 104, Stock: 1718, Name: 中纖, Total Time: 376
WARN[2021/10/01 03:10:33] Balance: -123, Stock: 2610, Name: 華航, Total Time: 775
WARN[2021/10/01 03:10:33] Balance: -190, Stock: 2371, Name: 大同, Total Time: 960
WARN[2021/10/01 03:10:33] Balance: 248, Stock: 1444, Name: 力麗, Total Time: 157
WARN[2021/10/01 03:10:33] Balance: -60, Stock: 2605, Name: 新興, Total Time: 497
WARN[2021/10/01 03:10:33] Balance: -610, Stock: 1308, Name: 亞聚, Total Time: 159
WARN[2021/10/01 03:10:33] Balance: 114, Stock: 1734, Name: 杏輝, Total Time: 81
WARN[2021/10/01 03:10:33] Balance: -151, Stock: 1305, Name: 華夏, Total Time: 87
WARN[2021/10/01 03:10:33] Balance: 52, Stock: 1710, Name: 東聯, Total Time: 127
WARN[2021/10/01 03:10:33] Balance: -228, Stock: 1802, Name: 台玻, Total Time: 2039
WARN[2021/10/01 03:10:33] Balance: 49, Stock: 2406, Name: 國碩, Total Time: 1449
WARN[2021/10/01 03:10:33] Balance: 20, Stock: 2618, Name: 長榮航, Total Time: 664
WARN[2021/10/01 03:10:33] Balance: -179, Stock: 2409, Name: 友達, Total Time: 126
WARN[2021/10/01 03:10:33] Balance: -42, Stock: 1409, Name: 新纖, Total Time: 245
WARN[2021/10/01 03:10:33] Balance: 93, Stock: 1711, Name: 永光, Total Time: 628
WARN[2021/10/01 03:10:33] Balance: 154, Stock: 2027, Name: 大成鋼, Total Time: 2032
WARN[2021/10/01 03:10:33] Balance: 38, Stock: 1905, Name: 華紙, Total Time: 299
WARN[2021/10/01 03:10:33] Balance: 171, Stock: 2023, Name: 燁輝, Total Time: 355
WARN[2021/10/01 03:10:33] Balance: -37, Stock: 6126, Name: 信音, Total Time: 791
WARN[2021/10/01 03:10:33] Balance: 313, Stock: 1455, Name: 集盛, Total Time: 107
WARN[2021/10/01 03:10:33] Balance: -26, Stock: 3481, Name: 群創, Total Time: 39
WARN[2021/10/01 03:10:33] Balance: 662, Stock: 4141, Name: 龍燈-KY, Total Time: 250
WARN[2021/10/01 03:10:33] Balance: 76, Stock: 3033, Name: 威健, Total Time: 987
WARN[2021/10/01 03:10:33] Balance: -153, Stock: 2002, Name: 中鋼, Total Time: 352
WARN[2021/10/01 03:10:33] Balance: -42, Stock: 3062, Name: 建漢, Total Time: 176
WARN[2021/10/01 03:10:33] Balance: -216, Stock: 6116, Name: 彩晶, Total Time: 2720
WARN[2021/10/01 03:10:33] Total Balance: 676, TradeCount: 33, Cond: {400 150 55 0 -1 7 7 45 55}
```

```log
WARN[2021/10/01 03:49:41] Balance: -293, Stock: 2408, Name: 南亞科, Total Time: 329
WARN[2021/10/01 03:49:41] Balance: 225, Stock: 6182, Name: 合晶, Total Time: 779
WARN[2021/10/01 03:49:41] Balance: 95, Stock: 2436, Name: 偉詮電, Total Time: 58
WARN[2021/10/01 03:49:41] Balance: -81, Stock: 2606, Name: 裕民, Total Time: 493
WARN[2021/10/01 03:49:41] Balance: 217, Stock: 2368, Name: 金像電, Total Time: 1788
WARN[2021/10/01 03:49:41] Balance: 472, Stock: 6217, Name: 中探針, Total Time: 123
WARN[2021/10/01 03:49:41] Balance: 559, Stock: 2486, Name: 一詮, Total Time: 795
WARN[2021/10/01 03:49:41] Balance: -272, Stock: 2303, Name: 聯電, Total Time: 436
WARN[2021/10/01 03:49:41] Balance: 92, Stock: 2354, Name: 鴻準, Total Time: 117
WARN[2021/10/01 03:49:41] Total Balance: 1014, TradeCount: 9, Cond: {500 150 55 0 -3 5 5 45 55 5 6.5 4}
```

### Trade Bot Service

![callvis](./assets/callvis/callvis.svg "callvis")

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
