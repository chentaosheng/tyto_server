package daemonutil

import (
	"github.com/sevlyar/go-daemon"
	"tyto/core/logging"
)

// 变为守护进程
func Start(logger logging.Logger, name string, code string) (*daemon.Context, bool) {
	logger.Error("start daemon failed, system unsupported:", name, code)
	return nil, false
}

func Stop(logger logging.Logger, dctx *daemon.Context) bool {
	logger.Error("stop daemon failed, system unsupported")
	return false
}
