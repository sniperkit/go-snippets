type RoundTripper func(r *http.Request) (*http.Response, error)

func (rt RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

/*...*/

c := http.Client{
	Transport: RoundTripper(func(r *http.Request) (*http.Response, error) {
		/* ... */
	}),
}
