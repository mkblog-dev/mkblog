build-cli:
	go build -o ./bin/mkblog ./cli

playground-clean:
	rm -rf dist

playground-build: playground-clean
	./bin/mkblog -d playground -o dist

playground-serve:
	http-server dist
