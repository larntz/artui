# ArTUI

[![build](https://github.com/larntz/artui/actions/workflows/build-release.yaml/badge.svg)](https://github.com/larntz/artui/actions/workflows/build-release.yaml)
[![build](https://github.com/larntz/artui/actions/workflows/build.yaml/badge.svg)](https://github.com/larntz/artui/actions/workflows/build.yaml)

## Overview 

This is a TUI app for interacting with ArgoCD and managing apps.

![screenshot](screenshots/artui.png)

## Commands

- `:r` refresh app
- `:hr` hard refresh app
- `:s` sync app
- `:q` quit

## Navigation

- `j`,`k` up/down application list
- up arrow, down arrow scroll application overview
- `/` search application list

## Compatility

This is currently being developed against ArogCD version 2.5.4. 

## Dependencies

Go v1.18.5.

To setup dependencies use the replace() block from the [argocd project `go.mod`](https://github.com/argoproj/argo-cd/blob/v2.5.4/go.mod) file. Then run `go mod tidy`.
