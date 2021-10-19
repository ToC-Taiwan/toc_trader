// Package network package network
package network

import (
	"net"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

var lock sync.Mutex

// CheckPortIsOpen CheckPortIsOpen
func CheckPortIsOpen(host string, port string) bool {
	defer lock.Unlock()
	lock.Lock()
	logger.GetLogger().Infof("Checking host on %s:%s...", host, port)
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		logger.GetLogger().Warn(err)
	}
	if conn != nil {
		defer func() {
			if err := conn.Close(); err != nil {
				logger.GetLogger().Error(err)
			}
		}()
		return true
	}
	return false
}
