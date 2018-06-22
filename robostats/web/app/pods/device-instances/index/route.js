import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    return this.store.find('deviceInstance');
  },

  actions: {
    showSessions: function(device) {
      this.transitionTo('device-instances.show', device);
    }
  }

});
