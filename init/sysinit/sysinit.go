// Package sysinit all init should put here
package sysinit

import (
	_ "gitlab.tocraw.com/root/toc_trader/init/sysparminit" // sysparminit

	_ "gitlab.tocraw.com/root/toc_trader/init/dbinit" // dbinit

	_ "gitlab.tocraw.com/root/toc_trader/init/globalinit" // globalinit

	_ "gitlab.tocraw.com/root/toc_trader/init/taskinit" // taskinit
)
