# Revs

A command line app for accessing and responding to review requests on GitHub.

![Revs sample image](img/sample.gif)

## Features

* Filter review requests
* Open pull request
* Mark as read
* Unsubscribe

## Install

```sh
go install github.com/Shopify/revs@main
```

or

```sh
go install github.com/Shopify/revs@<tag>
```

This will install to `$GOBIN`, `$GOPATH/bin` or `~/go/bin`. Make sure whichever is in your path.

## Development

Revs development is setup to be done in Spin.

```sh
spin up revs
```

Optionally, you can save a GitHub token as a Spin secret that will be  auto-mounted in the instance.

```sh
spin secrets create -u revs-token
```

