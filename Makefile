.PHONY: types

all: hello node_modules bower_components public/hello.js public/index.html public/vendor.js public/hello.css

run: all
	./hello

hello: hello.go

%: %.go
	go build $*.go

node_modules:
	npm -q update

bower_components:
	bower -q update

public/hello.js: app/hello.ls
	gulp

public/index.html: app/index.jade
	gulp

public/vendor.js: bower.json
	gulp

public/hello.css: app/hello.styl
	gulp

types:
	java -jar bower_components/closure-compiler/compiler.jar \
		--warning_level=VERBOSE \
		--externs support/externs/angular-1.2.js \
		--externs support/externs/externs.js \
		--angular_pass \
		tmp/hello.js
