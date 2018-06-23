import Ember from 'ember';
import Session from 'simple-auth/session';
import Authenticator from 'simple-auth-oauth2/authenticators/oauth2';

var CustomSession = Session.extend({
  user: function() {
    var userId = this.get('userId');
    if (!Ember.isEmpty(userId)) {
      return this.container.lookup('store:main').find('user', userId);
    }
  }.property('userId'),
});

var CustomAuthenticator = Authenticator.extend({
  authenticate: function(creds) {
    return new Ember.RSVP.Promise(function(resolve, reject) {
      Ember.$.ajax({
        url: '/user/login',
        type: 'POST',
        data: {
          grant_type: 'password',
          email: creds.identification,
          password: creds.password
        }
      }).then(function(response) {
        Ember.run(function() {
          resolve({
            access_token: response.access_token,
            userId: response.user_id,
          });
        });
      }, function(xhr/*, status, err*/) {
        Ember.run(function() {
          reject(xhr.responseText);
        });
      });
    });
  }
});

export function initialize(container/*, application */) {
  container.register('session:custom', CustomSession);
  container.register('authenticator:custom', CustomAuthenticator);
}

export default {
  name: 'authentication',
  before: 'simple-auth',
  initialize: initialize
};


