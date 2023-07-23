
# TCP pooling

[![tag](https://img.shields.io/github/tag/samber/go-tcp-pool.svg)](https://github.com/samber/go-tcp-pool/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18.0-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/go-tcp-pool?status.svg)](https://pkg.go.dev/github.com/samber/go-tcp-pool)
![Build Status](https://github.com/samber/go-tcp-pool/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/go-tcp-pool)](https://goreportcard.com/report/github.com/samber/go-tcp-pool)
[![Coverage](https://img.shields.io/codecov/c/github/samber/go-tcp-pool)](https://codecov.io/gh/samber/go-tcp-pool)
[![Contributors](https://img.shields.io/github/contributors/samber/go-tcp-pool)](https://github.com/samber/go-tcp-pool/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/go-tcp-pool)](./LICENSE)

✨ Drop-in replacement to `net.Conn` with pooling and auto-reconnect.

## 🚀 Install

```sh
go get github.com/samber/go-tcp-pool
```

This library is v1 and follows SemVer strictly. No breaking changes will be made to exported APIs before v2.0.0.

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/samber/go-tcp-pool](https://pkg.go.dev/github.com/samber/go-tcp-pool)

### Create a TCP connection pool

```sh
# Start a tcp server
ncat -l 9999 -k
```

```go
import pool "github.com/samber/go-tcp-pool"

conn, err := pool.Dial("tcp", "localhost:9999")
if err != nil {
    log.Fatal(err)
}

conn.SetPoolSize(10)
conn.SetMaxRetries(10)
conn.SetRetryInterval(10 * time.Millisecond)

// a tcp connection will be used in a round-robin manner
n, err := conn.Write([]byte("Hello, world!\n"))
if err != nil {
    log.Fatal(err)
}

// will always return an error
conn.Read(...)
```

## 🚀 @TODO

- [x] Implement round-robin connection pool
- [x] Implement auto-reconnect
- [ ] Implement Read()
- [ ] Implement other load-balancing strategies
  - Max idle time
  - MinConn + MaxConn

## 🤝 Contributing

- Ping me on Twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/go-tcp-pool)
- Fix [open issues](https://github.com/samber/go-tcp-pool/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/go-tcp-pool)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
