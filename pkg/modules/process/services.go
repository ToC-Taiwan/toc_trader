// Package process main process
package process

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// RestartService RestartService
func RestartService() {
	global.ExitChannel <- "exit"
}
