# TODO
# add maven

# Short term:
# assume etap-stubservice checked out on host and mapped vol
# git clone https://github.com/EurostarDigital/etap_stubservice.git
# get maven
PATH=$PATH:/opt/maven/bin
cd /opt
wget http://www.mirrorservice.org/sites/ftp.apache.org/maven/maven-3/3.3.9/binaries/apache-maven-3.3.9-bin.zip
unzip apache-maven-3.3.9-bin.zip && rm apache-maven-3.3.9-bin.zip
ln -s /opt/apache-maven-3.3.9 /opt/maven

# fix the code:
cd /etap_stubservice/src/main/java/com/eurostar/etapstubs/transformers/
mv SBELoadTravelResponseTransformer.java SbeLoadTravelResponseTransformer.java # javac on linux cares about case!
# remove hard-coded (user-specific) path for generated output
sed -i 's!Users/Kye/Desktop/ETAP_data.txt!var/tmp/ETAP_data.txt!' src/test/java/com/eurostar/etapstubs/DataGeneratorTest.java
cd /etap_stubservice
mvn clean install # if successful, jars will be created under ./target dir

You can run from here (or if you move the jar, amend the -cp option below (for classpath)
so java can find the stubs jar.


java  \
-Dsbe_loadtravel_fixedDelay=2000  \
-Dsbe_pah_fixedDelay=1000  \
-Dsbe_mticket_fixedDelay=1000  \
-Dpassbook_email_fixedDelay=200  \
-Dpassbook_create_fixedDelay=1000  \
-Dnotification_purge_fixedDelay=200  \
-Dnotification_update_fixedDelay=200  \
-cp "target/etap-stubs-1.0.jar:wiremock-1.57-standalone.jar" com.github.tomakehurst.wiremock.standalone.WireMockServerRunner  \
--verbose  \
--port 9001  \
--extensions com.eurostar.etapstubs.transformers.ProviderResponseTransformer

