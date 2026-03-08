all: test

test:
	go test -v -count=1 -covermode=atomic ..

test-docker:
	docker build -t golang-design/x/clipboard ..
	docker run --rm --name cb golang-design/x/clipboard
	docker rmi golang-design/x/clipboard