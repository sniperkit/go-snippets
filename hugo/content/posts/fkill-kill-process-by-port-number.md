---
title: "Fabulosly Kill (fkill) Process by Port Number"
date: 2018-03-26T12:49:55+08:00
tags: ["devops"]
draft: true
disableToc: true
---

![fkill logo](/images/bettercap/fkill-logo.jpg#featured)

`TODO: write content` 

https://www.npmjs.com/package/fkill-cli

<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">You can now search and kill processes by port number with the latest `fkill-cli` version, thanks to <a href="https://twitter.com/kevvayo?ref_src=twsrc%5Etfw">@kevvayo</a>. <a href="https://t.co/yzVkFWoh2c">https://t.co/yzVkFWoh2c</a> <a href="https://t.co/fzc8I4ZNEF">pic.twitter.com/fzc8I4ZNEF</a></p>&mdash; Sindre Sorhus (@sindresorhus) <a href="https://twitter.com/sindresorhus/status/974343649600323584?ref_src=twsrc%5Etfw">March 15, 2018</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>


<p class="text-center">***</p>

*Outline:*

1. interactive search: search by name or port
2. fkill: fabulously kill
3. unix way: 
    - find process by port: `lsof -i:8080`
    - kill the process: `kill $(lsof -t -i:8080)` (use -t to only get the process id)
