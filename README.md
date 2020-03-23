# monako

![Build Monako](https://github.com/snipem/monako/workflows/Build%20Monako/badge.svg?branch=develop)
[![Go Report Card](https://goreportcard.com/badge/github.com/snipem/monako)](https://goreportcard.com/report/github.com/snipem/monako)
[![GoDoc](https://godoc.org/github.com/snipem/monako?status.svg)](https://godoc.org/github.com/snipem/monako)
[![Maintainability](https://api.codeclimate.com/v1/badges/1ff16e0c4f8a871bfac3/maintainability)](https://codeclimate.com/github/snipem/monako/maintainability)

![monako logo](res/logo/cover.png)

----

A less opinionated document aggregator and publisher. Easier to use and to adapt than Antora, Hugo, JBake and the likes.

----

## Usage

```help
$ monako -h
Usage of monako:
  -base-url string
        Custom base URL
  -config string
        Configuration file (default "config.monako.yaml")
  -menu-config string
        Menu file for monako-book theme (default "config.menu.md")
  -trace
        Enable trace logging
```

A Docker image is available from [Dockerhub](https://hub.docker.com/repository/docker/snipem/monako).

## Design Goals

* Make it simple and stupid
* Git clone existing repositories
  * Don't attempt to change their structure
  * Transform all documents below a specific structure
* Use Antora like configuration file for main structure
* Support Mermaid

## Features

* Multi repository document fetching
* Bundled dependencies (minus Asciidoctor)

## Development

Init with `make init`

## TODOs

* Fail on wrong `docdir` ("TODO")
* Fail on build fail to prevent false deployment to production
* Move TOC fix to [monako-book](https://github.com/snipem/monako-book)
* Reduce uncritical warnings by Asciidoctor

## Improvements

* Support edit this page links with links to origin repository. Book theme provides a [similar feature](https://github.com/alex-shpak/hugo-book/search?q=BookRepo&unscoped_q=BookRepo).
