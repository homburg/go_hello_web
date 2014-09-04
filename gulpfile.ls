require! gulp
require! "main-bower-files"
$ = do (require "gulp-load-plugins")

gulp.task "default", ->
	gulp.src "app/hello.ls"
		.pipe $.livescript!
		.pipe $.size showFiles: true
		.pipe $.ngAnnotate!
		.pipe $.closureCompiler(
			compilerPath: "bower_components/closure-compiler/compiler.jar"
			fileName: "hello.js"
		)
		.pipe $.size showFiles: true
		.pipe gulp.dest "public"

	gulp.src "app/hello.jade"
		.pipe $.jade!
		.pipe gulp.dest "public"

	gulp.src "app/hello.styl"
		.pipe $.stylus!
		.pipe gulp.dest "public"

	gulp.src mainBowerFiles!
		.pipe $.size showFiles: true
		.pipe $.concat "vendor.js"
		.pipe gulp.dest "public"

