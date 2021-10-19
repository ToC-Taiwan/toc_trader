// Package restyinit package restyinit
package restyinit

import (
	"github.com/go-resty/resty/v2"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

func init() {
	client := resty.New()
	client.SetLogger(logger.GetLogger())
	global.RestyClient = client
}
