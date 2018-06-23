package toolkit

type TryCatch struct {
	try     func()
	catch   func(interface{})
	finally func()
}

func Try(f func()) *TryCatch {
	o := TryCatch{}
	o.try = f
	return &o
}

func (o *TryCatch) Catch(f func(interface{})) *TryCatch {
	o.catch = f
	return o
}

func (o *TryCatch) Finally(f func()) *TryCatch {
	o.finally = f
	return o
}

func (o *TryCatch) Run() {
	if o.finally != nil {
		defer o.finally()
	}
	if o.catch != nil {
		defer func() {
			if r := recover(); r != nil {
				o.catch(r)
			}
		}()
	}
	o.try()
}
