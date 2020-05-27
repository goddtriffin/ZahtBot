# ZahtBot

The Offical ZahtBot.

## Setup

In your shell, run: `export ZAHT_BOT_TOKEN=<token>`. Replace `<token>` with your actual Discord bot token.

## Usage

- `make usage`
  - display Makefile target info
- `make buildlocal`
  - builds the binary locally
- `make runlocal`
  - runs the binary locally
- `make builddocker`
  - builds the binary and Docker container
- `make rundocker`
  - creates and runs a new Docker container
- `make startdocker`
  - resumes a stopped Docker container
- `make stopdocker`
  - stops the Docker container
- `make removedocker`
  - removes the Docker container
- `make memusage`
  - displays the memory usage of the currently running Docker container

## DCA File Format

[https://github.com/bwmarrin/dca](https://github.com/bwmarrin/dca)
