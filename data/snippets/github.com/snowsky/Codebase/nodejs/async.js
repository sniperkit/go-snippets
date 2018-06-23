/*
  async.whilst(
    function () { 
      if(pagetoken == undefined) return;
    },
    function (cb) {
      pagetoken = getMessages(auth, userId, label, pagetoken);
    },
    function (err) {
      console.log("Error get messages: " + err);
        // 5 seconds have passed
    }
  );
*/

/*
          async.series([
            function (callback) {
              getMessages(auth, userId, label, "");
              callback();
            },
            function (callback) {
              getAllMessages(auth, userId, label);
              callback();
            }
          ], function(err, results){
            console.log('Error get all messages: ' + err);
    // results is now equal to ['one', 'two']
          }); 
*/
