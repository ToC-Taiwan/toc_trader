// Package tickanalyze package tickanalyze
package tickanalyze

import (
	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
	"gitlab.tocraw.com/root/toc_trader/tools/common"
)

// GenerateMA GenerateMA
// func GenerateMA(input quote.Quote) float64 {
// 	maArr := talib.Ma(input.Close, 2, talib.SMA)
// 	return common.Round(maArr[len(maArr)-1], 2)
// }

// GenerateRSI GenerateRSI
func GenerateRSI(input quote.Quote) float64 {
	rsiArr := talib.Rsi(input.Close, len(input.Close)-1)
	return common.Round(rsiArr[len(rsiArr)-1], 2)
}

// GenerateBBAND GenerateBBAND
// func GenerateBBAND(input quote.Quote) []float64 {
// 	high, medium, low := talib.BBands(input.Close, 2, 2, 2, talib.SMA)
// 	var tmp []float64
// 	tmp = append(tmp, common.Round(high[len(high)-1], 2), common.Round(medium[len(medium)-1], 2), common.Round(low[len(low)-1], 2))
// 	return tmp
// }
