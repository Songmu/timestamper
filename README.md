timestamper
=======

[![Build Status](https://travis-ci.org/Songmu/timestamper.png?branch=main)][travis]
[![Coverage Status](https://coveralls.io/repos/Songmu/timestamper/badge.png?branch=main)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/timestamper?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/timestamper
[coveralls]: https://coveralls.io/r/Songmu/timestamper?branch=main
[license]: https://github.com/Songmu/timestamper/blob/main/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/timestamper

text transformer to put timestamps. It is very useful for logging.

## Description

The text transformer to put timestamps. The timestamper implements
`golang.org/x/text/transform.Transform` interface.

## Synopsis

### Easy Usage

```Go
var s transform.Transformer = timestamper.New()
var w io.WriteCloser = transform.NewWriter(os.Stdout, s)
fmt.Fprint(w, "Hello\nWorld!")
// Output:
// 2019-02-11T01:14:54.093021+09:00 Hello
// 2019-02-11T01:14:54.093151+09:00 World!
```

### Functional Option

```Go
s1 := timestamper.New(timestamper.UTC()) // use UTC timestamp
s2 := timestamper.New(timestamper.Layout("06-01-02 15:04:05 ")) // specify custom layout
```

## Installation

    % go get github.com/Songmu/timestamper

## Author

[Songmu](https://github.com/Songmu)
