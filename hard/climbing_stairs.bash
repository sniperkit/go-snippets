while read line || [ -n "$line" ]; do
  if [[ -n "$line" ]]; then
    BC_LINE_LENGTH=0 bc <<< 'a=0;b=1;c=1;for(i=0;i<'$line';i+=1){a=b;b=c;c=a+b};b'
  fi
done <$1
