import Ember from 'ember';

export default Ember.View.extend({
  didInsertElement: function() {
    var series = this.get('controller.model.time_serie');
    var cpu = series.values['cpu'];
    cpu.unshift('CPU');
    
    var chart = c3.generate({
      bindto: '#chart',
      data: {
        columns: [
          cpu
        ]
      }
    });
  }
});
