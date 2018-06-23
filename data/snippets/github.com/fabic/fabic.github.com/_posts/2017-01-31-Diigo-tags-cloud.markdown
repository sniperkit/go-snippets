---
layout: post
title: "Diigo bookmarks: tags cloud"
description: "Just my online bookmarks."
category: diary
published: true
tags: [links, dark]
---

<!-- Moved Diigo tags cloud to a dedicated post.
     See /index.html where I tried to have it be loaded asynchronously
     upon scroll into view, but _they_ resort do `document.write(...)`
     which is prohibeted for external scripts (that are added asynchronously). -->

<div class="diigo-tags" style="width:97%;">
  <div class="diigo-banner sidebar-title center"
       style="font: bold 12px arial;margin-bottom:5px;">
    <a href="https://www.diigo.com/cloud/fabicxx"
       style="font-variant: small-caps;">
      <img src="https://www.diigo.com/images/ii_blue.gif"
           width="16" height="16"
           alt="Diigo icon"/>&nbsp;Diigo Bookmark tags</a>
  </div>
  <!-- Note: Pass QS arg. `count=128` so as to restrict the number of displayed tags.-->
  <script type="text/javascript"
          src="https://www.diigo.com/tools/tagrolls_script/fabicxx?sort=freq;icon;size=11-23;color=56bbff-167bc2;title=Diigo%20Bookmark%20tags;name;showadd;v=3"></script>
</div>
