package nqm

import (
	f "github.com/crosserclaws/intest/common/db/facade"
	log "github.com/crosserclaws/intest/common/logruslog"
	tb "github.com/crosserclaws/intest/common/textbuilder"
)

var DbFacade *f.DbFacade

var t = tb.Dsl

var logger = log.NewDefaultLogger("warn")
