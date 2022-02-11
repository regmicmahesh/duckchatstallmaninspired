<p align="center">
<img src="https://user-images.githubusercontent.com/64764773/153582967-b7c7fa11-a91d-4f40-9521-48dc9432dfe2.png" width="100" align="center">
</p>
<h1 align="center">DuckChat - Stallman Inspired </h1>

[![Latest Release](https://github.com/regmicmahesh/duckchatstallmaninspired/actions/workflows/pull-request-merge.yml/badge.svg?branch=main)](https://github.com/regmicmahesh/duckchatstallmaninspired/actions/workflows/pull-request-merge.yml)

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
