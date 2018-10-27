[![Build Status](https://travis-ci.org/0xC0D3D00D/goresp.svg?branch=master)](https://travis-ci.org/0xC0D3D00D/goresp)
[![Coverage Status](https://coveralls.io/repos/github/0xC0D3D00D/goresp/badge.svg?branch=master)](https://coveralls.io/github/0xC0D3D00D/goresp?branch=master)
[![codebeat badge](https://codebeat.co/badges/5bdcc4c1-864f-40e3-b6ca-d36b2e59c851)](https://codebeat.co/projects/github-com-0xc0d3d00d-goresp-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xc0d3d00d/goresp)](https://goreportcard.com/report/github.com/0xc0d3d00d/goresp)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2310/badge)](https://bestpractices.coreinfrastructure.org/projects/2310)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2F0xC0D3D00D%2Fgoresp.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2F0xC0D3D00D%2Fgoresp?ref=badge_shield)

# goresp

Go implementation for RESP (REdis Serialization Protocol)

# Documentation
For documentation see [Godoc](https://godoc.org/github.com/0xc0d3d00d/goresp).

# Data types
This table has examples for how data will be encoded and decoded.

RESP Representation                   | Human-readable | RESP Type      | Go Representation
--------------------------------------|----------------|----------------|-----------------------------------------------------------
"\r\n"                                | (empty)        | undefined      | nil
":1\r\n"                              | 1              | Integer        | int64(1)
"+abc\r\n"                            | "abc"          | Simple String  | []byte{'a','b','c'}
"+\r\n"                               | ""             | Simple String  | []byte{}
"-abc\r\n"                            | "abc"          | Error          | error (msg=abc)
"-\r\n"                               | ""             | Error          | error
"$5\r\nabc\r\n\r\n"                   | "abc\r\n"      | Bulk String    | []byte{'a','b','c'}
"$0\r\n\r\n"                          | ""             | Bulk String    | []byte{}
"*0\r\n"                              | []             | Array          | []interface{}{}
"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"    | ["foo", "bar"] | Array          | []interface{}{[]byte{'f','o','o'}, []byte{'b','a','r'}}
"*2\r\n:1\r\n$3\r\nfoo\r\n"           | [6, "foobar"]  | Array (mixed)  | []interface{}{int64(1), []byte{'f', 'o', 'o'}}

# Benchmarks
Benchmark name                              | (1)        | (2)         | (3) 		    | (4)
--------------------------------------------|-----------:|------------:|-----------:|---------:
BenchmarkReadSmallString                    | 20000000   |       88.7  |      16    |    2
BenchmarkReadInteger                        | 10000000   |        126  |      16    |    3
BenchmarkReadBulkString                     | 20000000   |        107  |      16    |    3
BenchmarkReadArray                          |  3000000   |        580  |     176    |   14

- (1): Total Repetitions achieved in constant time, higher means more confident result
- (2): Single Repetition Duration (ns/op), lower is better
- (3): Heap Memory (B/op), lower is better
- (4): Average Allocations per Repetition (allocs/op), lower is better

# TODO List
[ ] Change Simple string decoded type to string

# License
Unless otherwise noted, All source files in this library are distributed under the Apache License 2.0 found in the LICENSE file.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2F0xC0D3D00D%2Fgoresp.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2F0xC0D3D00D%2Fgoresp?ref=badge_large)
