import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    return this.store.find('deviceSession');
  },

  actions: {
    showEvents: function(session) {
      this.transitionTo('device-sessions.show', session);
    }
  }
});
