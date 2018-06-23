import DS from 'ember-data';

var DeviceSession = DS.Model.extend({
  class_id: DS.attr('string'),
  created_at: DS.attr('date'),
  start_time: DS.attr('date'),
  end_time: DS.attr('date'),
  user_id: DS.attr('string'),
  session_key: DS.attr('string'),
  instance_id: DS.attr('string'),
  events: DS.hasMany('deviceEvent', {async: true})
});

DeviceSession.reopenClass({
  FIXTURES: [
  {
    id: 1,
    start_at: new Date(),
    end_at: new Date(),
    events: ["1", "2", "3", "4"]
  },
  {
    id: 2,
    start_at: new Date(),
    end_at: new Date(),
    events: ["5", "6"] 
  }
  ]
});

export default DeviceSession;
