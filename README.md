# Koanf zflag Provider

[![GoDoc](https://godoc.org/github.com/zulucmd/koanf-zflag?status.svg)](https://godoc.org/github.com/zulucmd/koanf-zflag)
[![Go Report Card](https://goreportcard.com/badge/github.com/zulucmd/koanf-zflag)](https://goreportcard.com/report/github.com/zulucmd/koanf-zflag)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/zulucmd/koanf-zflag?sort=semver)](https://github.com/zulucmd/koanf-zflag/releases)
[![Build Status](https://github.com/zulucmd/koanf-zflag/actions/workflows/validate.yml/badge.svg)](https://github.com/zulucmd/koanf-zflag/actions/workflows/validate.yml)

<!-- toc -->

- [Introduction](#introduction)
- [Installation](#installation)
- [Documentation](#documentation)
  - [Reading from command line](#reading-from-command-line)

<!-- /toc -->

## Introduction

koanf-zflag is a Koanf provider to retrieve configuration from [zulucmd/zflag](https://github.com/zulucmd/zflag).

## Installation

koanf-zflag is available using the standard `go get` command.

Install by running:

```bash
go get github.com/zulucmd/koanf-zflag
```

## Documentation

### Reading from command line

The following example shows the use of kzflag.Provider, a wrapper over the spf13/pflag library, an advanced commandline
lib. For Go's built in flag package, use basicflag.Provider.

<!-- BEGIN EMBED FILE: example_test.go -->
<!-- END EMBED_FILE: example_test.go -->
