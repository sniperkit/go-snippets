import DS from 'ember-data';

var DeviceEvent = DS.Model.extend({
  local_time: DS.attr('number'),
  latlng: DS.attr(),
   created_at: DS.attr('date'),
});

//{"id":"54c542d612fa742501000019",
  //"data":{"cpu":0.9244193434715271,"height":2.8205637487124195},
  //"local_time":1,
  //"latlng":[19.42705,-99.16619],
  //"created_at":"2015-01-25T14:24:06.348-05:00"},


DeviceEvent.reopenClass({
  FIXTURES: [
  {
    id: 1,
    altitude: 33.3,
    direction: 10
  },
  {
    id: 2,
    altitude: 33.3,
    direction: 20
  },
  {
    id: 3,
    altitude: 12.5,
    direction: 30
  },
  {
    id: 4,
    altitude: 45,
    direction: 100
  },
  {
    id: 5,
    altitude: 122.4,
    direction: 2
  },
  {
    id: 6,
    altitude: 0,
    direction: -3
  },
  ]
});


export default DeviceEvent;
