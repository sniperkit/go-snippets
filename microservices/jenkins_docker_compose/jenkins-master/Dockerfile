FROM jenkins:2.7.1
MAINTAINER Robert Watkins

USER root
RUN mkdir /var/log/jenkins
RUN mkdir /var/cache/jenkins
RUN chown -R  jenkins:jenkins /var/log/jenkins
RUN chown -R  jenkins:jenkins /var/cache/jenkins

USER jenkins
ENV JAVA_OPTS="-Xmx4096m"

# install plugins; the plugins.txt file can be exported from Jenkins like this:
# JENKINS_HOST=username:password@myhost.com:port
# curl -sSL "http://$JENKINS_HOST/pluginManager/api/xml?depth=1&xpath=/*/*/shortName|/*/*/version&wrapper=plugins" | perl -pe 's/.*?<shortName>([\w-]+).*?<version>([^<]+)()(<\/\w+>)+/\1 \2\n/g'|sed 's/ /:/' > jenkins-master/conf/plugins.txt

#COPY conf/plugins.txt /var/jenkins_home/plugins.txt
#RUN /usr/local/bin/plugins.sh /var/jenkins_home/plugins.txt

# "For 2.x-derived images, you may also want to" - we're 2.0 dervied, so we want this
RUN echo 2.0 > /usr/share/jenkins/ref/jenkins.install.UpgradeWizard.state

# Put the log file into the log directory, which will be in the data volume
# Move the WAR out of the persisted jenkins data dir
ENV JENKINS_OPTS="--logfile=/var/log/jenkins/jenkins.log --webroot=/var/cache/jenkins/war"
