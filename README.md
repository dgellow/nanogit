[![Build Status](https://travis-ci.org/QRCLabs/nanogit.svg?branch=master)](https://travis-ci.org/QRCLabs/nanogit)
[![Coverage Status](https://coveralls.io/repos/github/QRCLabs/nanogit/badge.svg?branch=master)](https://coveralls.io/github/QRCLabs/nanogit?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/qrclabs/nanogit)](https://goreportcard.com/report/github.com/qrclabs/nanogit)

# nanogit - A lightweight git server with simple setup and configuration

**HIGHLY EXPERIMENTAL, DO NOT USE IN PRODUCTION**

## Introduction

`nanogit` is a small server to manage git projects.

It has following features:

- Organizations to group repositories and manage rights
- Teams, a group of user in an organization
- Entire config in one file, in a human readable format ([YAML](https://en.wikipedia.org/wiki/YAML))

## Install

### From `go get`

```
$ go get github.com/qrclabs/nanogit
```

## Usage

```
$ nanogit server

# If you want to use a different config file
$ nanogit server --config /path/to/custom/configfile.yml

# When the server is running you can begin to use git commands
$ git clone git@localhost:1337/QRCLabs/myproject.git
```

## Configuration

Server and access settings are configured in a [YAML](https://en.wikipedia.org/wiki/YAML) file. The default one is `config.yml` at the root of the project.

```yaml
server:
  port: 1337
  host: localhost
  root: /var/nanogit/
  user: nanogit
  group: nanogit

orgs:
  - id: fixme
    description: FIXME Hackerspace
    team:
      - name: default
        write: yes
        read: yes
      - name: ctf
        write: yes
        read: yes
      - name: comity
        write: yes
        read: yes
  - id: qrclabs
    description: QRC Labs company
    team:
      - name: default
        write: no
        read: yes
      - name: dev
        write: yes
        read: yes
      - name: admin
        write: no
        read: yes

users:
  - name: dgellow
    sshkeys:
      - from: https
        val: github.com/dgellow.keys
    orgs:
      - id: qrclabs
        teams:
          - dev
          - admin
      - id: fixme
  - name: notgcmalloc
    sshkeys:
      - from: hardcoded
        val: ssh-rsa AAAAB3NzaC1[truncated for the sake of readability]+MWYbwK1Tgx
      - from: file
        val: /path/to/file
    orgs:
      - id: fixme
        teams:
          - ctf

```
