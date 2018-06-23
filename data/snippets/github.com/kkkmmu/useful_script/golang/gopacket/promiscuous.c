#include <stdio.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>
#include <linux/if_ether.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <string.h>

#define ETH_NAME "ens33"

int main(int argc, char **argv) {
	int fd, rv;
	struct ifreq ifr;

	fd=socket(PF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
	if(fd  <0){ 
		printf("Cannot open socket: %s\n", strerror(errno));
		return -1;
	}
	strcpy(ifr.ifr_name, ETH_NAME);

	rv = ioctl(fd, SIOCGIFFLAGS, &ifr);
	if (rv <0){
		close(fd);
		printf("Cannot Get Interface flags: %s\n", strerror(errno));
		return-1;
	}

	if(ifr.ifr_flags & IFF_RUNNING){
		printf("Interface is link up\n");
	}else{
		printf("Interface link down\n");
	}

	ifr.ifr_flags |= IFF_PROMISC;
	rv = ioctl(fd, SIOCSIFFLAGS, &ifr);
	if (rv < 0){
		printf("Cannot Set Interface flags: %s \n", strerror(errno));
		return -1;
	}
	printf("Setting interface ::: %s ::: to promisc\n\n", ifr.ifr_name);

	char buf[1024] = {0};
	while (recvfrom(fd, buf, 1024, 0, NULL, NULL) > 0) {
		printf("Received a packet: %s\n", strerror(errno));
	}

	return 0;
}
