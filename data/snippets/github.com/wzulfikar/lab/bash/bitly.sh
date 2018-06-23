# shorten url using bitly via 
# npm module `bitly` (npm i -g bitly)
bitly () {
	GREEN='\033[0;32m'

	# reject if url is empty
	if [ -z "$1" ]; then
	    echo "Invalid command - please provide a link to shorten.\nExample: bitly https://google.com"
	    return
	fi

	# sample result:
	# "You can access your short url at http://cnet.co/2ozs55t"
	result=$(eval "/usr/local/bin/bitly -u $1")

	link=${result##* } # only get last word of result
	echo $GREEN$link
}