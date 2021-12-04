// Package restclient package restclient
package restclient

import (
	"github.com/go-resty/resty/v2"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
)

var client *resty.Client

// GetClient GetClient
func GetClient() *resty.Client {
	if client != nil {
		return client
	}
	new := resty.New()
	new.SetLogger(logger.GetLogger())
	client = new
	return client
}
