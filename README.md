# monako

![Build Monako](https://github.com/snipem/monako/workflows/Build%20Monako/badge.svg?branch=develop)
[![codecov](https://codecov.io/gh/snipem/monako/branch/master/graph/badge.svg)](https://codecov.io/gh/snipem/monako)
[![Go Report Card](https://goreportcard.com/badge/github.com/snipem/monako)](https://goreportcard.com/report/github.com/snipem/monako)
[![GoDoc](https://godoc.org/github.com/snipem/monako?status.svg)](https://godoc.org/github.com/snipem/monako)
[![Maintainability](https://api.codeclimate.com/v1/badges/1ff16e0c4f8a871bfac3/maintainability)](https://codeclimate.com/github/snipem/monako/maintainability)

![monako logo](https://github.com/snipem/monako/raw/master/assets/logo/cover.png)

----

A less opinionated document aggregator and publisher. Easier to use and to adapt than Antora, Hugo, JBake and the likes.

----

## Purpose

Monako abstracts the complexity of collecting documentation from a configurable amount of Git repositories (origins).

Monako uses [Hugo](https://gohugo.io) as a static site generator but hides it to the user. It also uses the great [Hugo Book Theme](https://github.com/alex-shpak/hugo-book) which I forked at [Monako Book](https://github.com/snipem/monako-book).

## How Monako works

![How Monako works](https://github.com/snipem/monako/raw/master/assets/monako.png)

## Usage

```help
$ monako -h
Usage of monako:
  -base-url string
        Custom base URL
  -config string
        Configuration file (default "config.monako.yaml")
  -fail-on-error
        Fail on document conversion errors
  -menu-config string
        Menu file for monako-book theme (default "config.menu.md")
  -compose
        Only compose the Monako structure
  -render
        Only render HTML files from an existing Monako structure
  -trace
        Enable trace logging
```

A Docker image is available from [Dockerhub](https://hub.docker.com/repository/docker/snipem/monako).

## Configuration

### Configuration of Origins

```yaml
---
  baseURL : "https://example.com/"
  title : "My Projects"

  whitelist:
    - ".md"
    - ".adoc"
    - ".jpg"
    - ".svg"
    - ".png"

  origins:
  - src: https://github.com/snipem/commute-tube
    branch: master
    docdir: .
    targetdir: docs/commute

  - src: https://github.com/snipem/monako
    branch: develop
    docdir: doc
    targetdir: docs/monako
```

### Configuration of Menus

```markdown
- **My Projects**
  - [Commute Tube]({{< relref "/docs/commute/readme" >}})
  - [Monako]({{< relref "/docs/monako/readme" >}})
  ...
```

### Configuration of Documents

Monako supports all [Hugo Frontmatter](https://gohugo.io/content-management/front-matter/) types (YAML, TOML and JSON).
Monako converts them to YAML.

Add frontmatter as you wish at long as it's supported by Hugo and the Theme.

#### Monako specific options

Hide Git links like "edit this page" and "last edit by". Add this line to the frontmatter of the document:

```yaml
MonakoGitLinks = false
```

### Screenshot

![Screenshot of a documentation site built with Monako](https://github.com/snipem/monako/raw/master/assets/screenshot.png)

## Design Goals

* Make it simple and stupid, single binary with no dependencies ✓
* Git clone existing repositories ✓
  * Don't attempt to change their structure ✓
  * Transform all documents below a specific structure ✓
* Use Antora like configuration file for main structure ✓
* Support Mermaid ✓

## Features

* Multi repository document fetching
* Bundled dependencies (minus Asciidoctor)

## Development

Init with `make init`

## Improvements

* Support edit this page links with links to origin repository. Book theme provides a [similar feature](https://github.com/alex-shpak/hugo-book/search?q=BookRepo&unscoped_q=BookRepo).
