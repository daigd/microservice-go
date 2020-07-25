package config

import (
	"log"
	"os"

	kitlog "github.com/go-kit/kit/log"
)

var Logger *log.Logger
var KitLogger kitlog.Logger

func init() {
	// 初始化go默认日志组件 及 go-kit 的日志组件
	Logger = log.New(os.Stderr, "", log.LstdFlags)
	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "my-ts", kitlog.DefaultTimestampUTC)
	KitLogger = kitlog.With(KitLogger, "my-caller", kitlog.DefaultCaller)
}
