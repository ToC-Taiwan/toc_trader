// Package rest package rest
package rest

import (
	"github.com/go-resty/resty/v2"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
)

// Client Client
var Client *resty.Client

// GetClient GetClient
func GetClient() *resty.Client {
	if Client != nil {
		return Client
	}
	client := resty.New()
	client.SetLogger(logger.GetLogger())
	Client = client
	return Client
}
