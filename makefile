PROG_NAME := concurrency_patterns

build:
	go build -o $(PROG_NAME) *.go

ARG?=
run: build
	./$(PROG_NAME) $(ARG)

clear:
	rm ./$(PROG_NAME)
