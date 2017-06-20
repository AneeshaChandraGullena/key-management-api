[![Build Status](https://travis.ibm.com/Alchemy-Key-Protect/key-management-api.svg?token=xzhLDR8UgxkQseLgi8S2&branch=develop)](https://travis.ibm.com/Alchemy-Key-Protect/key-management-api)
![Coverage](https://pages.github.ibm.com/Alchemy-Key-Protect/key-management-api/coverage-report/develop/badge.svg)
# key-management-api

Golang based API server

## Dependencies
We use [glide](https://glide.sh) for dependency management, so install it.
1.  `glide install`
2.  go build

## Run tests
1.  See our .travis.yml file

## Usage
1. Create your fork of this repo.
2. Setup and install go 1.6.  Make sure your GOPATH is set correctly: https://golang.org/doc/code.html#GOPATH
3. `go build`
4. `./key-management-api 2> >(tee kp-api.log >&2)`  # This logs to file and console
