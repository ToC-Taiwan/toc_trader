// Package tickanalyze package tickanalyze
package tickanalyze

import (
	"errors"

	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
)

// GenerateRSI GenerateRSI
func GenerateRSI(input quote.Quote) (rsi float64, err error) {
	rsiArr := talib.Rsi(input.Close, len(input.Close)-1)
	if len(rsiArr) == 0 {
		return 0, errors.New("no rsi")
	}
	return utils.Round(rsiArr[len(rsiArr)-1], 2), err
}

// GetForwardRSIStatus GetForwardRSIStatus
func GetForwardRSIStatus(input []float64, rsiHighLimit float64) bool {
	var result int
	var resultArr, tmp []float64
	part := 11
	divide := len(input) / part
	for i := 1; i <= part; i++ {
		tmp = input[len(input)-i*divide : len(input)-(i-1)*divide]
		var quoteArr quote.Quote
		quoteArr.Close = tmp
		rsi, err := GenerateRSI(quoteArr)
		if err != nil {
			logger.GetLogger().Error(err)
			return false
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
		if v == 100 {
			result++
		}
	}
	return float64(result)/float64(part) > rsiHighLimit
}

// GetReverseRSIStatus GetReverseRSIStatus
func GetReverseRSIStatus(input []float64, rsiLowLimit float64) bool {
	var result int
	var resultArr, tmp []float64
	part := 11
	divide := len(input) / part
	for i := 1; i <= part; i++ {
		tmp = input[len(input)-i*divide : len(input)-(i-1)*divide]
		var quoteArr quote.Quote
		quoteArr.Close = tmp
		rsi, err := GenerateRSI(quoteArr)
		if err != nil {
			logger.GetLogger().Error(err)
			return false
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
		if v == 0 {
			result--
		}
	}
	return float64(result)/float64(part) < rsiLowLimit
}

// GenerareMAByCount GenerareMAByCount
func GenerareMAByCount(input []float64, n int) (lastMa float64, err error) {
	if len(input) == 0 || len(input) == 1 {
		return 0, errors.New("input is empty or length is 1")
	}
	maArr := talib.Ma(input, n, talib.SMA)
	if len(maArr) == 0 {
		return 0, errors.New("no ma")
	}
	return maArr[len(maArr)-1], err
}
