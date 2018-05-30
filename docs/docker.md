# Docker

## Overview

The docker image for `scribe` is a whopping `14MB`. The following documentation explains how to use `scribe` as a docker image:

## Building

To build the docker image, simply run the following command from the project directory:

```
$ docker build
```

## Configuration

To supply `scribe` with a [compendium](../README.md#glossary-of-terms), you'll need to mount a docker volume. By default, the docker image for `scribe` expects the compendium to exist at the following path: `/compendium.yaml`. Here is an example command for how to run scribe with a read-only compendium:

```
$ docker run -v /path/to/my/compendium.yaml:/compendium.yaml:ro scribe
```
