# Intentionally vulnerable Golang project

[![Build Status](https://travis-ci.org/sonatype-nexus-community/intentionally-vulnerable-golang-project.svg?branch=master)](https://travis-ci.org/sonatype-nexus-community/intentionally-vulnerable-golang-project) [![CircleCI](https://circleci.com/gh/sonatype-nexus-community/intentionally-vulnerable-golang-project.svg?style=shield)](https://circleci.com/gh/sonatype-nexus-community/intentionally-vulnerable-golang-project)

This is just a minimal repo for testing Sonatype's `nancy` against an intentionally vulnerable list of 
dependencies, and as well showing a small example of how to use it in Travis-CI and CircleCI

Project is currently setup to use both `dep` and `go mod` so you should be able to use either one. 

To see how `nancy` will output when finding vulnerabilities, check out [this build on Travis-CI](https://travis-ci.org/github/sonatype-nexus-community/intentionally-vulnerable-golang-project/builds/671448888) or [this build on CircleCI](https://circleci.com/gh/sonatype-nexus-community/intentionally-vulnerable-golang-project/26)

P2 dual-write test 2026-09-21T19:23:53Z
P2 dual-write test (synchronize) 19:24:29Z

P3 envelope test 2026-09-22T17:55:32Z
P3 envelope test (synchronize) 17:55:49Z

W3 clone test 2026-09-22T18:25:11Z
W3 clone test (synchronize) 18:25:26Z

W4 reader test 2026-09-23T15:57:48Z
W4 reader test (synchronize) 15:58:19Z

P5 dispatch test 2026-09-23T19:55:42Z
P5 dispatch test (synchronize) 19:56:30Z

Fanout regression test 2026-09-25T01:31:51Z
Fanout regression test (synchronize) 01:32:48Z
