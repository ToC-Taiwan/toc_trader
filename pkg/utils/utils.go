// Package utils is common utils
package utils

import (
	"math"
	"net"
	"strconv"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
)

// StrToInt64 StrToInt64
func StrToInt64(input string) (ans int64, err error) {
	ans, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		return ans, err
	}
	return ans, err
}

// StrToFloat64 StrToFloat64
func StrToFloat64(input string) (ans float64, err error) {
	ans, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return ans, err
	}
	return ans, err
}

// Round Round
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// MicroDateTimeToTimeStamp MicroDateTimeToTimeStamp
func MicroDateTimeToTimeStamp(dateTime string) (timeStamp int64, err error) {
	dataTimeStamp, err := time.ParseInLocation(global.LongTimeLayout, dateTime[:19], time.Local)
	if err != nil {
		return timeStamp, err
	}
	microSec, err := StrToInt64(dateTime[20:])
	if err != nil {
		return timeStamp, err
	}
	return dataTimeStamp.UnixNano() + microSec*1000, err
}

var checkLock sync.Mutex

// CheckPortIsOpen CheckPortIsOpen
func CheckPortIsOpen(host string, port string) bool {
	defer checkLock.Unlock()
	checkLock.Lock()
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

// GetHostIP GetHostIP
func GetHostIP() string {
	var results []string
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.GetLogger().Error(err)
	}
	var addrs []net.Addr
	for _, i := range ifaces {
		if i.HardwareAddr.String() == "" {
			continue
		}
		addrs, err = i.Addrs()
		if err != nil {
			logger.GetLogger().Error(err)
		}
		for _, addr := range addrs {
			if ip := addr.(*net.IPNet).IP.To4(); ip != nil {
				if ip[0] != 127 && ip[0] != 169 {
					results = append(results, ip.String())
				}
			}
		}
	}
	return results[len(results)-1]
}
