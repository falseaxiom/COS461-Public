/*****************************************************************************
 * server-c.c                                                                 
 * Name:  Monique Legaspi
 * NetId: mlegaspi
 *****************************************************************************/

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netdb.h>
#include <netinet/in.h>
#include <errno.h>

#define QUEUE_LENGTH 10
#define RECV_BUFFER_SIZE 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 * Return 0 on success, non-zero on failure
*/
int server(char *server_port) {
  // necessary variables
  int success, mysock, newsock, recsize, b, l;
  struct addrinfo hints, *res, *a;
  struct sockaddr_storage client_addr;
  socklen_t addr_len;
	char buf[RECV_BUFFER_SIZE];

  // hints setup
  memset(&hints, 0, sizeof hints);  // empty struct
  hints.ai_family = AF_INET;        // IPv4
  hints.ai_protocol = IPPROTO_TCP;  // TCP??
  hints.ai_socktype = SOCK_STREAM;  // TCP??
  hints.ai_flags = AI_PASSIVE;      // use my ip

  // STEP 1: get info about server ip, throw error if fail
  success = getaddrinfo(NULL, server_port, &hints, &res);
  if (success != 0) {
    fprintf(stderr, "server: getaddrinfo error: %s\n", gai_strerror(success));
		return 1;
  }

  // STEP 2 & 3: make socket and bind!
  // loop through res, trying to bind to valid port
  for (a = res; a != NULL; a = a->ai_next) {
    // make socket, throw error if fail
    mysock = socket(a->ai_family, a->ai_socktype, a->ai_protocol);
    if (mysock == -1) {
      perror("server: socket");
      continue;
    }

    // bind to socket, throw error if fail
    b = bind(mysock, a->ai_addr, a->ai_addrlen);
    // printf("bind: %d\n", b);
    if (b == -1) {
      perror("server: bind");
      close(mysock);
      continue;
    }

    // if everything's fine, we can exit the loop!
    break;
  }
  // if a gets to end of linked list, server failed to bind -> throw error
  if (a == NULL) {
    fprintf(stderr, "server: failed to bind\n");
    exit(1);
  }

  // STEP 4: listen for incoming connections, throw error if fail
  l = listen(mysock, QUEUE_LENGTH);
  // printf("listen: %d\n", l);
  if (l == -1) {
    perror("listen");
    close(mysock);
    exit(1);
  }

  // STEP 5: wait for connections, accept when we get one
  while(1) {
    addr_len = sizeof(client_addr);
    newsock = accept(mysock, (struct sockaddr *)&client_addr, &addr_len);
    // printf("accept: %d\n", newsock);
    if (newsock == -1) {
      perror("accept");
      continue;
    }

    // STEP 6: receive message from client, throw error if fail
    while(1) {
      recsize = recv(newsock, buf, RECV_BUFFER_SIZE, 0);
      // printf("receive: %d\n", recsize);
      if (recsize == 0) break; // size == 0 this means the message is over!
      if (recsize == -1) {
        perror("recv");
        close(newsock);
        break;
      }

      // print message
      write(1, buf, recsize);
      fflush(stdout);
    }

    // STEP 7: close the new socket
    close(newsock);
  }

  return 0;
}

/*
 * main():
 * Parse command-line arguments and call server function
*/
int main(int argc, char **argv) {
  char *server_port;

  if (argc != 2) {
    fprintf(stderr, "Usage: ./server-c [server port]\n");
    exit(EXIT_FAILURE);
  }

  server_port = argv[1];
  return server(server_port);
}
