package logger

import insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"

var lvlString = map[insightV1.SeverityLevel]string{
	insightV1.SeverityLevel_DEBUG:   "debu",
	insightV1.SeverityLevel_INFO:    "info",
	insightV1.SeverityLevel_WARNING: "warn",
	insightV1.SeverityLevel_ERROR:   "fail",
}

func level(lvl insightV1.SeverityLevel) string {
	return lvlString[lvl]
}
