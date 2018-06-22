import {
  moduleFor,
  test
} from 'ember-qunit';

moduleFor('route:device-events/index', 'DeviceEventsIndexRoute', {
  // Specify the other units that are required for this test.
  // needs: ['controller:foo']
});

test('it exists', function() {
  var route = this.subject();
  ok(route);
});
