# returns copiable encrypted string (base64) from given args
enc () {
	echo "$*" | openssl enc -base64 -e -aes-256-cbc -nosalt
}

# decrypt result of above encryption (base64, aes-256-cbc, nosalt)
dec () {
	echo "$*" | openssl enc -base64 -d -aes-256-cbc -nosalt
}
