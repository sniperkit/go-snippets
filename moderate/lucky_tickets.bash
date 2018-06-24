while read line || [ -n "$line" ]; do
  if [[ -n "$line" ]]; then
    BC_LINE_LENGTH=0 bc <<< 'define f(x) {if (x <= 1) return (1); return (f(x-1)*x)} define c(n,k) {return (f(n)/(f(n-k)*f(k)))};m='$line';r=0;for (i=0;i<m/2;i++) {r+=-1^i*c(m,i)*c(11*m/2-1-10*i,m-1)};r'
  fi
done <$1
