/*****************************************************************************
 * client-c.c                                                                 
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

#define SEND_BUFFER_SIZE 2048


/* TODO: client()
 * Open socket and send message from stdin.
 * Return 0 on success, non-zero on failure
*/
int client(char *server_ip, char *server_port) {

  // necessary variables
  int success, mysock, offset, readsize, s;
  struct addrinfo hints, *res, *a;
  char buf[SEND_BUFFER_SIZE];


  // hints setup
  memset(&hints, 0, sizeof hints);  // empty struct
  hints.ai_family = AF_INET;        // IPv4
  hints.ai_protocol = IPPROTO_TCP;  // TCP??
  hints.ai_socktype = SOCK_STREAM;  // TCP??
  hints.ai_flags = AI_PASSIVE;      // use my ip

  // STEP 1: get info about server ip, throw error if fail
  success = getaddrinfo(server_ip, server_port, &hints, &res);
  if (success != 0) {
    fprintf(stderr, "client: getaddrinfo error: %s\n", gai_strerror(success));
    exit(1);
  }

  // STEP 2 & 3: make a socket and connect to server!
  // loop through res, trying to connect to a valid server
	for(a = res; a != NULL; a = a->ai_next) {
    // make socket, throw error if fail
    mysock = socket(a->ai_family, a->ai_socktype, a->ai_protocol);
		if (mysock == -1) {
			perror("client: socket");
			continue;
		}

    // connect to server, throw error if fail
		if (connect(mysock, a->ai_addr, a->ai_addrlen) == -1) {
			perror("client: connect");
			close(mysock);
			continue;
		}

    // if everything's fine, we can exit the loop!
		break;
	}
  // if a gets to end of linked list, client failed to connect -> throw error
  if (a == NULL) {
    fprintf(stderr, "client: failed to connect\n");
    return 2;
  }

  // STEP 4: read & send message
  while(1) {
    // read file incrementally, resetting buf each time
    // if read 0 bytes, exit loop
    memset(buf, 0, SEND_BUFFER_SIZE);
    readsize = fread(buf, 1, SEND_BUFFER_SIZE, stdin);
    if (readsize == 0) break;
    if (readsize == -1) {
      perror("client: read");
      close(mysock);
      exit(1);
    }

    // send message portion, throw error if fail
    s = send(mysock, buf, readsize, 0);
    // printf("send: %d\n", s);
    if (s == -1) {
      perror("send");
      close(mysock);
      exit(1);
    }
  }

  // STEP 5: close the socket
  close(mysock);

  return 0;
}

/*
 * main()
 * Parse command-line arguments and call client function
*/
int main(int argc, char **argv) {
  char *server_ip;
  char *server_port;

  if (argc != 3) {
    fprintf(stderr, "Usage: ./client-c [server IP] [server port] < [message]\n");
    exit(EXIT_FAILURE);
  }

  server_ip = argv[1];
  server_port = argv[2];
  return client(server_ip, server_port);
}
