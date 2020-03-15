# monako

![Run Monako](https://github.com/snipem/monako/workflows/Run%20Monako/badge.svg?branch=develop)

A less opininated document aggregator and publisher. Easier to use and to adapt than Antora, Hugo, JBake and the likes.

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

* Really big repos seem to produce `runtime error: invalid memory address or nil pointer dereference` error.
* Fail on wrong docdir ("TODO")
* Make white list configurable