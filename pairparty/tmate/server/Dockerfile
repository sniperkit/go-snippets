FROM fedora

# Set correct environment variables.
ENV HOME /root

# Regenerate SSH host keys. baseimage-docker does not contain any, so you
# have to do that yourself. You may also comment out this instruction; the
# init system will auto-generate one during boot.
#RUN /etc/my_init.d/00_regen_ssh_host_keys.sh

# Use baseimage-docker's init system.
#CMD ["/sbin/my_init"]

RUN yum clean all
RUN yum makecache
RUN yum update -y 
RUN yum -y group install 'C Development Tools and Libraries'
RUN yum -y install git-core libtool libevent-devel ncurses-devel zlib-devel automake libssh2-devel cmake ruby
RUN rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN mkdir -p /opt/code/tmate-slave

COPY . /opt/code/tmate-slave
WORKDIR /opt/code/tmate-slave

RUN ./create_keys.sh && \
		./autogen.sh && \
		./configure && \
		make

RUN ./message.sh

CMD ["./tmate-slave", "-p 2222"]
#RUN mkdir /etc/service/tmate-slave
#ADD tmate-slave.sh /etc/service/tmate-slave/run

#RUN mkdir -p /etc/my_init.d
#ADD message.sh /etc/my_init.d/message.sh
