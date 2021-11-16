###############################################################################
# client-python.py
# Name:  Monique Legaspi
# NetId: mlegaspi
###############################################################################

import sys
import socket

SEND_BUFFER_SIZE = 2048

def client(server_ip, server_port):
    """TODO: Open socket and send message from sys.stdin"""

    # STEP 1: get info about server ip
    hints = socket.getaddrinfo(server_ip, server_port, socket.AF_INET, socket.SOCK_STREAM, 0, socket.AI_PASSIVE)
    
    # socket/connect loop
    for serv in hints:
        # STEP 2: make socket
        try:
            s = socket.socket(serv[0], serv[1])
        except:
            continue

        # STEP 3: connect to server
        try:
            s.connect((server_ip, server_port))
        except:
            continue

        #if everything's fine, exit loop
        break

    # STEP 4: read & send message
    while 1:
        chunk = sys.stdin.read(SEND_BUFFER_SIZE)
        if (chunk == ""):
            break
        s.send(chunk)

    # STEP 5: close the socket
    s.close()

    pass


def main():
    """Parse command-line arguments and call client function """
    if len(sys.argv) != 3:
        sys.exit("Usage: python client-python.py [Server IP] [Server Port] < [message]")
    server_ip = sys.argv[1]
    server_port = int(sys.argv[2])
    client(server_ip, server_port)

if __name__ == "__main__":
    main()
