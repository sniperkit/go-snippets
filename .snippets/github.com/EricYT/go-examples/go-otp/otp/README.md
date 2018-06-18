# Erlang/OTP simple implement by Go.
---
I don't think it's a right way to do this again.
Erlang is a dynamic type language, so functions
dispatch is easy and simple to understand, but
it is hard to implement by go, and Go has no
pattern match. Reflect is a wonderful tool for
go, but it has drawback I don't know how to
use it to make Go flexible like Erlang. Maybe
you have a answer for that, please tell me!
A C/S architecture for Erlang is
a nature way. Every process acts as a box has a
special ID, we conmunicate with it by this one.
Go has same thing named 'channel'.
