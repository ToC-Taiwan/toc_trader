// Package tickanalyze package tickanalyze
package tickanalyze

import (
	"errors"

	"github.com/markcheno/go-quote"
	"github.com/markcheno/go-talib"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
)

// GenerateRSI GenerateRSI
func GenerateRSI(input quote.Quote) (rsi float64, err error) {
	rsiArr := talib.Rsi(input.Close, len(input.Close)-1)
	if len(rsiArr) == 0 {
		return 0, errors.New("no rsi")
	}
	return common.Round(rsiArr[len(rsiArr)-1], 2), err
}
