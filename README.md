Babycron
========

[![Build Status](https://github.com/meowgorithm/babycron/workflows/build/badge.svg)](https://github.com/meowgorithm/babycron/actions)

Run a single cron job in the foreground, sending output to stdout. Works well
in Docker.

## Usage

Pass two arguments: the cron schedule and a string of the task to run.

```bash
# Run a script every minute
babycron '*/1 * * * *' 'sh path/to/script.sh'

# Run a script on start, and then again every six hours
babycron --run-on-start '* */6 * * *' 'sh path/to/script.sh'

# Gzip and backup a Redis DB every six hours
babycron '* */6 * * *' 'cat /data/dump.rdb | gzip | pipedream -b backups -p backup.rdb.gz'
```

Note that if you’re running a script you *must* include the interpreter in the
second argument (i.e. `sh` or `/bin/sh`), regardless if you have a `#!` and
executable permissions. Babycron will search your `PATH`, so things like `bash`
or `sh`, without the full path to the interpreter, is usually fine.

In Docker:

```docker
ENTRYPOINT [ "babycron", "*/1 * * * *", "sh path/to/script.sh" ]
```

Output and errors are sent to stdout, so they'll appear in your Docker logs.

## Installation

macOS and Linux users can use Homebrew:

```bash
brew install meowgoritm/tap/babycron
```

Additional binaries can be found on the [releases](https://github.com/meowgorithm/babycron/releases) page (Linux x86_64/ARM, macOS x86_64, Windows x86_64/i386).

Or just use `go get`:

```bash
go get github.com/meowgoritm/babycron
```

## License

MIT
