// plugin_test.go
package stp

import (
	asicdmock "utils/asicdClient/mock"
)

func UsedForTestOnlySetupAsicDPlugin() {

	// Use MOCK plugin
	SetAsicDPlugin(&asicdmock.MockAsicdClientMgr{})
}
