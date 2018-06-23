import DS from 'ember-data';

export default DS.Model.extend({
  user_id: DS.attr('string'),
  class_id: DS.attr('string'),
  created_at: DS.attr('date'),
  sessions: DS.hasMany('deviceSession', {async: true})
});


