package connect

type (
	//PushConnectPool is used as description push mode's connect pool
	PushConnectPool struct {
		url  string
		pool []*pushConnect
	}
)

// CreatePushConnectPool is used as instance push connect pool
func CreatePushConnectPool(url string) *PushConnectPool {
	pcp := &PushConnectPool{
		url:  url,
		pool: make([]*pushConnect, defaultPoolSize),
	}

}
