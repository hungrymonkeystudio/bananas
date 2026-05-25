// properly find the resource path

package paths

import (
	"github.com/WarrenWu4/bananatype/internal/foundation/logger"
)

var (
	Build = "dev"
)

func GetResourcePath() string {
	// FIX: not entirely reliable since where the pkg manager
	// installs the machine is varied
	if Build == "prod" {
		logger.Log(logger.DEBUG, "Using prod resource path: /usr/share/bananatype")
		return "/usr/share/bananatype"
	} else {
		logger.Log(logger.DEBUG, "Using dev resource path: ./resources")
		return "./resources"
	}
}
