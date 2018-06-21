echo "foo:123:bar:789" | awk -F: '{print $0}'
echo "foo:123:bar:789" | awk -F: '{print $1}'
echo "foo:123:bar:789" | awk -F: '{print $2}'
echo "foo:123:bar:789" | awk -F: '{print $3}'
echo "foo:123:bar:789" | awk -F: '{print $4}'
echo "foo:123:bar:789" | awk -F: '{print NF}'
echo "foo:123:bar:789" | awk -F: '{print $NF}'
echo "foo:123:bar:789" | awk -F: '{print $1, $NF}'
echo "foo:123:bar:789" | awk -F: '{print $(NF-1)}'


echo "foo;123;bar;789" | awk -F: '{print $(NF-1)}'
