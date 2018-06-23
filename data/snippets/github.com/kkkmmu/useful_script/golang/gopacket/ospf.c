#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <ifaddrs.h>
#include <sys/socket.h>
#include <netdb.h>

#define SUCCESS 0
#define FAILURE -1

#include <netinet/in.h>
#define MCAST_ALL_SPF_ROUTERS "224.0.0.5"

static int join_router_4_multicast(int fd)
{
	int ret = SUCCESS;
	int one = 1;
	struct ip_mreqn mreqn;

	/* FIXME: we does not set any of imr_address or
	 * 	 * imr_ifindex because 0 works fine for us now.
	 * 	 	 * This should be fixed with a proper routine */
	memset(&mreqn, 0, sizeof(struct ip_mreqn));

	ret = inet_aton(MCAST_ALL_SPF_ROUTERS, &mreqn.imr_multiaddr);
	if (ret < 0) {
		printf("failed to convert router multicast address");
		return FAILURE;
	}

	/*
	   ret = setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &one , sizeof(int));
	   if (ret < 0) {
	   printf("failed to add multicast membership");
	   return FAILURE;
	   }
	*/

	ret = setsockopt(fd, IPPROTO_IP, IP_ADD_MEMBERSHIP, &mreqn, sizeof(mreqn));
	if (ret < 0) {
		printf("failed to add multicast membership");
		return FAILURE;
	}

	return ret;
}

int main() {
	int fd, on = 1, ret, flags;
	char buf[1024] = {0};

	fd = socket(PF_INET, SOCK_RAW, 89);
	if (fd < 0) {
		printf("Failed to create raw IPPROTO_OSPF socket");
		return FAILURE;
	}

	/*
	ret = setsockopt(fd, IPPROTO_IP, IP_HDRINCL, &on, sizeof(on));
	if (ret < 0) {
		printf("Failed to set IP_HDRINCL socket option");
		return FAILURE;
	}

	ret = setsockopt(fd, IPPROTO_IP, IP_PKTINFO, &on, sizeof(on));
	if (ret < 0) {
		printf("Failed to set IP_PKTINFO socket option");
		return FAILURE;
	}
	*/

	ret = join_router_4_multicast(fd);
	if (ret != 0) {
		printf("cannot join multicast groups");
		return FAILURE;
	}

	while (recvfrom(fd, buf, 1024, 0, NULL, NULL) > 0) {
		printf("Received a packet\n");
	}

	return 0;
}
