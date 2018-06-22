NAME   = parser-java
SCRIPT = parser
JAR_PATH = ./target/javaparser-1.0-SNAPSHOT-jar-with-dependencies.jar

all: clean build

build:
	mvn -DskipTests -Prelease-profile install

package: clean build
	mkdir ${NAME}
	cp ${JAR_PATH} ${NAME}/javaparser.jar
	echo "#!/bin/sh" > ${NAME}/${SCRIPT}
	echo 'java -jar $$(dirname $$0)/javaparser.jar $$@' >> ${NAME}/${SCRIPT}
	chmod +x ${NAME}/${SCRIPT}
	zip -r ${NAME}.zip ${NAME}
	rm -rf ./${NAME}

clean:
	mvn clean && rm -f ${NAME}.zip

