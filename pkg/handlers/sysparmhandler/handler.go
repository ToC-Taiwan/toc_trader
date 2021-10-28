// Package sysparmhandler sysparmhandler
package sysparmhandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/parameters"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/system/sysparm", UpdateSysparm)
}

// UpdateSysparm UpdateSysparm
// @Summary UpdateSysparm
// @tags sysparm
// @accept json
// @produce json
// @param body body []sysparm.Parameters true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /system/sysparm [post]
func UpdateSysparm(c *gin.Context) {
	var req []sysparm.Parameters
	var res handlers.ErrorResponse
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := json.Unmarshal(byteArr, &req); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	for _, v := range req {
		if err := parameters.UpdateSysparm(v.Key, v.Value); err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}
