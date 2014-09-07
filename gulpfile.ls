require! gulp
require! "main-bower-files"
$ = do (require "gulp-load-plugins")

getJs = ->
  gulp.src "app/hello.ls"
    .pipe $.livescript!
    .pipe $.size showFiles: true
    .pipe $.ngAnnotate!

gulp.task "default", ->
  getJs!
    .pipe gulp.dest "public"

gulp.task "closure-compiler", ->
  getJs!
    .pipe $.closureCompiler(
      compilerPath: "bower_components/closure-compiler/compiler.jar"
      fileName: "hello.js"
      warning_level: "VERBOSE"
      externs:
        "support/externs/angular-1.2.js"
        "support/externs/externs.js"
    )
    .pipe $.size showFiles: true
    .pipe gulp.dest "tmp"

  gulp.src "app/index.jade"
    .pipe $.jade!
    .pipe gulp.dest "public"

  gulp.src "app/hello.styl"
    .pipe $.stylus!
    .pipe gulp.dest "public"

  gulp.src mainBowerFiles!
    .pipe $.size showFiles: true
    .pipe $.concat "vendor.js"
    .pipe gulp.dest "public"

