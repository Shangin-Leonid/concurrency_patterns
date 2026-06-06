PROG_NAME := concurrency_patterns.out

build:
	go build -o $(PROG_NAME) *.go

ARG?=
run: build
	./$(PROG_NAME) $(ARG)

clean:
	rm ./$(PROG_NAME)
