# global io

'use strict'

angular.module 'feedBag'
  .factory 'socket', (socketFactory) ->

    # socket.io now auto-configures its connection when we omit a connection url
    ioSocket = io 'http://localhost:3000',
      # Send auth token on connection, you will need to DI the Auth service above
      # 'query': 'token=' + Auth.getToken()
      path: '/socket.io'

    socket = socketFactory ioSocket: ioSocket

    socket: socket

    ###
    Register listeners to sync an array with updates on a model

    Takes the array we want to sync, the model name that socket updates are sent from,
    and an optional callback function after new items are updated.

    @param {String} modelName
    @param {Array} array
    @param {Function} callback
    ###
    syncUpdates: (modelName, array, callback) ->

      ###
      Syncs item creation/updates on 'model'
      ###
      socket.on modelName, (items) ->
        _.each items, (item) -> array.push(item)
        callback? array

    ###
    Removes listeners for a models updates on the socket
    @param modelName
    ###
    unsyncUpdates: (modelName) ->
      socket.removeAllListeners modelName
