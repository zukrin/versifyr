# versifyr

Many times we need to manage the version of a project in many files. This tool can set the version of a project in project files.

Versifyr allows to set values as defined from key=value pairs passed at command line.

Files to be managed are listed in the configuration file `.versifyr/configuration.json` in the root of the project:

```yaml

files:
  - name: chart.yaml
    type: yaml
    path: chart/Chart.yaml
  - name: pom.xml
    type: xml
    path: pom.xml
  - name: Version.java
    type: java
    path: src/main/java/sample/Version.java

```

The path is relative to the root of the project. Each file may contain well commented lines to identify the version to be managed
in the form `$versifyr:template=<template>$`. The template will replace the followin line:

```yaml

apiVersion: v2
name: orchestrator
description: A Helm chart for Kubernetes

# A chart can be either an 'application' or a 'library' chart.
#
# Application charts are a collection of templates that can be packaged into versioned archives
# to be deployed.
#
# Library charts provide useful utilities or functions for the chart developer. They're included as
# a dependency of application charts to inject those utilities and functions into the rendering
# pipeline. Library charts do not define any templates and therefore cannot be deployed.
type: application

# This is the chart version. This version number should be incremented each time you make changes
# to the chart and its templates, including the app version.
# Versions are expected to follow Semantic Versioning (https://semver.org/)
# $versifyr:template=version: {{.version}}$
version: 4.0.1

# This is the version number of the application being deployed. This version number should be
# incremented each time you make changes to the application. Versions are not expected to
# follow Semantic Versioning. They should reflect the version the application is using.
# $versifyr:template=appVersion: {{.version}}$
appVersion: 4.0.1

```

Files can be any text file. `versifyr` can set any value passed at command line.

Some values are always available:

`actualdate`: es. 2023-05-22 
`actualtime`: es. 16:15:57 
`actualtimestamp`: es. 2023-05-22 16:15:57 
`latesttag`: es. 4.0.1, in case of git repository


Se examples in `examples` folder.

Usage

```sh

NAME:
   versifyr - A new cli application

USAGE:
   versifyr [global options] command [command options] [arguments...]

VERSION:
   2.0.1

AUTHOR:
   Stefano Zuccaro <zukrin@gmail.com>

COMMANDS:
   init, i  init project configuration
   show, s  show actual configuration and file content
   set, s   set values as key=value to be replaced in files
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d     set output to debug (default: false)
   --nochange, -n  simulate changes (default: false)
   --help, -h      show help
   --version, -v   print the version

COPYRIGHT:
   (c) 2023 Stefano Zuccaro


```

## installation

```sh
>  go install -v github.com/zukrin/versifyr/cmd/versifyr@latest
```

## init

Use `versifyr init` to create the configuration file `.versifyr/configuration.json` in the root of the project

## show

Use `versifyr show` to show the actual configuration and the file content:

```sh
versifyr show
      1 actual situation
      ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

      1.1 version.go

      ┃ package versifyr
      ┃ 
      ┃ // $versifyr:template=const Version = "{{ .version }}"$
      ┃ const Version = "v0.0.1"
      ┃ 
      ┃ // $versifyr:template=const Sample = "{{ .sample }}"$
      ┃ const Sample = "something"
      ┃ 
      ┃ // $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
      ┃ const Compiled = "2023-05-23 10:29:22"

```

## set

Use `versifyr set` to set values as key=value to be replaced in files. Use the option `-n` to simulate changes.

```sh

> versifyr -n set version="v0.0.1" sample="something" 
setting values
using values map[actualdate:2023-05-23 actualtime:10:28:58 actualtimestamp:2023-05-23 10:28:58 latesttag:unknown sample:something version:v0.0.1]
replaced into version.go line 3 with const Version = "v0.0.1"
replaced into version.go line 6 with const Sample = "something"
replaced into version.go line 9 with const Compiled = "2023-05-23 10:28:58"
      1 transformed files
      ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

      1.1 version.go

      ┃ package versifyr
      ┃ 
      ┃ // $versifyr:template=const Version = "{{ .version }}"$
      ┃ const Version = "v0.0.1"
      ┃ 
      ┃ // $versifyr:template=const Sample = "{{ .sample }}"$
      ┃ const Sample = "something"
      ┃ 
      ┃ // $versifyr:template=const Compiled = "{{ .actualtimestamp }}"$
      ┃ const Compiled = "2023-05-23 10:28:58"


```sh