import DS from 'ember-data';

var DeviceClass = DS.Model.extend({
    name: DS.attr('string'),
    user: DS.belongsTo('user'),
    api_key: DS.attr('string'),
    devices: DS.hasMany('deviceInstance', {async: true})
});

DeviceClass.reopenClass({
  FIXTURES: [
  {
    id: 1,
    name: "Drones",
    api_key: "abcd",
    devices: ["1", "2"]
  },
  {
    id: 2,
    name: "Fridges",
    api_key: "abcde",
    devices: ["3", "4", "5"]
  },
  {
    id: 3,
    name: "Television",
    api_key: "abcdef",
    devices: ["6", "7"]
  },
  ]
});

export default DeviceClass;
