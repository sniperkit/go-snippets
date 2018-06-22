  window.users = {
    0: ""
  };
  $(function() {
    username = $.cookie("username");
    window.ws = new WebSocket("ws://localhost:8000/events/{{.Id}}?username="+username);

    window.ws.onmessage = function(e) {
      var tokens = e.data.split(":");
      var command=tokens[0];
      var params=tokens[1];

      switch(command) {
        case "w":
          var tokens = params.split("#");
          var user = parseInt(tokens[0], 10);
          var key = tokens[1];
          window.GO_RACER_NS.carretMoveLogic(key, user)
          break;
        case "users":
          var tokens = params.split("&");
          $(tokens).each(function(index, val){
            var tokensInside = val.split("#");
            var id = parseInt(tokensInside[0],10);
            var name = tokensInside[1];
            window.users[id] = name;
          })
          //TODO: Re-Render Users
          break;
        case "start":

          //TODO: Write start logic
          var template = $('#template');
          for(var id in window.users) {
            var pre = $("<pre/>").html(template.html()).prop("id", 'code_'+id).addClass('code_field');
            if(id != window.current_user_id){
              pre.addClass('another_users_field');
            }
            $(document.body).append(pre);
          }
          var js = $("<script/>").attr('src', "https://google-code-prettify.googlecode.com/svn/loader/run_prettify.js?callback=js_ident")
          $(document.body).append(js);
          break;
        case "current_user":
          window.current_user_id = parseInt(params, 10);
          break;
      }

      $(document.body).append(e.data) ;
      // window.GO_RACER_NS.carretMoveLogic(e.data);
      // $(document.body).append(e.data)
      // window.GO_RACER_NS.carretMoveLogic(e.data, 'self');
    };
    // setTimeout(function(){
    //   window.GO_RACER_NS.user1();
    // }, 2000);
    // window.GO_RACER_NS.user2();
  });

  window.exports = {};
  window.exports['js_ident'] = window.GO_RACER_NS.prepareGameField;