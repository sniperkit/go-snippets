default: start
project:=ms-workspace-demo
service:=ms-demo-golang

.PHONY: start
start: 
	docker-compose -p ${project} up -d

.PHONY: stop
stop: 
	docker-compose -p ${project} down

.PHONY: restart
restart: stop start

.PHONY: logs
logs: 
	docker-compose -p ${project} logs -f ${service}

.PHONY: ps
ps: 
	docker-compose -p ${project} ps

.PHONY: shell
shell: 
	docker-compose -p ${project} exec ${service} sh

.PHONY: build
build:
	docker-compose -p ${project} build --no-cache

.PHONY: dep-add
dep-add:
	docker-compose -p ${project} exec ${service} dep ensure -add ${package}

.PHONY: dep-update
dep-update:
	docker-compose -p ${project} exec ${service} dep ensure -update ${package}

.PHONY: dep-update-all
dep-update-all:
	docker-compose -p ${project} exec ${service} dep ensure -update
