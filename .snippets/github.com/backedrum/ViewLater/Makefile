ASSEMBLY=./viewlater

.DEFAULT_GOAL: $(ASSEMBLY)

$(ASSEMBLY):
	$(GOPATH)/bin/gopherjs build app.go -o ${ASSEMBLY}/js/app.js
	mkdir ${ASSEMBLY}/css
	mkdir ${ASSEMBLY}/fonts
	cp manifest.json ${ASSEMBLY}
	cp index.html ${ASSEMBLY}
	cp css/*.css ${ASSEMBLY}/css
	cp js/*.js ${ASSEMBLY}/js
	cp fonts/*.* ${ASSEMBLY}/fonts
clean:
	if [ -d ${ASSEMBLY} ] ; then rm -rf ${ASSEMBLY} ; fi
