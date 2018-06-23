package queue

// Recover prevents handler from propagating panic by calling `catch` on it
func Recover(catch Handler) Middleware {
	return func(next Handler) Handler {
		return func(val interface{}) (err error) {
			defer func() {
				if p := recover(); p != nil {
					err = catch(p)
				}
			}()
			return next(val)
		}
	}
}
