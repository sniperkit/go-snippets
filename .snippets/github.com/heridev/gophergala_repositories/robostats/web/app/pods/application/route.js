import Ember from 'ember';
import ApplicationRouteMixin from 'simple-auth/mixins/application-route-mixin';

export default Ember.Route.extend(ApplicationRouteMixin, {
  actions: {
    error: function(error/*, transition*/) {
      var status = error.status;
      if (error && (status >=  500 && status < 599)) {
        alert("Ha ocurrido un error. Por favor recarga la página e intenta de nuevo");
      }
    }
  }
});
