# qiloop

[![Build Status](https://travis-ci.org/lugu/qiloop.svg?branch=master)](https://travis-ci.org/lugu/qiloop)
[![Documentation](https://godoc.org/github.com/lugu/qiloop?status.svg)](http://godoc.org/github.com/lugu/qiloop)
[![license](https://img.shields.io/github/license/lugu/qiloop.svg?maxAge=2592000)](https://github.com/lugu/qiloop/blob/master/LICENSE)
[![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)
[![Release](https://img.shields.io/github/tag/lugu/qiloop.svg)](https://github.com/lugu/qiloop/releases)

[![CircleCI](https://circleci.com/gh/lugu/qiloop/tree/master.svg?style=shield)](https://circleci.com/gh/lugu/qiloop/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugu/qiloop)](https://goreportcard.com/report/github.com/lugu/qiloop)
[![codecov](https://codecov.io/gh/lugu/qiloop/branch/master/graph/badge.svg)](https://codecov.io/gh/lugu/qiloop)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b192466a26dbced44274/test_coverage)](https://codeclimate.com/github/lugu/qiloop/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/b192466a26dbced44274/maintainability)](https://codeclimate.com/github/lugu/qiloop/maintainability)

**`qiloop`** is an implementation of QiMessaging written in [Go](https://golang.org).

QiMessaging is a network protocol used to build rich distributed
applications. It was created by Aldebaran Robotics (currently
[SoftBank Robotics](https://www.softbankrobotics.com/emea/en/index))
and is the foundation of the NAOqi SDK. For more details about
QiMessaging, visit this [analysis of the
protocol](https://github.com/lugu/qiloop/blob/master/doc/NOTES.md).

## Installation

    go get github.com/lugu/qiloop/...

## Tutorials

By default, `qiloop` comes with two proxies: ServiceDirectory and
LogManager.

Follow the [ALVideoDevice tutorial](https://github.com/lugu/qiloop/blob/master/doc/TUTORIAL.md)
to learn how to create a proxy to an existing service.

Follow the [clock tutorial](https://github.com/lugu/qiloop/blob/master/doc/SERVICE_TUTORIAL.md)
to create your own service.

## Examples

The [examples directory](https://github.com/lugu/qiloop/blob/master/examples/)
illustrates some basic usages of qilooop:

-   [method call](https://github.com/lugu/qiloop/blob/master/examples/method)
    illustrates how to call a method of a service: this example lists
    the services registered to the service directory.

-   [signal registration](https://github.com/lugu/qiloop/blob/master/examples/signal)
    illustrates how to subscribe to a signal: this example prints a
    log each time a service is added to the service directory.

-   [ping pong service](https://github.com/lugu/qiloop/blob/master/examples/pong)
    illustrates how to implement a service.

-   [space service](https://github.com/lugu/qiloop/blob/master/examples/space)
    illustrates the client side objects creation.

-   [clock service](https://github.com/lugu/qiloop/blob/master/examples/space)
    completed version of the clock tutorial.

## Authentication

If you need to provide a login and a password to authenticate yourself
to a server, create a file `$HOME/.qi-auth.conf` with you login on the
first line and your password on the second.

## Status

This is work in progress, you have been warned.

The client and the server side is working: one can implement a service
from an IDL and generate a specialized proxy for this service.
A service directory is implemented as part of the standalone server.

What is working:

-   TCP and TLS connections
-   client proxy generation
-   server stub generation
-   method, signals and properties
-   Authentication
-   IDL parsing and generation
