HiYoga CLI
==========

[![Build Status](https://travis-ci.org/tazjin/hiyoga.svg?branch=master)](https://travis-ci.org/tazjin/hiyoga)

This is a simple CLI tool for accessing the [HiYoga][] API. It currently supports listing classes for a specified number
of days in advance, as well as describing the different types of classes available.

![Example screenshot](http://i.imgur.com/0NOk1SK.png)

Just `go get github.com/tazjin/hiyoga`.

## Usage

```
USAGE:
   hiyoga [global options] command [command options] [arguments...]

COMMANDS:
     list-classes, lc       list upcoming yoga classes
     list-class-types, lct  list available yoga class types
     help, h                Shows a list of commands or help for one command
```

## Authentication support

The HiYoga CLI supports authentication for doing things like listing bookings and booking new classes.

To make use of this, place a configuration file like this in `$HOME/.hiyoga`:

```json
{
  "username": "some.email@example.com",
  "password": "TremendousPassword2017"
}

```

## Disclaimer

This project is not in any way affiliated with HiYoga. I don't work for the company and if this breaks in any way it's
entirely my fault as this API is not designed for public consumption.

[HiYoga]: https://www.hiyoga.no/
