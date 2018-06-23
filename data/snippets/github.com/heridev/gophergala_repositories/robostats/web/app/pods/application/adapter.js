import DS from 'ember-data';
import Ember from 'ember';

export default DS.ActiveModelAdapter.extend({
  ajaxError: function(xhr) {
    var error = this._super(xhr);

    if (xhr && xhr.status === 422) {
      var jsonErrors = Ember.$.parseJSON(xhr.responseText)["errors"];
      return new DS.InvalidError(jsonErrors);
    } else {
      return error;
    }
  }
});
