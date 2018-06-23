/* 

http://pegjs.majda.cz/online 

This @is:a nik. #WHo2 ^is:hi ~this?   What! and by  2/12/12 2 now 

*/

start
 = sentences
 
sentences
 = sentence+
 
sentence
 = words '.' _
 / words '?' _
 / words '!' _
 / words _
 
words
 = (word _)+
 / word

word 
 = goal
 / user
 / tag
 / context
 / date
 / chars digits
 / chars
 / digits
 
goal
 = '^' word ':' word
 / '^' word


user
 = '@' word ':' role
 / '@' word


role
 = 'r'
 / 'a'
 / 'c'
 / 'i'

tag
 = '#' word

context
 = '~' word

date
 = 'by' _  digit '/'  digit '/' digit digit
 / 'by' _  digit digit '/'  digit '/' digit digit
 / 'by' _  digit '/'  digit digit '/' digit digit
 / 'by' _  digit digit '/'  digit digit '/' digit digit
 / 'by' _  digit '/'  digit '/' digit digit digit digit
 / 'by' _  digit digit '/'  digit '/' digit digit digit digit
 / 'by' _  digit '/'  digit digit '/' digit digit digit digit
 / 'by' _  digit digit '/'  digit digit '/' digit digit digit digit



/* ===== Lexical Elements ===== */
 
 
chars
  = chars:char+ { return chars.join(""); }
 
char
  = [A-Za-z]
 
digits
  = digit+
 
digit
  = [0-9]
 
 
/* ===== Whitespace ===== */
 
_ "whitespace"
  = whitespace*
 
whitespace
  = [ \t\n\r] 
