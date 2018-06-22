package log

import (
	"testing"
)

func Test_logrus(t *testing.T) {
	// 打印堆栈信息
	Logger().Debug("This is a debug log.")
	Logger().Info("A group of walrus emerges from the ocean")
	Logger().Warn("The group's number increased tremendously!")
	Logger().Error("The group's number increased tremendously!")
	Logger().Warnf("The group's %s increased tremendously!", "number")
}
