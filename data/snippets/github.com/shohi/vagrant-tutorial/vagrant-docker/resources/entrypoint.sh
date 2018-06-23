#!/bin/sh

# ref https://stackoverflow.com/questions/43381244/entrypoint-with-environment-variables-is-not-acepting-new-params

java -Dcom.sun.management.jmxremote=true \
     -Dcom.sun.management.jmxremote.port=9010 \
     -Dcom.sun.management.jmxremote.rmi.port=9010 \
     -Dcom.sun.management.jmxremote.local.only=false \
     -Dcom.sun.management.jmxremote.authenticate=false \
     -Dcom.sun.management.jmxremote.ssl=false \
     -Djava.rmi.server.hostname=$SERVER_IP \
     -jar app.war