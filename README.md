# dockfmt

[![Travis
CI](https://travis-ci.org/jessfraz/dockfmt.svg?branch=master)](https://travis-ci.org/jessfraz/dockfmt)

Dockerfile format.

This is a work in progress so calm yourself if you want to file 80 bajillion
issues.

## Installation

#### Binaries

- **darwin** [386](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-darwin-386) / [amd64](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-darwin-amd64)
- **freebsd** [386](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-freebsd-386) / [amd64](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-freebsd-amd64)
- **linux** [386](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-linux-386) / [amd64](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-linux-amd64) / [arm](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-linux-arm) / [arm64](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-linux-arm64)
- **solaris** [amd64](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-solaris-amd64)
- **windows** [386](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-windows-386) / [amd64](https://github.com/jessfraz/dockfmt/releases/download/v0.3.0/dockfmt-windows-amd64)

#### Via Go

```bash
$ go get github.com/jessfraz/dockfmt
```

## Usage

**Help output**

```console
$ dockfmt -h
NAME:
   dockfmt - Dockerfile format.

USAGE:
   dockfmt [global options] command [command options] [arguments...]

VERSION:
   v0.3.0

AUTHOR(S):
   @jessfraz <no-reply@butts.com>

COMMANDS:
     base         list the base image used in Dockerfile(s)
     dump         dump parsed Dockerfile(s)
     format, fmt  format the Dockerfile(s)
     maintainer   list the maintainer for Dockerfile(s)
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -D    run in debug mode
   --help, -h     show help
   --version, -v  print the version

```

### Format

**help output**

```console
$ dockfmt fmt -h
NAME:
   dockfmt format - format the Dockerfile(s)

USAGE:
   dockfmt format [command options] [arguments...]

OPTIONS:
   --diff, -d   display diffs instead of rewriting files
   --list, -l   list files whose formatting differs from dockfmt's
   --write, -w  write result to (source) file instead of stdout
```

**get a diff**

```console
$ dockfmt format -d htop/Dockerfile
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

**list multiple files with different output**

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

