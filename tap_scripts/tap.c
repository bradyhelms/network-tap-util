#include <arpa/inet.h>
#include <netinet/if_ether.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <time.h>
#include <unistd.h>

#define BUFSIZE 65536

// pcap global header structure
struct pcap_header {
  unsigned int magic_number;    // 0xa1b2c3d4
  unsigned short version_major; // 2
  unsigned short version_minor; // 4
  unsigned int thiszone;        // 0
  unsigned int sigfigs;         // 0
  unsigned int snaplen;         // 65535
  unsigned int network;         // 1 (Ethernet)
};

// pcap packet header structure
struct pcap_packet_header {
  unsigned int ts_sec;   // timestamp seconds
  unsigned int ts_usec;  // timestamp microseconds
  unsigned int incl_len; // length of packet data
  unsigned int orig_len; // original length of packet data
};

// Function to write the pcap global header
void write_pcap_header(FILE *file) {
  struct pcap_header header = {
      .magic_number = 0xa1b2c3d4,
      .version_major = 2,
      .version_minor = 4,
      .thiszone = 0,
      .sigfigs = 0,
      .snaplen = 65535,
      .network = 1 // Ethernet
  };
  fwrite(&header, sizeof(header), 1, file);
}

// Function to write a packet's data to the pcap file
void write_packet(FILE *file, unsigned char *data, int length) {
  struct pcap_packet_header packet_header;
  struct timespec ts;
  clock_gettime(CLOCK_REALTIME, &ts);

  // Timestamp: seconds and microseconds
  packet_header.ts_sec = ts.tv_sec;
  packet_header.ts_usec = ts.tv_nsec / 1000; // convert to microseconds
  packet_header.incl_len = length;
  packet_header.orig_len = length;

  fwrite(&packet_header, sizeof(packet_header), 1, file); // Write packet header
  fwrite(data, length, 1, file);                          // Write packet data
}

void print_header(void);

int main(int argc, char **argv) {

  // Parse command line options
  int duration = 0;
  if (argc > 1) {
    duration = atoi(argv[1]);
    if (duration <= 0) {
      fprintf(stderr, "Invalid duration. Please provide a positive integer.\n");
      return 1;
    }
  }

  print_header();

  int sockfd = socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
  if (sockfd < 0) {
    perror("socket");
    return 1;
  } else {
    printf("Socket succesfully opened.\n");
  }

  FILE *pcap_file = fopen("capture.pcap", "wb");
  if (!pcap_file) {
    perror("fopen");
    return 1;
  } else {
    printf("Creating file 'capture.pcap'.\n");
  }

  write_pcap_header(pcap_file); // Write the global pcap header

  unsigned char buffer[BUFSIZE];
  time_t start_time = time(NULL);

  printf("Starting packet capture. Press Ctrl+C to stop capture.");

  while (1) {
    if (duration > 0 && (time(NULL) - start_time) >= duration) {
      break;
    }
    int length = recv(sockfd, buffer, BUFSIZE, 0);
    if (length < 0) {
      perror("recv");
      break;
    }
    write_packet(pcap_file, buffer,
                 length); // Write each packet to the pcap file
  }

  printf("Capture stopped.\n");
  printf("Writing to file.\n");

  close(sockfd);
  fclose(pcap_file);
  return 0;
}

void print_header(void) {
  printf("   _____                         __ _____ \n"
         "  / ____|                       /_ | ____|\n"
         " | |  __ _ __ ___  _   _ _ __    | | |__  \n"
         " | | |_ | '__/ _ \\| | | | '_ \\   | |___ \\ \n"
         " | |__| | | | (_) | |_| | |_) |  | |___) |\n"
         "  \\_____|_|  \\___/ \\__,_| .__/   |_|____/ \n"
         "                        | |               \n"
         "                        |_|               \n");

  printf(
      " _   _      _                      _      _____           \n"
      "| \\ | |    | |                    | |    |_   _|          \n"
      "|  \\| | ___| |___      _____  _ __| | __   | | __ _ _ __  \n"
      "| . ` |/ _ \\ __\\ \\ /\\ / / _ \\| '__| |/ /   | |/ _` | '_ \\ \n"
      "| |\\  |  __/ |_ \\ V  V / (_) | |  |   <    | | (_| | |_) |\n"
      "\\_| \\_/\\___|\\__| \\_/\\_/ \\___/|_|  |_|\\_\\   \\_/\\__,_| .__/ \n"
      "                                                   | |    \n"
      "                                                   |_|    \n");
}
