NAME = srccat
EXEC = srccat.jar
VERSION = 1.0.0
DIR = ${NAME}-${VERSION}

build:
	sbt assembly
	mv target/scala-2.11/srccat-assembly-1.0.0.jar ./${EXEC}
	
package: clean build
	test -d ${DIR} || mkdir ${DIR}
	cp ${EXEC} ${DIR}/
	cp README.md ${DIR}/
	tar czvf ${DIR}.tar.gz ${DIR}
	rm -rf ${DIR}
	
clean:
	sbt clean