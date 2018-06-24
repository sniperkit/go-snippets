(define (oddnum n)
  (if (<= n 99)
    (begin
      (display n)
      (newline)
      (oddnum (+ n 2)))))
(oddnum 1)
