Babycron
========

Run a single cron job in the foreground, outputting sent to stdout. Works
nicely in Docker.

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
correct permissions in your script. Additionally, Babycron will find program in
your `PATH`, so usually `bash` or `sh` is find.


In Docker:

```docker
ENTRYPOINT [ "babycron", "*/1 * * * *", "sh path/to/script.sh" ]
```

All output and errors are sent to stdout, so they'll appear in your Docker logs.


## Installation

macOS and Linux users can use Homebrew:

```bash
brew install meowgoritm/tap/babycron
```

Our you can download a binary from the [releases](https://github.com/meowgoritm/babycron/releases) page.

Or just use `go get`:

```bash
go get github.com/meowgoritm/babycron
```

## License

MIT
