timestamper
=======

[![Build Status](https://github.com/Songmu/timestamper/workflows/test/badge.svg?branch=main)][actions]
[![Coverage Status](https://coveralls.io/repos/Songmu/timestamper/badge.svg?branch=main)][coveralls]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/timestamper?status.svg)][godoc]

[actions]: https://github.com/Songmu/timestamper/actions?workflow=test
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
var w io.WriteCloser = timestamper.NewWriter(os.Stdout)
defer w.Close()
fmt.Fprint(w, "Hello\nWorld!")
// Output:
// 2019-02-11T01:14:54.093021+09:00 Hello
// 2019-02-11T01:14:54.093151+09:00 World!
```

`NewWriter` serializes concurrent `Write` and `Close` calls, including writes to
the underlying writer. Use one shared writer for a shared output stream.
`New` remains available for `transform.NewReader` and `transform.NewWriter`,
but sharing a `transform.Writer` requires external synchronization.

### Functional Option

```Go
s1 := timestamper.New(timestamper.UTC()) // use UTC timestamp
s2 := timestamper.New(timestamper.Layout("06-01-02 15:04:05 ")) // specify custom layout
```

## Installation

    % go get github.com/Songmu/timestamper
    % go install github.com/Songmu/timestamper/cmd/timestamp

## Usage

How to use the included command line tool called timestamp.

```console
% tail -f path/to/log | timestamp
...log with timestamp here...
```

Use `timestamp -version` to show the version and revision.

## Author

[Songmu](https://github.com/Songmu)
