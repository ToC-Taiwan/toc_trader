// Package sysparm package sysparm
package sysparm

// DefaultKey DefaultKey
var DefaultKey = []string{
	"runmode",
	"reset",
	"dbuser",
	"dbpassword",
	"dbhost",
	"dbport",
	"database",
	"dbencode",
	"dbtimezone",
	"target_condition",
	"black_stock_arr",
	"black_category_arr",
	"cleanevent_cron",
	"restart_sinopac_toc_trader_cron",
	"http_port",
	"py_server_port",
	"py_server_host",
}

// DefaultSetting DefaultSetting
var DefaultSetting = map[string]interface{}{
	"runmode":                         "debug",
	"reset":                           0,
	"dbuser":                          "postgres",
	"dbpassword":                      "asdf0000",
	"dbhost":                          "172.20.10.10",
	"dbport":                          "5432",
	"database":                        "tradebot_debug",
	"dbencode":                        "utf8",
	"dbtimezone":                      "Asia/Taipei",
	"target_condition":                `[{"limit_price_low":50,"limit_price_high":100,"limit_volume":10000}]`,
	"black_stock_arr":                 `["1314","2317"]`,
	"black_category_arr":              `["17"]`,
	"cleanevent_cron":                 "0 0 4 * * ?",
	"restart_sinopac_toc_trader_cron": "0 20 8 * * ?",
	"http_port":                       "6670",
	"py_server_port":                  "3333",
	"py_server_host":                  "127.0.0.1",
}
