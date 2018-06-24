while read line || [ -n "$line" ]; do
  declare -a a
  v=true
  for ((i=0; i<${#line}; i++)); do
    b="${line:$i:1}"
    case "$b" in
    ')') if [ ${#a[*]} -gt 0 ] && [ ${a[-1]} == '(' ]; then
           unset a[${#a[*]}-1]
         else
           v=false
           break
         fi
         ;;
    ']') if [ ${#a[*]} -gt 0 ] && [ ${a[-1]} == '[' ]; then
           unset a[${#a[*]}-1]
         else
           v=false
           break
         fi
         ;;
    '}') if [ ${#a[*]} -gt 0 ] && [ ${a[-1]} == '{' ]; then
           unset a[${#a[*]}-1]
         else
           v=false
           break
         fi
         ;;
    *)   a[${#a[*]}]=$b
         ;;
    esac
  done
  if $v && [ ${#a[*]} -eq 0 ]; then
    echo "True"
  else
    echo "False"
  fi
  unset a
done <$1
