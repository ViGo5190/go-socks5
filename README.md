# go-socks5

[![Build Status](https://travis-ci.org/vigo5190/go-socks5.svg?branch=master)](https://travis-ci.org/vigo5190/go-socks5)

Simple golang socks5 realization.


How to use
----------

```bash
    docker pull vigo5190/gosocks5
    docker run -d -p 5190:8008 vigo5190/gosocks5
```

How to use (not docker)
------------------------

Build:

```bash
    make
```

Run:

```bash
    ./go-socks5 -port=8009 -addr=0.0.0.0
```