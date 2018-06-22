class SocketHandler
  constructor: (sockUrl, messageHandler) ->
    @sock = new WebSocket(sockUrl)
    @handler = messageHandler
    @sock.onopen = (m) ->
    @sock.onmessage = (m) =>
      @handler(m)
    @sock.onerror = (m) ->
    @sock.onclose = (m) ->
  send: (msg) ->
    @sock.send(msg)

window.SocketHandler = SocketHandler
