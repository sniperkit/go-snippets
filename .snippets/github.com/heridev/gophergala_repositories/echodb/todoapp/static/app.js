/*global $ */
/*jshint unused:false */
var app = app || {};
var ENTER_KEY = 13;
var ESC_KEY = 27;

$(function () {
  'use strict';

  // kick things off by creating the `App`
  new app.AppView();

  // if a browser supports websocket we're going to use it :)
  if (window["WebSocket"]) {
    var conn = new WebSocket("ws://" + url + "/ws/todo");

    conn.onclose = function(evt) {
      $('#todoapp h1').css("color", "rgba(255, 255, 255, 0.3)")
    }

    conn.onmessage = function(evt) {
      var data = {};
      // web2.0 blink ;)
      $( "#todoapp h1 span" ).animate({opacity: 0.5}, 500, function() {
        $( "#todoapp h1 span" ).animate({opacity: 1}, 500, function() {
          $( "#todoapp h1 span" ).animate({opacity: 0.8}, 500, function() {
            $( "#todoapp h1 span" ).animate({opacity: 1}, 500, function() {
            })
          })
        })
       });
      try {
        data = JSON.parse(evt.data)
        if(data.__action) {
          data.__doc.id = data.__doc._id;
          if(data.__action === "create") {
            app.todos.add(data.__doc)
          }
          else if(data.__action === "update") {
            var model = app.todos.findWhere({id: data.__doc._id})
            if(model) {
              model.set('title', data.__doc.title)
              model.set('completed', data.__doc.completed)
              model.set('order', data.__doc.order)
            }
          }
          else if(data.__action === "delete") {
            var model = app.todos.findWhere({id: data.__doc._id})
            if(model) {
              model.trigger('destroy', model, model.collection, {});
            }
          }
          else if(data.__action === "destroy" && data.__collection && data.__collection === "todo") {
            app.todos.each(function(m) {
              m.trigger('destroy', m, m.collection, {});
            })

          }
        }
      }
      catch (e) {
      }
    }
    conn.onerror = function(err) {
      $('#todoapp h1').css("color", "rgba(255, 25, 25, 0.3)")
    }
    conn.onopen = function() {
      $('#todoapp h1').css("color", "rgba(25, 255, 25, 0.3)")
    }
  }
});
