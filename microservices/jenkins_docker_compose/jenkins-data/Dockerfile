# Needs to match the linux version used in jenkins-master
FROM debian:jessie 
MAINTAINER Robert Watkins

# user 1000 must match the user id for the jenkins user in jenkins-master
RUN useradd -d "/var/jenkins_home" -u 1000 -m -s /bin/bash jenkins

RUN mkdir -p /var/log/jenkins
RUN chown -R jenkins:jenkins /var/log/jenkins

VOLUME ["/var/log/jenkins", "/var/jenkins_home"]

USER jenkins

CMD ["echo", "Data container for Jenkins"]
