// Uses encoding/json.Marshal to return JSON to the client.
func RenderLowerJson(o interface{}) revel.Result {
	return RenderLowerJsonResult{o, true}
}

type RenderLowerJsonResult struct {
	obj            interface{}
	lowerFirstChar bool
}

func (r RenderLowerJsonResult) Apply(req *revel.Request, resp *revel.Response) {
	var b []byte
	var err error
	if revel.Config.BoolDefault("results.pretty", false) {
		b, err = jsonutils.MarshalIndent(r.obj, "", " ", r.lowerFirstChar)
	} else {
		b, err = jsonutils.Marshal(r.obj, r.lowerFirstChar)
	}

	if err != nil {
		revel.ErrorResult{Error: err}.Apply(req, resp)
		return
	}

	resp.WriteHeader(200, "application/json")
	resp.Out.Write(b)
}
