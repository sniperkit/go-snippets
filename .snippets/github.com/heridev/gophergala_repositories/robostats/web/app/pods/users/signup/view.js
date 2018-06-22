import Ember from 'ember';

export default Ember.View.extend({
  willDestroyElement: function() {
    Ember.$('body').removeClass('login');
  },

  didInsertElement: function() {
    Ember.$('body').addClass('login');
  }
});
