$(function () {
  initWebsockets();

  $('.main-wall').error(function() {
    // Show them the default image if the canvas doesn't exist yet
    $(this).attr('src', 'https://s3.amazonaws.com/gollage/placeholder.png');
  });

  $('.mini-wall').error(function() {
    // Show them the default image if the canvas doesn't exist yet
    $(this).attr('src', 'https://s3.amazonaws.com/gollage/thumbholder.png');
  });

  $('.image-form input[type=file]').on('change', function() {
    // Make the submit button clickable
    $(this).parents('.image-form').find('input[type=submit]').removeAttr('disabled');
  });

  // If links exists
  if (typeof links !== "undefined" && links.length > 0) {
    var holder = $('.wall-holder');
    for (var i = 0; i < links.length; i++) {
      var clickBox = $('<a class="click-box"></a>');
      clickBox.css({
        width: links[i].DispWidth,
        height: links[i].DispHeight,
        top: links[i].YOffset,
        left: links[i].XOffset,
      });
      clickBox.attr('href', "http://" + links[i].Url);
      holder.append(clickBox);
    }
  }
});

