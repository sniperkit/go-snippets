import eventlet


def handle(client):
    while True:
        c = client.recv(1)
        if not c: break
        client.sendall(c)


server = eventlet.listen(('0.0.0.0', 6000))
pool = eventlet.GreenPool(10000)
while True:
    new_sock, address = server.accept()
    pool.spawn_n(handle, new_sock)
