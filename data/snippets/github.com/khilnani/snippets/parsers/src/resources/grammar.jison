/* 
http://zaach.github.io/jison/try/usf/index.html 
SLR(1) or LR(1)

http://zaach.github.io/jison/try/

This is a nik. WHo2 is hi  this?   What! and by    2 now
This @is:a nik. #WHo2 ^is:hi ~this?   What! and by  2/12/12 2 now.

*/

/* lexical grammar */
%lex
%%

// [a-zA-Z0-9]+	      return 'WORDY'
([0-9]|[0-9][0-9])\/([0-9]|[0-9][0-9])\/([0-9][0-9]|[0-9][0-9][0-9][0-9])       return 'DATE'
[0-9]                 return 'DIGIT'
[0-9]+                return 'DIGITS'
[0-9]+("."[0-9]+)?\b  return 'FLOAT'
'by'                  return 'BY'
[A-Za-z]+[0-9]?       return 'STRING'
'.'                   return 'DOT'
'?'                   return 'QUESTION'
'!'                   return 'EXCLAMATION'
'^'                   return 'CARET'
'@'                   return 'AMP'
'#'                   return 'HASH'
'~'                   return 'TILDE'
'/'                   return 'FSLASH'
':r'                  return 'R'
':a'                  return 'A'
':c'                  return 'C'
':i'                  return 'I'
':'                   return 'COLON'
[ \t\n\r]+            return 'WS'
<<EOF>>               return 'EOF'
// (.|\n)             return 'CHAR'

/lex



%start expressions

%% /* language grammar */

expressions
 : sentences EOF
 ;

sentences
 : sentence
 | sentences sentence
 ;

sentence
 : words DOT
 | words DOT WS
 | words QUESTION
 | words QUESTION WS
 | words EXCLAMATION
 | words EXCLAMATION WS
 | words
 ;

words
 : word
 | (word WS)+
 ;

word
 : goal
 | user
 | tag
 | context
 | datestring
 | STRING
 | DIGIT
 | DIGITS
 | FLOAT
 ;

goal
 : CARET word COLON word
 ;

user
 : AMP word R
 | AMP word A
 | AMP word C
 | AMP word I
 ;

tag
 : HASH word
 ;

context
 : TILDE word
 ;

datestring
 : BY WS DIGIT WS
 | BY WS DIGITS WS
 | BY WS DATE 
 ;


