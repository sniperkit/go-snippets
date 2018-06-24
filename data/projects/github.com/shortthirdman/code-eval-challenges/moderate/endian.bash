if [ `echo -n I|hexdump -o|awk '{print substr($2,6,1); exit}'` -eq 0 ]; then
    echo "BigEndian"
else
    echo "LittleEndian"
fi
