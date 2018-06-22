import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    return this.store.createRecord('device-class');
  },

  deactivate: function() {
    var klass = this.get('controller.model');
    if (klass.get('isNew')) {
      klass.transitionTo('loaded.created.uncommitted');
      klass.destroyRecord();
    }
  },

  setupController: function(controller, model) {
    controller.set('errorMessage', null);
    controller.set('model', model);
  },

  actions: {
    save: function(klass) {
      var that = this;
      var $form = Ember.$('#device-class-form');
      if ($form.valid()) {
        klass.save().then(function() {
          that.transitionTo('device-classes');
        }, function(response) {
          that.set('controller.errorMessages', response.errors);
        });
      }
    },

    cancel: function() {
      this.transitionTo('device-classes');
    }
  }
});
