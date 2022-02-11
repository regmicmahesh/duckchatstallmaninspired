<p align="center">
<img src="https://user-images.githubusercontent.com/64764773/153582967-b7c7fa11-a91d-4f40-9521-48dc9432dfe2.png" width="100" align="center">
</p>
<h1 align="center">DuckChat - Stallman Inspired </h1>

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)]() [![Latest Release](https://github.com/regmicmahesh/duckchatstallmaninspired/actions/workflows/pull-request-merge.yml/badge.svg?branch=main)](https://github.com/regmicmahesh/duckchatstallmaninspired/actions/workflows/pull-request-merge.yml) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/regmicmahesh/duckchatstallmaninspired/graphs/commit-activity)

## Demo of current release.

![](./screenshot.gif)

This application runs completely on top of TCP , and we don't bother with high level protocols such as http, websockets and others. We've a simple run and we don't bother with message segments like request header, request body, crlf and stuffs.

We run a simple protocol that is your message ends with `\n`. How simple is that?

For now, there are few commands such as 

| Command | Description | Example | 
| --- | --- | --- |
| `/join` | Join? | `/join` | 
| `/quit` | Quit? | `/quit` | 
| `/whisper` | DM? | `/whisper mahesh k xa yar` | 
| `/nick` | Change Nickname? | `/nick mahesh` |
