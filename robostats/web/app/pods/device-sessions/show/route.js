import Ember from 'ember';

export default Ember.Route.extend({
  model: function(params) {
    return this.promise(params.device_session_id);
  },

  setupController: function(controller, model) {
    controller.set('model', model);
    if (!model.time_serie) {
      window.location.reload();
    }
  },

  promise: function(id) {
    var dataURL = 'http://api.dev.robostats.io/device_sessions/time_series?session_id='+id+'&key[]=cpu';
    var token = this.get('session.content.access_token');

    return Ember.$.ajax({
      type: 'GET',
      dataType: "json",
      url: dataURL,
      headers: {
        "Authorization": "Bearer "+token
      }
    });
  }
});
