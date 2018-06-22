PARSERS=java go
PARSERS_DIR=parsers

package: clean
	test -d ${PARSERS_DIR} || mkdir ${PARSERS_DIR}
	for p in ${PARSERS}; do make -C $$p package && mv $$p/parser-$$p.zip ${PARSERS_DIR}/; done

clean:
	rm -rf ./${PARSERS_DIR}
