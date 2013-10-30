.PHONY: all login realm install clean

all: login realm

login:
	go build -a -o login.a github.com/Blackrush/gofus/login

realm:
	go build -a -o realm.a github.com/Blackrush/gofus/realm

install:
	go install -a github.com/Blackrush/gofus/login
	go install -a github.com/Blackrush/gofus/realm

clean:
	rm realm.a login.a