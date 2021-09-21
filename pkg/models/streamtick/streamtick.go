// Package streamtick package streamtick
package streamtick

import (
	"gitlab.tocraw.com/root/toc_trader/tools/common"
	"gorm.io/gorm"
)

// Old Old
type Old struct {
	gorm.Model
	StockNum  string  `gorm:"column:stock_num"`
	AmountSum float64 `gorm:"column:amount_sum"`
	Close     float64 `gorm:"column:close"`
	TickType  int64   `gorm:"column:tick_type"`
	VolumeSum int64   `gorm:"column:volume_sum"`
	Volume    int64   `gorm:"column:volume"`
	TimeStamp int64   `gorm:"column:timestamp"`
}

// StreamTick StreamTick
type StreamTick struct {
	gorm.Model
	StockNum        string  `gorm:"column:stock_num;index:idx_streamtick"`
	TimeStamp       int64   `gorm:"column:timestamp;index:idx_streamtick"`
	Open            float64 `gorm:"column:open"`
	AvgPrice        float64 `gorm:"column:avg_price"`
	Close           float64 `gorm:"column:close"`
	High            float64 `gorm:"column:high"`
	Low             float64 `gorm:"column:low"`
	Amount          float64 `gorm:"column:amount"`
	AmountSum       float64 `gorm:"column:amount_sum"`
	Volume          int64   `gorm:"column:volume"`
	VolumeSum       int64   `gorm:"column:volume_sum"`
	TickType        int64   `gorm:"column:tick_type"`
	ChgType         int64   `gorm:"column:chg_type"`
	PriceChg        float64 `gorm:"column:price_chg"`
	PctChg          float64 `gorm:"column:pct_chg"`
	BidSideTotalVol int64   `gorm:"column:bid_side_total_vol"`
	AskSideTotalVol int64   `gorm:"column:ask_side_total_vol"`
	BidSideTotalCnt int64   `gorm:"column:bid_side_total_cnt"`
	AskSideTotalCnt int64   `gorm:"column:ask_side_total_cnt"`
	Suspend         int64   `gorm:"column:suspend"`
	Simtrade        int64   `gorm:"column:simtrade"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (StreamTick) TableName() string {
	return "tick_stream"
}

// SinoPacStreamTick SinoPacStreamTick
type SinoPacStreamTick struct {
	Topic     string    `json:"topic"`
	AmountSum []float64 `json:"AmountSum"`
	Close     []float64 `json:"Close"`
	Date      string    `json:"Date"`
	Simtrade  int64     `json:"Simtrade"`
	TickType  []int64   `json:"TickType"`
	Time      string    `json:"Time"`
	VolSum    []int64   `json:"VolSum"`
	Volume    []int64   `json:"Volume"`
}

// ProtoToStreamTick ProtoToStreamTick
func (c *StreamTickProto) ProtoToStreamTick() (result *StreamTick, err error) {
	timeStamp, err := common.MicroDateTimeToTimeStamp(c.Tick.DateTime)
	if err != nil {
		return result, err
	}
	tmp := StreamTick{
		StockNum:        c.Tick.Code,
		TimeStamp:       timeStamp,
		Open:            c.Tick.Open,
		AvgPrice:        c.Tick.AvgPrice,
		Close:           c.Tick.Close,
		High:            c.Tick.High,
		Low:             c.Tick.Low,
		Amount:          c.Tick.Amount,
		AmountSum:       c.Tick.TotalAmount,
		Volume:          c.Tick.Volume,
		VolumeSum:       c.Tick.TotalVolume,
		TickType:        c.Tick.TickType,
		ChgType:         c.Tick.ChgType,
		PriceChg:        c.Tick.PriceChg,
		PctChg:          c.Tick.PctChg,
		BidSideTotalVol: c.Tick.BidSideTotalVol,
		AskSideTotalVol: c.Tick.AskSideTotalVol,
		BidSideTotalCnt: c.Tick.BidSideTotalCnt,
		AskSideTotalCnt: c.Tick.AskSideTotalCnt,
		Suspend:         c.Tick.Suspend,
		Simtrade:        c.Tick.Simtrade,
	}
	return &tmp, err
}

// ToStreamTick ToStreamTick
// func (c *SinoPacStreamTick) ToStreamTick() (result *StreamTick, err error) {
// 	var stockNum string
// 	if strings.HasPrefix(c.Topic, "MKT/idcdmzpcr01/TSE") {
// 		stockNum = strings.ReplaceAll(c.Topic, "MKT/idcdmzpcr01/TSE/", "")
// 	} else {
// 		stockNum = strings.ReplaceAll(c.Topic, "MKT/idcdmzpcr01/OTC/", "")
// 	}
// 	replaceDate := strings.ReplaceAll(c.Date, "/", "-")
// 	utcTimeSec, err := time.ParseInLocation(global.LongTimeLayout, replaceDate+" "+c.Time[:8], time.Local)
// 	if err != nil {
// 		return result, err
// 	}
// 	ms, err := common.StrToFloat64("0" + c.Time[8:])
// 	if err != nil {
// 		return result, err
// 	}
// 	if c.Simtrade == 1 {
// 		logger.Logger.WithFields(map[string]interface{}{
// 			"TickType": c.TickType[0],
// 			"Volume":   c.Volume[0],
// 			"Close":    c.Close[0],
// 			"Name":     global.AllStockNameMap.GetName(stockNum),
// 		}).Info("SimTrade")
// 		return result, err
// 	}

// 	result = &StreamTick{
// 		StockNum:  stockNum,
// 		AmountSum: c.AmountSum[0],
// 		Close:     c.Close[0],
// 		TickType:  c.TickType[0],
// 		VolumeSum: c.VolSum[0],
// 		Volume:    c.Volume[0],
// 		TimeStamp: int64(float64(utcTimeSec.UnixNano() + int64(ms*1000*1000*1000))),
// 	}
// 	return result, err
// }

// PtrArr PtrArr
type PtrArr []*StreamTick

// GetOutInRatio GetOutInRatio
func (c *PtrArr) GetOutInRatio() float64 {
	var outSum, inSum int64
	for _, v := range *c {
		if v.TickType == 1 {
			outSum += v.Volume
		} else {
			inSum += v.Volume
		}
	}
	return float64(outSum) / float64(outSum+inSum)
}

// GetOutSum GetOutSum
func (c *PtrArr) GetOutSum() int64 {
	var outSum int64
	for _, v := range *c {
		if v.TickType == 1 {
			outSum += v.Volume
		}
	}
	return outSum
}

// GetInSum GetInSum
func (c *PtrArr) GetInSum() int64 {
	var inSum int64
	for _, v := range *c {
		if v.TickType == 2 {
			inSum += v.Volume
		}
	}
	return inSum
}

// GetLastClose GetLastClose
func (c *PtrArr) GetLastClose() float64 {
	var tmp []*StreamTick = *c
	return tmp[len(tmp)-1].Close
}

// GetAllCloseArr GetAllCloseArr
func (c *PtrArr) GetAllCloseArr() []float64 {
	var tmp []*StreamTick = *c
	var closeArr []float64
	for _, v := range tmp {
		closeArr = append(closeArr, v.Close)
	}
	return closeArr
}

// GetTotalTime GetTotalTime
func (c *PtrArr) GetTotalTime() float64 {
	tmp := *c
	return float64(tmp[len(tmp)-1].TimeStamp-tmp[0].TimeStamp) / 1000 / 1000 / 1000
}

// PtrArrArr PtrArrArr
type PtrArrArr []PtrArr

// Get Get
func (c *PtrArrArr) Get() []PtrArr {
	return *c
}

// GetCount GetCount
func (c *PtrArrArr) GetCount() int {
	return len(*c)
}

// GetLastNRow GetLastNRow
func (c *PtrArrArr) GetLastNRow(n int) []PtrArr {
	var tmp []PtrArr = *c
	var ans []PtrArr
	for i := len(*c) - 1; i > len(*c)-1-n; i-- {
		ans = append(ans, tmp[i])
	}
	return ans
}

// Append Append
func (c *PtrArrArr) Append(data PtrArr) {
	*c = append(*c, data)
}

// ClearAll ClearAll
func (c *PtrArrArr) ClearAll() {
	*c = []PtrArr{}
}

// GetCloseDiff GetCloseDiff
func (c *PtrArrArr) GetCloseDiff() float64 {
	var tmp []PtrArr = *c
	first := tmp[0][0].Close
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1].Close
	return last - first
}

// GetLastClose GetLastClose
func (c *PtrArrArr) GetLastClose() float64 {
	var tmp []PtrArr = *c
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1].Close
	return last
}

// GetLastTick GetLastTick
func (c *PtrArrArr) GetLastTick() *StreamTick {
	var tmp []PtrArr = *c
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1]
	return last
}
