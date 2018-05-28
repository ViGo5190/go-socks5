# go-socks5

[![Build Status](https://travis-ci.org/vigo5190/go-socks5.svg?branch=master)](https://travis-ci.org/vigo5190/go-socks5)

Simple golang socks5 realization.


How to use
----------

```bash
    docker pull vigo5190/gosocks5
    docker run -d -p 5190:8008 vigo5190/gosocks5
```

How to build (local)
----

```bash
    make
```


Configuration
--------------

By default your app run on `0.0.0.0:8008` without auth

How to get password :
```bash
    htpasswd -nbs  yourUsername yourPassword
```

You will get something like : `yourUsername:{SHA}pNvEyEKjh5XJ9cGgtK0l0WiuwmM=`
You need part after `:{SHA}`. for this example its `pNvEyEKjh5XJ9cGgtK0l0WiuwmM=
`

Example:

```toml
listen ="0.0.0.0:8008"

auth = true
[[users]]
    login = "vigo5190"
    pass = "Ys23Ag/5IOWqZCw9QGaVDdHwH00="
```