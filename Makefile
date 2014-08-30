all: hello node_modules bower_components public/hello.js public/hello.html public/vendor.js public/hello.css

hello: hello.go

%: %.go
	go build $*.go

node_modules:
	npm -q update

bower_components:
	bower -q update

public/hello.js: app/hello.ls
	gulp

public/hello.html: app/hello.jade
	gulp

public/vendor.js: bower.json
	gulp

public/hello.css: app/hello.styl
	gulp
