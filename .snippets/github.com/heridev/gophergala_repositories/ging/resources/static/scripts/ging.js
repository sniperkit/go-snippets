
"use strict";

// Query

var queryElement = document.getElementById("query");
var submitElement = document.getElementById("submitQuery");
var resultsElement = document.getElementById("result-wrapper");
var cache = {};
var autoQuery;
var autoQueryDelay = 400;

if (submitElement && queryElement) {
  submitElement.updateState = function() {
    var value = queryElement.value;
    submitElement.disabled = value.length < 3;
  };
  submitElement.updateState();

  var conn = new WebSocket("ws://" + window.location.host + "/stream/query");
  conn.onmessage = function (e) {
    var data = JSON.parse(e.data);
    cache[data.query] = data.result;
    resultsElement.innerHTML = data.result;
  };

  queryElement.addEventListener("input", function(e) {
    clearTimeout(autoQuery);
    autoQuery = setTimeout(function() {
      var queryString = queryElement.value.trim();
      if (queryString.length < 3) {
        resultsElement.innerHTML = "";
        return;
      }
      if (cache[queryString]) {
        resultsElement.innerHTML = cache[queryString];
        return;
      }
      conn.send(queryString);
    }, autoQueryDelay);

    submitElement.updateState();
  }, false);
}
