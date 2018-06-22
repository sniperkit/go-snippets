---
layout: post
title: "wget --header='Cookie:...'"
description: ""
category: diary
tags: [old]
---

[wget](http://www.gnu.org/software/wget/manual/wget.html#HTTP-Options) and
Cookie HTTP header for grabbing pages from a web site that needs e.g. login :

    wget --header='Cookie: PHPSESSID=2at5uflg6anjr30hun9g2k6f02' <url>
