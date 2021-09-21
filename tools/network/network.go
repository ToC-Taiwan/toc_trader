// Package network package network
package network

import (
	"net"
	"time"

	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// CheckPortIsOpen CheckPortIsOpen
func CheckPortIsOpen(host string, port string) bool {
	logger.Logger.Infof("Checking host on %s:%s...", host, port)
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		logger.Logger.Warn(err)
	}
	if conn != nil {
		defer func() {
			if err := conn.Close(); err != nil {
				logger.Logger.Error(err)
			}
		}()
		return true
	}
	return false
}
