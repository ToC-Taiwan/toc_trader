// Package entiretick package entiretick
package entiretick

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gorm.io/gorm"
)

// EntireTick EntireTick
type EntireTick struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	StockNum   string  `gorm:"column:stock_num;index:idx_entiretick"`
	Close      float64 `gorm:"column:close"`
	TickType   int64   `gorm:"column:tick_type"`
	Volume     int64   `gorm:"column:volume"`
	BidPrice   float64 `gorm:"column:bid_price"`
	BidVolume  int64   `gorm:"column:bid_volume"`
	AskPrice   float64 `gorm:"column:ask_price"`
	AskVolume  int64   `gorm:"column:ask_volume"`
	TimeStamp  int64   `gorm:"column:timestamp;index:idx_entiretick"`
	Open       float64 `gorm:"column:open"`
	High       float64 `gorm:"column:high"`
	Low        float64 `gorm:"column:low"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (EntireTick) TableName() string {
	return "tick_entire"
}

// ProtoToEntireTick ProtoToEntireTick
func (c *EntireTickProto) ProtoToEntireTick(stockNum string) (tick *EntireTick, err error) {
	// var tickType int64
	// if c.Close == c.AskPrice {
	// 	tickType = 1
	// } else {
	// 	tickType = 2
	// }
	tmp := EntireTick{
		StockNum:  stockNum,
		Close:     c.Close,
		TickType:  c.TickType,
		Volume:    c.Volume,
		TimeStamp: c.Ts,
		AskPrice:  c.AskPrice,
		AskVolume: c.AskVolume,
		BidPrice:  c.BidPrice,
		BidVolume: c.BidVolume,
	}
	return &tmp, err
}

// PtrArr PtrArr
type PtrArr []*EntireTick

// GetOutInRatio GetOutInRatio
func (c *PtrArr) GetOutInRatio() float64 {
	var outSum, inSum int64
	for _, v := range *c {
		switch v.TickType {
		case 0:
			continue
		case 1:
			outSum += v.Volume
		case 2:
			inSum += v.Volume
		}
	}
	return float64(outSum) / float64(outSum+inSum)
}

// GetOutSum GetOutSum
func (c *PtrArr) GetOutSum() int64 {
	var outSum int64
	for _, v := range *c {
		switch v.TickType {
		case 0:
			continue
		case 1:
			outSum += v.Volume
		}
	}
	return outSum
}

// GetInSum GetInSum
func (c *PtrArr) GetInSum() int64 {
	var inSum int64
	for _, v := range *c {
		switch v.TickType {
		case 0:
			continue
		case 2:
			inSum += v.Volume
		}
	}
	return inSum
}

// GetLastClose GetLastClose
func (c *PtrArr) GetLastClose() float64 {
	var tmp []*EntireTick = *c
	return tmp[len(tmp)-1].Close
}

// GetAllCloseArr GetAllCloseArr
func (c *PtrArr) GetAllCloseArr() []float64 {
	var tmp []*EntireTick = *c
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
func (c *PtrArrArr) GetLastTick() *EntireTick {
	var tmp []PtrArr = *c
	last := tmp[len(tmp)-1][len(tmp[len(tmp)-1])-1]
	return last
}

// ToStreamTick ToStreamTick
func (c *EntireTick) ToStreamTick() *streamtick.StreamTick {
	tickTime := time.Unix(0, c.TimeStamp).Add(-8 * time.Hour).UnixNano()
	tmp := streamtick.StreamTick{
		StockNum:        c.StockNum,
		TimeStamp:       tickTime,
		Open:            c.Open,
		AvgPrice:        0,
		Close:           c.Close,
		High:            c.High,
		Low:             c.Low,
		Amount:          0,
		AmountSum:       0,
		Volume:          c.Volume,
		VolumeSum:       0,
		TickType:        c.TickType,
		ChgType:         0,
		PriceChg:        0,
		PctChg:          0,
		BidSideTotalVol: c.BidVolume,
		AskSideTotalVol: c.AskVolume,
		BidSideTotalCnt: 0,
		AskSideTotalCnt: 0,
		Suspend:         0,
		Simtrade:        0,
	}
	return &tmp
}
