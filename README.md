# Vamos

__Author__: Kevin P. Albrecht - <http://www.kevinalbrecht.com>
__Web Site__: <http://www.github.com/onlyafly/vamos>

## What is Vamos?

Vamos is a Lisp written in Go.

## Running Test Suite

    $ make test

## Building From Source

    $ make

## Start REPL

    $ vamos

## Load forms from file then start REPL

    $ vamos -l foo.v

## Development

### Add a new dependency

    $ go get -u foo/bar
    $ godep save ./...

### Updating an existing dependency

    $ go get -u foo/bar
    $ godep update foo/bar
