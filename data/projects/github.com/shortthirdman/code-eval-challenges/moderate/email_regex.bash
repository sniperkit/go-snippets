email_regex="^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@\w*(\w+\.)+\w{2,4}$"

while read -r line || [ -n "$line" ]; do
    if [[ "$line" =~ $email_regex ]]; then
        echo true
    else
        echo false
    fi
done <$1
