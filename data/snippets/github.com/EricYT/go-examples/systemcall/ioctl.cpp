#include <iostream>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <sys/ioctl.h>

using namespace std;

int main(void) {
    int fd = ::open("/dev/sda", O_RDONLY | O_DIRECT);
    if(fd != -1) {
        uint8_t buf[512] = { 0xEC, 0, 0, 1, 0 };
        int ret = ::ioctl(fd, 0x31F, buf);

        for(int i = 0; i < 512; ++i) {
            //::putchar(buf[i]);
            printf("%0x", buf[i]);
        }
    } else {
        ::perror("open");
    }

    ::close(fd);
    return 0;
}
