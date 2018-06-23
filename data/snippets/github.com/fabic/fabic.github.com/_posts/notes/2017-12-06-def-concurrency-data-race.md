---
layout: note
title: "Data race (definition)"
tagline: ""
category : notes
tags : [note, memory, concurrency]
---

A __data race__ is a particular type of race condition in which these three
behaviors occur [[Rust](https://doc.rust-lang.org/book/second-edition/ch04-02-references-and-borrowing.html)]:

* __Two or more pointers__ access the same data at the same time.
* At least one of the pointers is being used to __write__ to the data.
* There’s no mechanism being used to __synchronize access__ to the data.
