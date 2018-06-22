<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="Content-Style-Type" content="text/css" />
  <meta http-equiv="content-script-type" content="text/javascript" />

  <link href="//maxcdn.bootstrapcdn.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
  <link rel="stylesheet" href="/static/style.css" type="text/css">
  {% if title %}
  <title>kindle my clippings - {{ title }}</title>
  {% else %}
  <title>kindle my clippings</title>
  {% endif %}
</head>
<body>
<div id="all">
  <div id="inner-content">

    <div id="header">
      <h1 class="title"><a href="/">kindle my clippings</a></h1>
    </div>

    <div id="content">
    {% for clip in clips %}
    <div class="entry-content">
        <div class="book-title"><a href="/book/{{ clip.Title }}">{{ clip.Title }}</a>
        &nbsp;&nbsp;&nbsp;<a class="amazon-link" href="http://www.amazon.co.jp/s/?url=search-alias%3Ddigital-text&field-keywords={{clip.Title}}"><i class="fa fa-external-link fa-lg"></i></a></div>
        <div class="author">by {{ clip.Author }}</div>
        <div class="added-date">Added on {{ clip.AddedOn|date:"Mon, 02 Jan 2006 15:04" }}</div>
        <blockquote class="clip">{{ clip.Content }}</blockquote>
    </div>
    {% endfor %}
    </div>
  </div>

  <div id="footer">
    <div class="container clearfix"></div>

    <div class="footer-bottom">
      Powered by <a href="http://golang.org/">Go</a>
    </div>
  </div>

</div>
</body></html>
