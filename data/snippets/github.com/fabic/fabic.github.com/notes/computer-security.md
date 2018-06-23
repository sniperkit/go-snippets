---
layout: page
title: "Computers, security (notes)"
tagline: "Personal notes about computers, IT & security."
category : notes
tags : [draft, security]
published: true
---


##### Bash's $RANDOM

At [SO](https://stackoverflow.com/a/1195035) :

```bash
echo $((1 + RANDOM % 10))
```

Someone mentions that doing arithmetics, and specificallly modulus arithmetics
on random numbers easily leads to reduced randomness :

> Be careful here. While this is fine in a pinch, doing arithmetic on random
> numbers can dramatically affect the randomness of your result. in the case of
> $RANDOM % 10, 8 and 9 are measurably (though marginally) less probable than
> 0-7, even if $RANDOM is a robust source of random data. – dimo414 Feb 22 '14
> at 22:51

And :

> By moduloing your random input, you are "pigeon-holing" the results. Since
> $RANDOM's range is 0-32767 the numbers 0-7 map to 3277 different possible
> inputs, but 8 and 9 can only be produced 3276 different ways (because 32768
> and 32769 aren't possible). This is a minor issue for quick hacks, but means
> the result is not uniformly random. Random libraries, like Java's Random,
> offer functions to properly return a uniform random number in the given
> range, rather than simply mod-ing a non-divisible number. – dimo414 May 26
> '14 at 16:07

> @J.F.Sebastian very true - the problem with modulo is it can break the
> uniformity of any RNG, not just bad PRNGs, but thanks for calling this out. –
> dimo414 Sep 24 '14 at 13:41

> Just for context, the basic pigeonholing for % 10 means 8 and 9 are about
> .03% less likely to occur than 0–7. If your shell script requires more
> accurate uniform random numbers than that, then by all means use a more
> complex and proper mechanism. – Nelson Sep 24 '14 at 23:48

Or with a [nice trick](https://stackoverflow.com/a/46320476) :

```bash
echo ${RANDOM:0:1} # random number between 1 and 9
echo ${RANDOM:0:2} # random number between 1 and 99
```

###### Related readings :

* <https://tobtu.com/decryptocat.php> DecryptoCat security troubles, by
  [tobtu.com](https://tobtu.com/)
* [2013: Anatomy of a pseudorandom number generator – visualising Cryptocat’s buggy PRNG (by Paul Ducklin)](https://nakedsecurity.sophos.com/2013/07/09/anatomy-of-a-pseudorandom-number-generator-visualising-cryptocats-buggy-prng/)
  __[MUST READ]__ (mentions Chi-squared statistical technique of assessing non-/random-/ness).

    ![Resulting colourmap of 20 millions old-school Cryptocat "random" floats (derived from PRNG values 0..250)](https://sophosnews.files.wordpress.com/2013/07/colourmap-mod251-4801.png)

##### Alternatives :

With _Coreutils'_ __shuf__ :

```bash
shuf -i 1-100 -n 1
```

From `/dev/urandom` or `/dev/random` (note that __/dev/random__ might stall
(block) in case it runs out of numbers, 'tis better to resort to __urandom__.) :

```bash
od -A n -t d -N 1 /dev/urandom
```

With __dd__ [solution](https://stackoverflow.com/a/32172294) for reading from
__/dev/urandom__ (drawback is it outputs useless unrelated stuff) :

```bash
$ dd if=/dev/urandom count=4 bs=1 | od -t d

4+0 records in
4+0 records out
0000000  1249408195
0000004
4 bytes copied, 4.5596e-05 s, 87.7 kB/s
```

[Nice (albeit less concise) method with AWK](https://stackoverflow.com/a/1197337) :

```bash
awk 'BEGIN {
   # seed
   srand()
   for (i=1;i<=1000;i++){
     print int(1 + rand() * 100)
   }
}'
```

And another nice solution (__fixme:__ it overwrites some `tmp` file which may exist!) :

```bash
( x=$(head -c 1 /dev/urandom > tmp &&
        hexdump -d tmp | head -n 1 | cut -c13-15) &&
    echo $(( 10#$x & 127 )) )
```

[And](https://stackoverflow.com/a/8990953) :

```bash
grep -m1 -ao '[0-9]' /dev/urandom | sed s/0/10/ | head -n1
```

Or [Perl](https://stackoverflow.com/a/23539217) :

```bash
$ perl -e 'print int rand 10, "\n"; '
```

## Pointers


## EOF

