###############################################################################
# server-python.py
# Name:  Monique Legaspi
# NetId: mlegaspi
###############################################################################

import sys
import socket

RECV_BUFFER_SIZE = 2048
QUEUE_LENGTH = 10

def server(server_port):
    """TODO: Listen on socket and print received message to sys.stdout"""

    # STEP 1: get info about server ip
    hints = socket.getaddrinfo(None, server_port, socket.AF_INET, socket.SOCK_STREAM, 0, socket.AI_PASSIVE)

    # socket/bind loop
    for serv in hints:
        # STEP 2: make socket
        try:
            s = socket.socket(serv[0], serv[1])
        except:
            continue
        
        # STEP 3: bind
        try:
            s.bind(serv[4])
        except:
            continue

        # if everything's fine, exit loop
        break

    # STEP 4: listen for incoming connections
    try:
        s.listen(QUEUE_LENGTH)
    except:
        exit

    # STEP 5: wait for connections, accept when we get one
    # while 1:
    # try:
    conn, addr = s.accept()
    # except:
    #     continue

    # STEP 6: receive message from client
    while 1:
        data = conn.recv(RECV_BUFFER_SIZE)
        if not data: break
        print(data)

    # STEP 7: close new socket
    conn.close()

    pass


def main():
    """Parse command-line argument and call server function """
    if len(sys.argv) != 2:
        sys.exit("Usage: python server-python.py [Server Port]")
    server_port = int(sys.argv[1])
    server(server_port)

if __name__ == "__main__":
    main()
