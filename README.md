go-weibo
========

[![Build Status](https://img.shields.io/travis/larrylv/go-weibo.svg?branch=master)][travis]
[![Coverage Status](http://img.shields.io/coveralls/larrylv/go-weibo.svg)][coveralls]

[travis]: https://travis-ci.org/larrylv/go-weibo
[coveralls]: https://coveralls.io/r/larrylv/go-weibo

go-weibo is a Go client library for access the [Weibo API][], inspired by [go-github][].

**Documentation:** <http://godoc.org/github.com/larrylv/go-weibo/weibo>

This library is __currently under development__, please do not use it in production environment.

## Usage ##

``` go
import "github.com/larrylv/go-weibo/weibo"
```

The go-weibo library does not directly handle authentication.  Instead, when
creating a new client, pass an `AccessToken` that will be added to every request's
header, and handle authentication for you.  The easiest and recommended way to
get an `AccessToken` is using the [goauth2][] library.

For example, to update a weibo:

```go
accessToken = "access_token"
client := weibo.NewClient(accessToken)

// Update a weibo
opts = &weibo.StatusRequest{Status: weibo.String("Hello, Weibo!")}
status, _, err := client.Statuses.Create(opts)
```

See the [goauth2 docs][] for complete instructions on using that library.

For complete usage of go-weibo, see the full [package docs][].

[Weibo API]: http://open.weibo.com/wiki/%E5%BE%AE%E5%8D%9AAPI
[go-github]: https://github.com/google/go-github/
[goauth2]: https://code.google.com/p/goauth2/
[goauth2 docs]: http://godoc.org/code.google.com/p/goauth2/oauth
[package docs]: http://godoc.org/github.com/larrylv/go-weibo/weibo

## License

This library is released under the MIT License.
