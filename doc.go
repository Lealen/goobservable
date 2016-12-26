// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package goob (github.com/Lealen/goob) implements data-driven templates for
generating HTML pages right on a client side. It was based on knockout.js,
but i hated it for some reason (it was terrible for using on a large scale
systems), so i created this equivalent in pure Go. It is using a GopherJS
to compile to javascript. Also it is designed to be very fast
and memory-efficient, so you will not waste your time on a project that will
not allow you to later expand your greatness application, also it is very easy
to use, so... go ahead and try it.

ATTENTION!!! This package is under heavy development, many things can break!
Use it at your own risk!

Introduction

This package allows you to use every benefits of html/template package
and also to create very fast, ellegant and super-dynamic web pages,
that really are just like an computer applications.

basic code for index.html looks like:

  <!DOCTYPE html>
  <html>
    <head>
      <meta charset="UTF-8">
      <title>test</title>
    </head>
    <body>
      <div id="app"></div>

      <script src="app.js"></script>
    </body>
  </html>

[[...TODO...]]


attributes that you can use in defined templates:

goob-bind
goob-data
goob-foreach
goob-watch
goob-alias


*/
package goob
