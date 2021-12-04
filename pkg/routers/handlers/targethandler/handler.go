// Package targethandler package targethandler
package targethandler

import (
	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/target", GetTarget)
}

// GetTarget GetTarget
// @Summary GetTarget
// @tags Target
// @accept json
// @produce json
// @param count header string true "count"
// @success 200 {object} []TargetResponse
// @failure 500 {object} handlers.ErrorResponse
// @Router /target [get]
func GetTarget(c *gin.Context) {
	var res handlers.ErrorResponse
	count := c.Request.Header.Get("count")
	countInt64, err := utils.StrToInt64(count)
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	total := len(global.TargetArr)
	response := []TargetResponse{}
	already := make(map[int64]bool)
	if int(countInt64) >= total {
		countInt64 = int64(total)
	}
	for i := 0; i < int(countInt64); i++ {
		randomBigInt, err := rand.Int(rand.Reader, big.NewInt((int64(total))))
		if err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		random := randomBigInt.Int64()
		if _, ok := already[random]; ok {
			i--
			continue
		}
		data := TargetResponse{
			StockNum: global.TargetArr[random],
			Close:    global.StockCloseByDateMap.GetClose(global.TargetArr[random], global.LastTradeDay.Format(global.ShortTimeLayout)),
		}
		already[random] = true
		response = append(response, data)
	}
	c.JSON(http.StatusOK, response)
}
