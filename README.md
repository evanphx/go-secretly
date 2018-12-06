# Go Secertly Young Person!

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/evanphx/go-secretly)

This package provides a simple abstraction over different ways of storing a secret value by name.

It supports the ability to store these values in:
* `file`: encrypted with chacha20poly1305 using a trivial key or user supplied
* `awsps`: stores and retrieves secrets from AWS Parameter Storage as SecureStrings
* `vault`: stores and retrieves secrets from HashiCorp vault. Supports generic, kv 1, and kv 2


Should it support another storage mechanism? Open an issue! Or better yet, send a PR!
