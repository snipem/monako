# monako

![Run Monako](https://github.com/snipem/monako/workflows/Run%20Monako/badge.svg?branch=develop)

A less opinionated document aggregator and publisher. Easier to use and to adapt than Antora, Hugo, JBake and the likes.

## Usage

```
$ monako -h
Usage of ./monako:
  -config string
        Configuration file, default: config.yaml (default "config.yaml")
  -hugo-config string
        Configuration file for hugo, default: config.toml (default "config.toml")
  -menu-config string
        Menu file for hugo-book theme, default: index.md (default "index.md")
  -trace
        Enable trace logging
```

## Design Goals

* Make it simple and stupid
* Git clone existing repositories
  * Don't attempt to change their structure
  * Transform all documents below a specific structure
* Use Antora like configuration file for main structure
* Support Mermaid

## Development

Init with `make init`

## TODOs

* Fail on wrong `docdir` ("TODO")
* Make white list configurable