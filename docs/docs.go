// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/balance": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Balance"
                ],
                "summary": "GetAllBalance",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/balance.Balance"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Balance"
                ],
                "summary": "ImportBalance",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/balance.Balance"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Balance"
                ],
                "summary": "DeletaAllBalance",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/condition": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TradeCondition"
                ],
                "summary": "GetLatestTradeCondition",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/simulate.Result"
                            }
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TradeCondition"
                ],
                "summary": "ImpoprtTradeCondition",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/simulate.Result"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TradeCondition"
                ],
                "summary": "DeletaAllResultAndCond",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/condition/latest": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TradeCondition"
                ],
                "summary": "GetLatestTradeCondition",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/simulationcond.AnalyzeCondition"
                            }
                        }
                    }
                }
            }
        },
        "/data/bid-ask": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data"
                ],
                "summary": "ReceiveBidAsk",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/sinopacapi.BidAskProto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/data/streamtick": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Data"
                ],
                "summary": "ReceiveStreamTick",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/sinopacapi.StreamTickProto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/manual/buy-later": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ManualTrade"
                ],
                "summary": "ManualBuyLaterStock",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/manualtradehandler.ManualBuyLaterBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/manual/sell": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ManualTrade"
                ],
                "summary": "ManualSellStock",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/manualtradehandler.ManualSellBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/system/restart": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MainSystem"
                ],
                "summary": "Restart",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/system/sysparm": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MainSystem"
                ],
                "summary": "UpdateSysparm",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/sysparm.Parameters"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/target": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Target"
                ],
                "summary": "GetTarget",
                "parameters": [
                    {
                        "type": "string",
                        "description": "count",
                        "name": "count",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/targethandler.TargetResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/trade-event": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TradeEvent"
                ],
                "summary": "ReciveTradeEvent",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/sinopacapi.EventProto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/trade-record": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TradeRecord"
                ],
                "summary": "UpdateTradeRecord",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/sinopacapi.TradeRecordArrProto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/trade/switch": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MainSystem"
                ],
                "summary": "GetTradeBotSwitch",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/global.SystemSwitch"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MainSystem"
                ],
                "summary": "UpdateTradeBotSwitch",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/global.SystemSwitch"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "balance.Balance": {
            "type": "object",
            "properties": {
                "discount": {
                    "type": "integer"
                },
                "forward": {
                    "type": "integer"
                },
                "original_balance": {
                    "type": "integer"
                },
                "reverse": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                },
                "trade_count": {
                    "type": "integer"
                },
                "trade_day": {
                    "type": "string"
                }
            }
        },
        "global.SystemSwitch": {
            "type": "object",
            "properties": {
                "buy": {
                    "type": "boolean"
                },
                "buy_later": {
                    "type": "boolean"
                },
                "mean_time_reverse_trade_stock_num": {
                    "type": "integer"
                },
                "mean_time_trade_stock_num": {
                    "type": "integer"
                },
                "sell": {
                    "type": "boolean"
                },
                "sell_first": {
                    "type": "boolean"
                },
                "use_bid_ask": {
                    "type": "boolean"
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "attachment": {},
                "response": {
                    "type": "string"
                }
            }
        },
        "manualtradehandler.ManualBuyLaterBody": {
            "type": "object",
            "properties": {
                "price": {
                    "type": "number"
                },
                "stock_num": {
                    "type": "string"
                }
            }
        },
        "manualtradehandler.ManualSellBody": {
            "type": "object",
            "properties": {
                "price": {
                    "type": "number"
                },
                "stock_num": {
                    "type": "string"
                }
            }
        },
        "simulate.Result": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "cond": {
                    "$ref": "#/definitions/simulationcond.AnalyzeCondition"
                },
                "cond_id": {
                    "type": "integer"
                },
                "forward_balance": {
                    "type": "integer"
                },
                "is_best_forward": {
                    "type": "boolean"
                },
                "is_best_reverse": {
                    "type": "boolean"
                },
                "negative_days": {
                    "type": "integer"
                },
                "positive_days": {
                    "type": "integer"
                },
                "reverse_balance": {
                    "type": "integer"
                },
                "total_days": {
                    "type": "integer"
                },
                "total_loss": {
                    "type": "integer"
                },
                "trade_count": {
                    "type": "integer"
                },
                "trade_day": {
                    "type": "string"
                }
            }
        },
        "simulationcond.AnalyzeCondition": {
            "type": "object",
            "properties": {
                "close_change_ratio_high": {
                    "type": "number"
                },
                "close_change_ratio_low": {
                    "type": "number"
                },
                "forward_out_in_ratio": {
                    "type": "number"
                },
                "history_close_count": {
                    "type": "integer"
                },
                "open_change_ratio": {
                    "type": "number"
                },
                "reverse_out_in_ratio": {
                    "type": "number"
                },
                "rsi_high": {
                    "type": "number"
                },
                "rsi_low": {
                    "type": "number"
                },
                "ticks_period_count": {
                    "type": "integer"
                },
                "ticks_period_limit": {
                    "type": "number"
                },
                "ticks_period_threshold": {
                    "type": "number"
                },
                "trim_history_close_count": {
                    "type": "boolean"
                },
                "volume_per_second": {
                    "type": "integer"
                }
            }
        },
        "sinopacapi.BidAskProto": {
            "type": "object",
            "properties": {
                "bid_ask": {
                    "$ref": "#/definitions/sinopacapi.BidAskProto_BidAskData"
                },
                "exchange": {
                    "type": "string"
                }
            }
        },
        "sinopacapi.BidAskProto_BidAskData": {
            "type": "object",
            "properties": {
                "ask_price": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "ask_volume": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "bid_price": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "bid_volume": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "code": {
                    "type": "string"
                },
                "date_time": {
                    "type": "string"
                },
                "diff_ask_vol": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "diff_bid_vol": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "simtrade": {
                    "type": "integer"
                },
                "suspend": {
                    "type": "integer"
                }
            }
        },
        "sinopacapi.EventProto": {
            "type": "object",
            "properties": {
                "event": {
                    "type": "string"
                },
                "event_code": {
                    "type": "integer"
                },
                "info": {
                    "type": "string"
                },
                "resp_code": {
                    "type": "integer"
                }
            }
        },
        "sinopacapi.StreamTickProto": {
            "type": "object",
            "properties": {
                "exchange": {
                    "type": "string"
                },
                "tick": {
                    "$ref": "#/definitions/sinopacapi.StreamTickProto_TickData"
                }
            }
        },
        "sinopacapi.StreamTickProto_TickData": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "ask_side_total_cnt": {
                    "type": "integer"
                },
                "ask_side_total_vol": {
                    "type": "integer"
                },
                "avg_price": {
                    "type": "number"
                },
                "bid_side_total_cnt": {
                    "type": "integer"
                },
                "bid_side_total_vol": {
                    "type": "integer"
                },
                "chg_type": {
                    "type": "integer"
                },
                "close": {
                    "type": "number"
                },
                "code": {
                    "type": "string"
                },
                "date_time": {
                    "type": "string"
                },
                "high": {
                    "type": "number"
                },
                "low": {
                    "type": "number"
                },
                "open": {
                    "type": "number"
                },
                "pct_chg": {
                    "type": "number"
                },
                "price_chg": {
                    "type": "number"
                },
                "simtrade": {
                    "type": "integer"
                },
                "suspend": {
                    "type": "integer"
                },
                "tick_type": {
                    "type": "integer"
                },
                "total_amount": {
                    "type": "number"
                },
                "total_volume": {
                    "type": "integer"
                },
                "volume": {
                    "type": "integer"
                }
            }
        },
        "sinopacapi.TradeRecordArrProto": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sinopacapi.TradeRecordProto"
                    }
                }
            }
        },
        "sinopacapi.TradeRecordProto": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "code": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "order_time": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "sysparm.Parameters": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "targethandler.TargetResponse": {
            "type": "object",
            "properties": {
                "close": {
                    "type": "number"
                },
                "stock_num": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.4.0",
	Host:        "",
	BasePath:    "/trade-bot",
	Schemes:     []string{},
	Title:       "ToC Trader",
	Description: "API docs for ToC Trader",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
