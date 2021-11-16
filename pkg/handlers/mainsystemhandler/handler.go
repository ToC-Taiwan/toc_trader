// Package mainsystemhandler main handler
package mainsystemhandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/parameters"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/system/restart", Restart)
	group.POST("/system/sysparm", UpdateSysparm)
}

// Restart Restart
// @Summary Restart
// @tags mainsystem
// @accept json
// @produce json
// @success 200
// @failure 500 {object} string
// @Router /system/restart [get]
func Restart(c *gin.Context) {
	var res handlers.ErrorResponse
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		res.Response = "you should be in the docker container(from toc_trader)"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	healthcheck.ExitService()
	c.JSON(http.StatusOK, nil)
}

// UpdateSysparm UpdateSysparm
// @Summary UpdateSysparm
// @tags mainsystem
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
