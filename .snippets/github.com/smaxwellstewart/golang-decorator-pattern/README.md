# golang-decorator-pattern

This is an example of the decorator pattern in golang. It is very useful when
a method has many orthogonal concerns.

What does 'many orthogonal concerns' mean? Say we wanna do a http request, pretty simple right.
But what if we also need to perform some logging of the request as well as add the requested action to an audit trail, notify another service of this request, increase the request count in a rate limiter, etc. None of these concerns are related, other than that we wanna do all of them whenever we call our method.
This is where the decorator pattern really helps.

By implementing the decorator pattern we subscribe to the open-closed principal.
Our method is open to future extension but closed to future modification. There's a lot of groovy benefits to obeying the open-closed pricniple.

For more useful info see: https://www.youtube.com/watch?v=xyDkyFjzFVc
