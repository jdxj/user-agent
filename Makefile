.PHONY: clean

fileName="user-agent.out.bak"

upload: build
	scp ./${fileName} root@hd1h.ssh.aaronkir.xyz:user-agent
build: main.go
	go build -o ${fileName} main.go
clean:
	rm -vf *.out *.log *.bak
