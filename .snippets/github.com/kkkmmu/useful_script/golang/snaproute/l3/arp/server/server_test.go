package server

import (
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"log/syslog"
	"testing"
	"time"
	"utils/logging"
)

func NewLogger(name string, tag string, listenToConfig bool) (*logging.Writer, error) {
	var err error
	srLogger := new(logging.Writer)
	srLogger.MyComponentName = name

	srLogger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err != nil {
		fmt.Println("Failed to initialize syslog - ", err)
		return srLogger, err
	}

	srLogger.MyLogLevel = sysdCommonDefs.INFO
	fmt.Println("Logging level ", srLogger.MyLogLevel, " set for ", srLogger.MyComponentName)
	return srLogger, err
}

func TestInitArpParams(t *testing.T) {
	t.Log("Testing initArpParams()")
	logger, err := NewLogger("arpdTest", "ARPTest", true)
	if err != nil {
		t.Errorf("Unable to initialize logger")
		return
	}
	ser := NewARPServer(logger)
	ser.initArpParams()
	if ser.snapshotLen != 65549 || ser.promiscuous != false ||
		ser.minCnt != 1 || ser.retryCnt != 5 ||
		ser.pcapTimeout != time.Duration(1)*time.Second ||
		ser.timerGranularity != 1 || ser.ConfRefreshTimeout != 600 ||
		ser.MinRefreshTimeout != 300 || ser.timeout != time.Duration(1)*time.Second ||
		ser.timeoutCounter != 600 || ser.retryCnt != 5 ||
		ser.minCnt != 1 || ser.probeWait != 5 ||
		ser.probeNum != 5 || ser.probeMax != 20 ||
		ser.probeMin != 10 || ser.dumpArpTable != false ||
		ser.arpSliceRefreshDuration != time.Duration(10)*time.Minute {
		t.Errorf("Arp parameters are not initialized properly")
	} else {
		t.Log("Successfully tested initArpParams()")
	}
}
