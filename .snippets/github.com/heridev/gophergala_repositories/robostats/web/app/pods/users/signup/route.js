import Ember from 'ember';

export default Ember.Route.extend({
  model: function() {
    return this.store.createRecord('user');
  },

  deactivate: function() {
    var user = this.get('controller.model');
    if (user.get('isNew')) {
      user.transitionTo('loaded.created.uncommited');
      user.deleteRecord();
    }
  },

  setupController: function(controller, model) {
    controller.set('model', model);
    controller.set('errorMessages', null);
  },

  actions: {
    save: function (){
      var user = this.get('controller.model');
      var $form = Ember.$('#user-signup-form');
      if ($form.valid()) {
        user.save().then(function(){
          console.log('user was saved');
        }, function() {
          console.log('errored out');
        });
      }
    },

    goHome: function() {
      this.transitionTo('index');
    }
  }
});
