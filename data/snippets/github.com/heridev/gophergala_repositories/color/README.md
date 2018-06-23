color
=====

Syntax highlighting for streamed terminal Go text.


```
cat file.go | color
```

or 

```
cat $(find . -type f -name *.go) | color | less -r

```

![screen shot][1]

I started off wanting to write a very limited in scope Go app with the above
functionality for my first Go application only to find that it's pretty much
all been done by https://github.com/koron/beni. Except that beni only works on
files passed as arguments not the STDIN.

Clearly my contribution is not a big deal and I'll probably just do a quick PR
on beni to get this functionality added.  But, it did give me a chance to work on
something Go related.  Hopefully next time I'll think of an idea that hasn't
already been done by someone else.

[1]:./color.png
