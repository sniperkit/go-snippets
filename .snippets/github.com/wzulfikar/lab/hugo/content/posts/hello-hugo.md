---
title: "Hello Hugo!"
date: 2018-03-03T02:33:20+08:00
tags: ["experiment", "go"]
draft: false
disableToc: true
---

<p class="text-center">Another experiment on how/where I write things I like to write <i>(Duh)</i>.</p>

![gohugo](/images/hugo-logo.png#featured) 

<p class="text-center red">***</p>


So, I've been through the path of Wordpress, Tumblr, Ghost, and recently, Gitbook. Each has its own merits and drawbacks. However, I see something interesting from Hugo, that I didn't see in my previous adventures with other writing tools. 

<!--more-->

Before we go to Hugo, I'd like to clarify the context by telling more about what I've found in my previous writing tools, and what I'm looking for. 

## Wordpress

Wordpress (https://wordpress.org) has been around since [2003](https://en.wikipedia.org/wiki/WordPress) and can be seen as one of the most famous blogging tools, or CMS (Content Management System). And by being most famous, it means that people can find help and support easily. 

It has a good [built-in visual editor](/images/wordpress-editor.jpg) and so many [themes](https://wordpress.org/themes/) to choose from. Many things can be integrated using ready-made plugins (ie. google analytics, image gallery, e-commerce, etc.). *Good things.* 

To get started with Wordpress is not so much difficult. Most hosting providers have built-in Wordpress package. Technical wise, an online Wordpress instance will need a hosting, database (MySQL), and domain. Since Wordpress stores its contents (posts, settings, etc) in database, doing a full backup means having a copy of files in the hosting and data stored in database.

I haven't really had a bad experience when using Wordpress. However, the visual editor seems too slow for me to work with, and an internet connection is required to access the editor. Well, I can just write somewhere and paste to Wordpress editor when I've internet connection, but I prefer not to.

With that downside, I started to look for something else.

## Tumblr

In [Tumblr](https://www.tumblr.com), things seem to be less complicated than Wordpress. You don't need to have your own hosting or to install your own database. The editor is good, getting started is just about creating new account. Your posts are hosted in Tumblr at no cost. 

Actually, I kind of felt comfortable using Tumblr. Not so long after, I realize that I've to be connected (internet) to write. Besides that, using Tumblr means that I don't actually have full-control of my blog to customize. For those who like to customize their blog, they may not like it.

Long story short, I started another search for writing tools.. 

## Ghost

Ghost is a blogging platform released on [2013](https://en.wikipedia.org/wiki/Ghost_(blogging_platform)). What I like the most about Ghost is its [clean, and user-friendly interface](https://user-images.githubusercontent.com/120485/28764244-344050c0-75d5-11e7-9314-45bc4177164e.png). Aside from that, it's actually an open source project (see: https://github.com/TryGhost/Ghost). 

To get started, you can install Ghost on your own server ([self-host](https://docs.ghost.org/v1/docs/getting-started-guide)), or subscribe to their managed service: [Ghost(Pro)](https://ghost.org). In my case, I used self-host version of Ghost. Installation is not difficult, things went well as expected. 

You may want to note that Ghost is coded in Node.js (a programming language) that's not as common as PHP (used to build Wordpress) when it comes to shared hosting. Hence, installing Ghost may not be available in most shared web hosting providers.

Overall, I don't really have a complain for Ghost. However, I actually still looking for something that I can put my writing offline, and post it later when I've internet connection. Also, I'd prefer to not put my contents in database, so it still can be seen anywhere.

## Gitbook

[Gitbook](https://www.gitbook.com) is kind of interesting. It's a platform to write book and publish content online. It has an offline editor where you can start writing and post later. In Gitbook, your contents are stored as files, in your computer. There's no need for database.

To get started with Gitbook, you can create an account there (https://gitbook.com) and download [Gitbook Editor](https://www.gitbook.com/editor). Open your Gitbook Editor, sign in with your account, and you can start writing. Try to explore Gitbook features from online resources and you'll find some interesting stuffs, like the ability to download your contents as PDF or ePub. 

There's something off with Gitbook tho. It's meant for publishing book. While what I'm looking for is something to publish blog-like contents, which have time-related contents. With that reasoning, I feel like there must be a more suitable tool for this.

## Hello Hugo! 

Hugo (http://gohugo.io) takes different approach compared to other tools we've discussed. It's a static site generator, something like [Jekyll](https://jekyllrb.com). 

Static site generator basically means that Hugo will takes your contents and build a static site from it. Since it's static, you can host your site at no cost in platforms like [GitHub](https://github.com), [Gitlab](https://gitlab.com), [Netlify](https://www.netlify.com), [surge.sh](https://surge.sh), etc.

Since we mentioned Jekyll, how does Hugo differ from Jekyll? 

Jekyll is built using Ruby (a programming language), which means that you'll need Ruby in your machine to use Jekyll. Yep, it's something *techy*. Hugo is built using Go (another programming language) that can produce a binary for multiple platform. People who developed Hugo have published it as a binary (a self executable program) that anyone can download. Hence, to use Hugo, one doesn't need to have Go programming language installed in their machine. They just need to have the Hugo binary. See more: [Hugo Quick Start](http://gohugo.io/getting-started/quick-start/).

> Hugo released its [v0.16](https://github.com/gohugoio/hugo/releases?after=v0.17) – the first Hugo release that wasn't tagged as `pre-release`, on June 2016.

With Hugo, all my contents are stored as files (markdown files, to be specific). There's no need for database. I can write my contents offline and publish it later. I can see my files anywhere, and I can use any editor I like (VIM, VSCode, SublimeText, etc).

If you wonder what markdown is, it's a form of formatting similar to how you [format message in Whatsapp](https://faq.whatsapp.com/en/android/26000002/). Try writing markdown yourself here: http://markdownlivepreview.com.

While I personally convinced to use Hugo (this blog itself is built using Hugo), it might not be for everyone since it involves some *techy* steps in its flow of writing contents. Nevertheless, I hope that knowing Hugo can somehow benefits you.

Yup, you've reached the end of this post! And to end this, there's also a good writing platform that you may want to check out: [Medium](https://medium.com). 

***Till next. See ya!***

![Gopher mascot](/images/gopher-head.png#featured "Gopher mascot")
<p class="text-center">**The Go gopher: an iconic mascot of Go project.*</p>

