import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    // get id from the session

    var dataURL = 'http://api.dev.robostats.io/device_sessions/time_series?session_id=54c542d612fa74250100000f&key[]=cpu';
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
