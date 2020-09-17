<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [dockfmt](#dockfmt)
  - [Installation](#installation)
      - [Binaries](#binaries)
      - [Via Go](#via-go)
  - [Usage](#usage)
    - [Format](#format)
      - [Get help](#get-help)
      - [Get a diff](#get-a-diff)
      - [List multiple files with different output](#list-multiple-files-with-different-output)
    - [Base image inspection](#base-image-inspection)
    - [Maintainer inspection](#maintainer-inspection)
    - [Stage inspection](#stage-inspection)
    - [Workdir inspection](#workdir-inspection)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# dockfmt

[![make-all](https://github.com/jessfraz/dockfmt/workflows/make%20all/badge.svg)](https://github.com/jessfraz/dockfmt/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/jessfraz/dockfmt/workflows/make%20image/badge.svg)](https://github.com/jessfraz/dockfmt/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/dockfmt)

Dockerfile format.

**NOTE:** This is a work in progress so calm yourself if you want to file 80 bajillion
issues.

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
- [Usage](#usage)
  * [Format](#format)
    + [Get help](#get-help)
    + [Get a diff](#get-a-diff)
    + [List multiple files with different output](#list-multiple-files-with-different-output)
  * [Base image inspection](#base-image-inspection)
  * [Maintainer inspection](#maintainer-inspection)
  * [Stage inspection](#stage-inspection)
  * [Workdir inspection](#workdir-inspection)

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/dockfmt/releases).

#### Via Go

```console
$ go get github.com/jessfraz/dockfmt
```

## Usage

```console
$ dockfmt -h
dockfmt -  Dockerfile format.

Usage: dockfmt <command>

Flags:

  -d, --debug  enable debug logging (default: false)

Commands:

  base        List the base image used in the Dockerfile(s).
  dump        Dump parsed Dockerfile(s).
  fmt         Format the Dockerfile(s).
  maintainer  List the maintainer for the Dockerfile(s).
  stages      List the stages in the Dockerfile.
  workdir     List the workdirs for the Dockerfile(s).
  version     Show the version information.
```

### Format

#### Get help

```console
$ dockfmt fmt -h
Usage: dockfmt fmt [OPTIONS] [DOCKERFILE...]

Format the Dockerfile(s).

Flags:

  -D, --diff   display diffs instead of rewriting files (default: false)
  -d, --debug  enable debug logging (default: false)
  -l, --list   list files whose formatting differs from dockfmt's (default: false)
  -w, --write  write result to (source) file instead of stdout (default: false)
```

#### Get a diff

```console
$ dockfmt fmt -d htop/Dockerfile
diff htop/Dockerfile dockfmt/htop/Dockerfile
--- /tmp/dockfmt143910590	2016-09-19 15:59:22.612250710 -0700
+++ /tmp/dockfmt412224773	2016-09-19 15:59:22.612250710 -0700
@@ -4,10 +4,11 @@
 # 	--pid host \
 # 	jess/htop
 #
-FROM alpine:latest
-MAINTAINER Jessie Frazelle <jess@linux.com>
+
+FROM	alpine:latest
+MAINTAINER	Jessie Frazelle <jess@linux.com>

-RUN apk --no-cache add \
+RUN	apk add --no-cache \
 	htop

-CMD [ "htop" ]
+CMD	["htop"]
```

#### List multiple files with different output

```console
$ dockfmt fmt -l */Dockerfile */*/Dockerfile
ab/Dockerfile
afterthedeadline/Dockerfile
android-tools/Dockerfile
ansible/Dockerfile
apt-file/Dockerfile
atom/Dockerfile
audacity/Dockerfile
awscli/Dockerfile
beeswithmachineguns/Dockerfile
buttslock/Dockerfile
camlistore/Dockerfile
cathode/Dockerfile
...
```

### Base image inspection

```console
$ dockfmt base */Dockerfile */*/Dockerfile
BASE                          COUNT
debian:stretch                50
alpine:latest                 30
debian:sid                    28
ubuntu:16.04                  12
alpine:edge                   7
python:2-alpine               3
ruby:alpine                   2
java:7-alpine                 2
r.j3ss.co/wine                1
kalilinux/kali-linux-docker   1
haskell                       1
mhart/alpine-node:5           1
r.j3ss.co/cpuminer            1
opensuse                      1
java:8-alpine                 1
golang:latest                 1
```

### Maintainer inspection

```console
$ dockfmt maintainer */Dockerfile */*/Dockerfile
MAINTAINER                                      COUNT
Jessie Frazelle <jess@jskdj.com>                113
Christian Koep <christian.koep@ksldkfj.de>      11
Justin Garrison <justinleegarrison@hskdl.com>   2
Daniel Romero <infoslack@jjskl.com>             1
Cris G c@skdlemfhtj.com                         1
```

### Stage inspection

```console
$ dockfmt stages Dockerfile
STAGE               INTERPOLATED
health-check        false
python-deps         false
stage-2             true
```

### Workdir inspection

```console
$ dockfmt workdir */Dockerfile */*/Dockerfile
WORKDIR   COUNT
/srv      3
/app      1
```
