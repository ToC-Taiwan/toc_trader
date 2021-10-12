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
        "/data/bid-ask": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tradebot"
                ],
                "summary": "ReceiveBidAsk",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/bidask.BidAskProto"
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
                    "tradebot"
                ],
                "summary": "ReceiveStreamTick",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/streamtick.StreamTickProto"
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
        "/data/target": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tradebot"
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
                                "$ref": "#/definitions/tradebothandler.TargetResponse"
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
        "/system/pyserver/host": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "mainsystem"
                ],
                "summary": "UpdatePyServerHost",
                "parameters": [
                    {
                        "type": "string",
                        "description": "host",
                        "name": "py_host",
                        "in": "header",
                        "required": true
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
                    "mainsystem"
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
                    "sysparm"
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
        "/system/trade/switch": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "mainsystem"
                ],
                "summary": "GetTradeBotCondition",
                "responses": {
                    "200": {
                        "description": ""
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
                    "mainsystem"
                ],
                "summary": "UpdateTradeBotCondition",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/mainsystemhandler.UpdateTradeBotConditionBody"
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
        "/trade-event": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tradeevent"
                ],
                "summary": "ReciveTradeEvent",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tradeevent.EventProto"
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
                    "traderecord"
                ],
                "summary": "UpdateTradeRecord",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/traderecord.TradeRecordArrProto"
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
        "/trade/manual/buy_later": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tradebot"
                ],
                "summary": "ManualBuyLaterStock",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tradebothandler.ManualBuyLaterBody"
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
        "/trade/manual/sell": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tradebot"
                ],
                "summary": "ManualSellStock",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tradebothandler.ManualSellBody"
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
        "bidask.BidAskProto": {
            "type": "object",
            "properties": {
                "bid_ask": {
                    "$ref": "#/definitions/bidask.BidAskProto_BidAskData"
                },
                "exchange": {
                    "type": "string"
                }
            }
        },
        "bidask.BidAskProto_BidAskData": {
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
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "attachment": {
                    "type": "object"
                },
                "response": {
                    "type": "string"
                }
            }
        },
        "mainsystemhandler.UpdateTradeBotConditionBody": {
            "type": "object",
            "properties": {
                "enable_buy": {
                    "type": "boolean"
                },
                "enable_buy_later": {
                    "type": "boolean"
                },
                "enable_sell": {
                    "type": "boolean"
                },
                "enable_sell_first": {
                    "type": "boolean"
                },
                "mean_time_trade_stock_num": {
                    "type": "integer"
                },
                "use_bid_ask": {
                    "type": "boolean"
                }
            }
        },
        "streamtick.StreamTickProto": {
            "type": "object",
            "properties": {
                "exchange": {
                    "type": "string"
                },
                "tick": {
                    "$ref": "#/definitions/streamtick.StreamTickProto_TickData"
                }
            }
        },
        "streamtick.StreamTickProto_TickData": {
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
        "tradebothandler.ManualBuyLaterBody": {
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
        "tradebothandler.ManualSellBody": {
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
        "tradebothandler.TargetResponse": {
            "type": "object",
            "properties": {
                "close": {
                    "type": "number"
                },
                "stock_num": {
                    "type": "string"
                }
            }
        },
        "tradeevent.EventProto": {
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
        "traderecord.TradeRecordArrProto": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/traderecord.TradeRecordProto"
                    }
                }
            }
        },
        "traderecord.TradeRecordProto": {
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
	Version:     "0.1.0",
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
	swag.Register(swag.Name, &s{})
}
