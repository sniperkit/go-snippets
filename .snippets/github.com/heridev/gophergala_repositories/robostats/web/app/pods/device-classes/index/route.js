import Ember from 'ember';
import AuthenticatedRouteMixin from 'simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
  model: function() {
    return this.store.find('device-class');
  },

  actions: {
    showDevices: function(klass) {
      this.transitionTo('device-classes.show', klass);
    },

    delete: function(klass) {
      klass.destroyRecord();
    }
  }
});
