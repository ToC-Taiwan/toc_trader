// Package tickanalyze package tickanalyze
package tickanalyze

import (
	"errors"

	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
)

// GenerateRSI GenerateRSI
func GenerateRSI(input quote.Quote) (rsi float64, err error) {
	rsiArr := talib.Rsi(input.Close, len(input.Close)-1)
	if len(rsiArr) == 0 {
		return 0, errors.New("no rsi")
	}
	return common.Round(rsiArr[len(rsiArr)-1], 2), err
}

// GetRSIStatus GetRSIStatus
func GetRSIStatus(input []float64, rsiHighLimit, rsiLowLimit float64) (highStatus, lowStatus bool) {
	var result int
	var resultArr, tmp []float64
	part := 10
	divide := len(input) / part
	for i := 1; i <= part; i++ {
		tmp = input[len(input)-i*divide : len(input)-(i-1)*divide]
		var quoteArr quote.Quote
		quoteArr.Close = tmp
		rsi, err := GenerateRSI(quoteArr)
		if err != nil {
			logger.GetLogger().Error(err)
			return false, false
		}
		resultArr = append(resultArr, rsi)
	}
	var totalRsi float64
	for _, v := range resultArr {
		totalRsi += v
	}
	average := totalRsi / float64(len(resultArr))
	for _, v := range resultArr {
		if v > average {
			result++
		}
	}
	return float64(result)/float64(part) > rsiHighLimit, float64(result)/float64(part) < rsiLowLimit
}
