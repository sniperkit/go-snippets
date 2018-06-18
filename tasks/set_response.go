package tasks

//SetResponse will set response on the context
func SetResponse(getResponse func() interface{}) LoadFunc {
	//the reason it is a func (getResponse) and not a normal arg is because:
	// if we call SetResponse(dtos.NewStores(storeList)) in the Chain, the DTO is created before storeList is populated
	return func(taskCtx *Context) {
		taskCtx.Response = getResponse()
	}
}
