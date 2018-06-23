import Ember from 'ember';

export default Ember.Route.extend({
  model: function(params) {
    var that = this;

    // if deviceClass is already in memory, the next request is not performed
    return this.store.find('deviceClass', params.device_class_id).then(function(klass) {
      return that.store.find('deviceInstance', {device_class_id: klass.get('id')});
    });
  },

  actions: {
    showSessions: function(device) {
      this.transitionTo('device-sessions.show', device);
    }
  }
});
