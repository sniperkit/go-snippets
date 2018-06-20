package owl

import (
	f "github.com/crosserclaws/intest/common/db/facade"
	log "github.com/crosserclaws/intest/common/logruslog"
)

var DbFacade *f.DbFacade
var logger = log.NewDefaultLogger("warn")
