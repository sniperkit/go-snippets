run-asm () {
    FILENAME=$1
    OUTPUT_FILE=asm_tmp.o
    nasm -f macho64 -o $OUTPUT_FILE $1 && ld $OUTPUT_FILE -o $OUTPUT_FILE && ./$OUTPUT_FILE
}