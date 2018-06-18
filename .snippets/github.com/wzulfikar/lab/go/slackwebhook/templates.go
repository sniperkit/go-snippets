package slackwebhook

// TODO:
// - redirect
// - google recaptcha

func templateFormSubmitted(redirect string) string {
	return `
<html>
Form submitted!
<br/>
<a href="` + redirect + `">← Back to site</a>
</html>
`
}

func templateFormError(err, site string) string {
	backToSite := ""
	if site != "" {
		backToSite = `<br/><a href="` + site + `">← Back to site</a>`
	}
	return `
<html>
<h4>Oops! Something went wrong :(</h4>
<br/>
` + err + `
` + backToSite + `
</html>
`
}
