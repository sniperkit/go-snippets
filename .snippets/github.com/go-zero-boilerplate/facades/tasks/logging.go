package tasks

//LogInfo will log an Info message
func LogInfo(getMsg func() string) LoadFunc {
	return func(taskCtx *Context) {
		taskCtx.Logger.Info(getMsg())
	}
}

//WithLoggerField adds Field to the context Logger
func WithLoggerField(getKeyVal func() (key string, value interface{})) LoadFunc {
	return func(taskCtx *Context) {
		taskCtx.Logger = taskCtx.Logger.WithField(getKeyVal())
	}
}

//WithLoggerFields adds Fields to the context Logger
func WithLoggerFields(getFields func() map[string]interface{}) LoadFunc {
	return func(taskCtx *Context) {
		taskCtx.Logger = taskCtx.Logger.WithFields(getFields())
	}
}
