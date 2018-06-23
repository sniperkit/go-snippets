# upload file to https://transfer.sh
# source: https://github.com/dutchcoders/transfer.sh/
upload () {
    # write to output to tmpfile because of progress bar
    tmpfile=$( mktemp -t transferXXX )
    filename=$(basename $1)

    # support optional arg, will return link to virustotal scan result.
    # eg. upload malicious-file.exe virustotal
    optional_arg=""
    if [ $# -eq 2 ]
      then
        optional_arg="/$2"
    fi

    curl --progress-bar --upload-file $1 https://transfer.sh/$filename$optional_arg >> $tmpfile;
    cat $tmpfile;
    rm -f $tmpfile;
}
