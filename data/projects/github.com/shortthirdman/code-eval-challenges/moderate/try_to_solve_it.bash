tr "abcdefghijklmnopqrstuvwxyz" "uvwxyznopqrstghijklmabcdef" <$1|while read -r line || [ -n "$line" ]; do
  echo "$line"
done
