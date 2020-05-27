.DEFAULT_GOAL := usage

.PHONY: usage
usage:
	@echo "Usage:"
	@echo "=========="
	@echo "make usage - display Makefile target info"
	@echo "make buildlocal - builds the binary locally"
	@echo "make runlocal - runs the binary locally"
	@echo "make builddocker - builds the binary and Docker container"
	@echo "make rundocker - creates and runs a new Docker container"
	@echo "make startdocker - resumes a stopped Docker container"
	@echo "make stopdocker - stops the Docker container"
	@echo "make removedocker - removes the Docker container"
	@echo "make memusage - displays the memory usage of the currently running Docker container"

.PHONY: buildlocal
buildlocal:
	CGO_ENABLED=0 go build -o bin/bot ./...

.PHONY: runlocal
runlocal: buildlocal
	./bin/bot -token=$(ZAHT_BOT_TOKEN)

.PHONY: builddocker
builddocker:
	docker build --tag zaht-bot --file build/Dockerfile .

.PHONY: rundocker
rundocker: builddocker
	docker run \
	--name "zaht_bot" \
	-d --restart unless-stopped \
	-e ZAHT_BOT_TOKEN \
	zaht-bot

.PHONY: startdocker
startdocker:
	docker start zaht_bot

.PHONY: stopdocker
stopdocker:
	docker stop zaht_bot

.PHONY: removedocker
removedocker:
	docker rm zaht_bot

.PHONY: memusage
memusage:
	docker stats zaht_bot --no-stream --format "{{.Container}}: {{.MemUsage}}"
