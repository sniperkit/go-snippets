# tesseract & imagemagick must be in path
# usage: `ocr (file/folder)`

OCR_TIFF_FILE=~/ocr_temp.tiff

extract_text () {
    file=$1

    # convert file to tiff using imagemagick
    eval $(convert "${file}" -resize 400% -type Grayscale $OCR_TIFF_FILE)

    # extract text from above tiff file
    eval $(tesseract -l eng $OCR_TIFF_FILE "${file}")
}

ocr () {
    PASSED=$1
    if [ -d "${PASSED}" ] ; then
        # PASSED is a directory

        for file in "$PASSED"/*.{jpg,jpeg,png} ; do
            # don't process if filename begins with "*."
            if [[ $file = *"/*."* ]]; then
                continue
            fi

            eval "extract_text ${file}"
        done
    else
        if [ -f "${PASSED}" ]; then
            # PASSED is a file

            file=$PASSED

            if [[ "$file" =~ \.(jpg|jpeg|png)$ ]]; then
                eval "extract_text ${file}"

                echo "\nTIFF file stored at $OCR_TIFF_FILE"

                echo "Output file: ${file}.txt"
                echo "------------------------------\n"
                
                eval "cat ${file}.txt"
            else
                echo "File not supported. Only jpg, jpeg & png are supported."
            fi

        else
            echo "${PASSED} is not valid.";
        fi
    fi
}

