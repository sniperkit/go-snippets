import Ember from 'ember';
import LoginControllerMixin from 'simple-auth/mixins/login-controller-mixin';

export default Ember.Controller.extend(LoginControllerMixin, {
  authenticator: 'authenticator:custom',

  actions: {
    authenticate: function() {
      var that = this;
      that._super().then(null, function() {
        that.set('errorMessage', 'Incorrect username or password');
      });
    }
  }
});
