# extract links from given url
# by `curl`-ing the page.
getlinks () {
	# use curl's -L switch to follow redirect
	link_regex='href="([^"#]+)"'
	links=$(curl -L $1 2>&1 | grep -o -E "$link_regex" | cut -d'"' -f2)

	# display links found
	echo $links

	links_found=$(echo $links | wc -l)
	# display number of links found
	echo "\nFound:$links_found link(s)."
}