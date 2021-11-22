// Package sysparm package sysparm
package sysparm

import (
	"encoding/json"

	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
)

// GlobalSettingMap GlobalSettingMap
type GlobalSettingMap map[string]string

// GetRunMode GetRunMode
func (c GlobalSettingMap) GetRunMode() string {
	return c["runmode"]
}

// GetResetParm GetResetParm
func (c GlobalSettingMap) GetResetParm() bool {
	reset, err := common.StrToInt64(c["reset"])
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if reset == 1 {
		return true
	}
	return false
}

// GetDBUser GetDBUser
func (c GlobalSettingMap) GetDBUser() string {
	return c["dbuser"]
}

// GetDBPass GetDBPass
func (c GlobalSettingMap) GetDBPass() string {
	return c["dbpassword"]
}

// GetDBHost GetDBHost
func (c GlobalSettingMap) GetDBHost() string {
	return c["dbhost"]
}

// GetDBPort GetDBPort
func (c GlobalSettingMap) GetDBPort() string {
	return c["dbport"]
}

// GetDBName GetDBName
func (c GlobalSettingMap) GetDBName() string {
	return c["database"]
}

// GetDBEncode GetDBEncode
func (c GlobalSettingMap) GetDBEncode() string {
	return c["dbencode"]
}

// GetDBTimeZone GetDBTimeZone
func (c GlobalSettingMap) GetDBTimeZone() string {
	return c["dbtimezone"]
}

// GetKbarPeriod GetKbarPeriod
func (c GlobalSettingMap) GetKbarPeriod() int64 {
	tmp, err := common.StrToInt64(c["kbar_period"])
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	return tmp
}

// GetCleanEventCron GetCleanEventCron
func (c GlobalSettingMap) GetCleanEventCron() string {
	return c["cleanevent_cron"]
}

// GetRestartSinopacCron GetRestartSinopacCron
func (c GlobalSettingMap) GetRestartSinopacCron() string {
	return c["restart_sinopac_cron"]
}

// GetHTTPPort GetHTTPPort
func (c GlobalSettingMap) GetHTTPPort() string {
	return c["http_port"]
}

// GetPyServerHost GetPyServerHost
func (c GlobalSettingMap) GetPyServerHost() string {
	return c["py_server_host"]
}

// GetPyServerPort GetPyServerPort
func (c GlobalSettingMap) GetPyServerPort() string {
	return c["py_server_port"]
}

// GetTargetCondArr GetTargetCondArr
func (c GlobalSettingMap) GetTargetCondArr() []TargetCondArr {
	targetArrString := c["target_condition"]
	var ans []TargetCondArr
	if err := json.Unmarshal([]byte(targetArrString), &ans); err != nil {
		logger.GetLogger().Panic(err)
	}
	return ans
}

// GetBlackStockMap GetBlackStockMap
func (c GlobalSettingMap) GetBlackStockMap() map[string]string {
	BlackStockArrString := c["black_stock_arr"]
	var ans []string
	if err := json.Unmarshal([]byte(BlackStockArrString), &ans); err != nil {
		logger.GetLogger().Panic(err)
	}
	ansMap := make(map[string]string)
	for _, v := range ans {
		ansMap[v] = v
	}
	return ansMap
}

// GetBlackCategoryMap GetBlackCategoryMap
func (c GlobalSettingMap) GetBlackCategoryMap() map[string]string {
	BlackCategoryArrString := c["black_category_arr"]
	var ans []string
	if err := json.Unmarshal([]byte(BlackCategoryArrString), &ans); err != nil {
		logger.GetLogger().Panic(err)
	}
	ansMap := make(map[string]string)
	for _, v := range ans {
		ansMap[v] = v
	}
	return ansMap
}
