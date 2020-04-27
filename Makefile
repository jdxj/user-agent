#.PHONY: clean

fileName="user-agent.out"

build: main.go
	go build -o ${fileName} main.go
clean:
	rm *.out *.log
