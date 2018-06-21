#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <sys/wait.h>

#include <unistd.h>
#include <signal.h>
#include <pthread.h>
#include <pty.h>
#include <arpa/inet.h>

/* Telnet definitions */
#define T_CMD   255
#define T_DONT  254
#define T_DO    253
#define T_WONT  252
#define T_WILL  251
#define T_ECHO  1
#define T_SGA   3

static unsigned char t_will_sga[]  = { T_CMD, T_WILL, T_SGA  };
static unsigned char t_will_echo[] = { T_CMD, T_WILL, T_ECHO };
static unsigned char t_dont_echo[] = { T_CMD, T_DONT, T_ECHO };

/* Network server */
static volatile int _netfd = -1;
static int _server_socket;

static int
nt_setup_socket(int port)
{
    struct sockaddr_in addr;
    int i;
    int sockfd;

    if ((sockfd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        perror("server: can't open stream socket");
        exit(1);
    }

    /*
     * Set socket option to reuse local address. This is supposed
     * to have the effect of freeing up the local address.
     */
    i = 1;
    if (setsockopt(sockfd, SOL_SOCKET,
                SO_REUSEADDR, (char *) &i, 4) < 0) {
        perror("setsockopt");
    }

    /* Set up server address */
    memset((void *) &addr,0x0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = htonl(INADDR_ANY);
    addr.sin_port = htons(port);

    /* Bind our local address and port */
    if (bind(sockfd, (struct sockaddr *) &addr, sizeof(addr)) < 0) {
        perror("server: can't bind local address");
        exit(1);
    }

    /* Only process one connection at a time */
    listen(sockfd, 1);

    return sockfd;
}

static int 
nt_wait_for_connection(int sockfd)
{
    struct sockaddr_in addr;
    socklen_t len;
    int newsockfd;

    len = sizeof(addr);
    if ((newsockfd = accept(sockfd, (struct sockaddr *) &addr, &len)) < 0) {
        perror("accept");
        exit(1);
    }

    /* Telnet response required for correct interaction with app */
    write(newsockfd, t_will_sga, sizeof(t_will_sga));
    write(newsockfd, t_will_echo, sizeof(t_will_echo));
    write(newsockfd, t_dont_echo, sizeof(t_dont_echo));

    return newsockfd;
}

static void
no_do_telnet(int ifd, int ofd)
{
    unsigned char cmd[8];
    int i;

    /* Read command */
    for (i = 1; i < 3; i++) {
        if (read(ifd, &cmd[i], 1) != 1) {
            return;
        }
    }

    /* Select negative response */
    switch (cmd[1]) {
        case T_WILL:
        case T_WONT:
            cmd[1] = T_DONT;
            break;
        case T_DO:
        case T_DONT:
            cmd[1] = T_WONT;
            break;
        default:
            /* Ignore */
            return;
    }

    /* Ignore responses to our own commands */
    switch (cmd[2]) {
        case T_ECHO:
        case T_SGA:
            return;
    }

    /* Send response */
    cmd[0] = T_CMD;
    write(ofd, &cmd, 3);
}

static void *
nt_net2tty(void *arg)
{
    int fd = *((int *)arg);
    unsigned char data;

    while (1) {
        while (_netfd < 0) {
            usleep(100000);
        }
        if (read(_netfd, &data, 1) != 1) {
            /* Broken pipe -- client quit */
            close(_netfd);
            _netfd = -1;
            raise(SIGUSR2);
        } else if (data == T_CMD) {
            no_do_telnet(_netfd, fd);
        } else {
            write(fd, &data, 1);
            fsync(fd);
        }
    }
    return NULL;
}

static void *
nt_tty2net(void *arg)
{
    int fd = *((int *)arg);
    unsigned char data;

    while (1) {
        if (read(fd, &data, 1) != 1) {
            /* Broken pipe -- app quit */
            close(fd);
            break;
        }
        if (_netfd >= 0) {
            write(_netfd, &data, 1);
            fsync(_netfd);
        }
    }
    return NULL;
}

static int
nt_start_application(char **args, int s)
{
    pid_t pid;

    if ((pid = fork()) < 0) {
        perror ("fork");
        return -1;
    }
    if (pid > 0) {
        return pid;
    }
    dup2(s, 0);
    dup2(s, 1);
    dup2(s, 2);
    execv(*args, args); 
    exit(0); 
}

static void 
nt_restart(void)
{
    int s;

    if ((s = nt_wait_for_connection(_server_socket)) < 0) {
        fprintf(stderr, "connection error\n"); 
        exit(0); 
    }
    /* Reenable the proxies with the new fd */
    _netfd = s; 
}

static void 
nt_signal(int sig)
{
    if (sig == SIGUSR2) {
        nt_restart(); 
    }
}

int 
main(int argc, char *argv[])
{
    int p;
    int pid;
    int do_fork = 0;
    int ttyfd, appfd;
    pthread_t id;

    argv++; 

    if (!*argv || !strcmp(*argv, "--help") || !strcmp(*argv, "-h")) {
        printf("usage: netserve [-d] <port> <program> [args]\n"); 
        exit(0); 
    }

    if (!strcmp(*argv, "-d")) {
        do_fork = 1; 
        argv++; 
    }

    if ((p = atoi(*argv++)) == 0) {
        fprintf(stderr, "bad port number\n"); 
        exit(0); 
    }

    if (do_fork) {
        /* Daemonize */
        pid = fork(); 
        if (pid < 0) {
            perror("fork()"); 
            exit(0);
        }
        if (pid > 0) {
            exit(0); 
        }
        setsid();
    }

    /* Broken pipes are not a problem */
    signal(SIGPIPE, SIG_IGN); 

    /* Get a pseudo tty */
    if (openpty(&ttyfd, &appfd, NULL, NULL, NULL) < 0) {
        perror("open pty"); 
        exit(0); 
    }

    /* Start the application up with sv[1] as its stdio */
    pid = nt_start_application(argv, appfd);

    /* Start proxy for input */
    if (pthread_create(&id, NULL, nt_net2tty, (void *)&ttyfd) < 0) {
        perror("pthread_create"); 
        exit(0); 
    }

    /* Start proxy for output */
    if (pthread_create(&id, NULL, nt_tty2net, (void *)&ttyfd) < 0) {
        perror("pthread_create"); 
        exit(0); 
    }

    /* Setup server */
    _server_socket = nt_setup_socket(p);

    /* SIGUSR2 restarts the server when the client connection closes */
    signal(SIGUSR2, signal);

    /* Start the first server */
    raise(SIGUSR2);

    /* Wait for our child to exit */
    waitpid(pid, &p, 0);

    return 0;
}

/* Example launching TestApplication on port 5611:
 *
 *   netserve -d 1234 TestApplication
 *
 * Now connect to TestApplication, e.g. like this:
 *
 *   telnet 127.0.0.1 1234
 *
 */
