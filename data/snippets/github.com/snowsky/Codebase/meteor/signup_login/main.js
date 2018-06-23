user = new Mongo.Collection("user");

if (Meteor.isClient) {
  // counter starts at 0
  Session.setDefault('counter', 0);

  Template.body.helpers({
    users: function () {
      return user.find();
    }
  });

  Template.body.events({
/*    "submit #login-form" : function (e, t) {
      e.preventDefault();
      var email = t.find('#login-email').value
        , password = t.find('#login-password').value;
        
      //
      
      Meteor.loginWithPassword(email, password, function (err) {
        if (err) {
          console.log("ok");          
        } else {
          console.log("fail");          
        }
      });
      return false;
    }
*/

    "submit #signupForm" : function (event) {
      var name = event.target.name.value
        , email = event.target.email.value
        , password = event.target.password.value;
      user.insert({
        name: name,
        email: email,
        password: password
      });
      return false;
    }
  });
}

if (Meteor.isServer) {
  Meteor.startup(function () {
    // code to run on server at startup
  });
}
