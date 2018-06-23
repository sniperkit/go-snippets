package ctrl

type Controller struct {
	Page     uint `param:"<in:query>"`
	PageSize uint `param:"<in:query>"`

	Total uint
}

func (t *Controller) GetPage() uint {
	if t.Page == 0 {
		t.Page = 1
	}
	return t.Page
}

func (t *Controller) GetPageSize() uint {
	if t.PageSize == 0 {
		t.PageSize = 10
	}
	if t.PageSize > 500 {
		t.PageSize = 10
	}
	return t.PageSize
}
