.PHONY: all deps login realm install clean

all: deps install

deps:
	go get -v ./...

install:
	go install -a -v ./...

login:
	go build -a -o login.a github.com/Blackrush/gofus/login

realm:
	go build -a -o realm.a github.com/Blackrush/gofus/realm

clean:
	rm realm.a login.a