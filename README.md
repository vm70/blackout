# Blackout - A Blackout / Erasure Poem Generator

Blackout is a command-line application that removes letters from a public-domain
poem to reveal a specific message. It's based on the article
["Using Regular Expressions to Make Blackout Poetry"](regex-blackout) by Vincent
'VM' Mercator and is written in Go.

## Installation

You can install Blackout through Go itself using the `go install` command.

```bash
go install github.com/vm70/blackout@latest
```

For UNIX users, make sure that your `$PATH` environment variable contains the
`$GOBIN` path for Go-installed binaries (with the default being `~/go/bin/`).

## Building From Source

You can also download Blackout and build it directly from its source code.

```bash
git clone https://github.com/vm70/blackout.git
cd blackout
go build .
./blackout --help
```

## Usage

When given an input message (e.g., `blackout poem`), Blackout will return a
public-domain poem blacked out to spell it.

```
[user@pc]$ ./blackout 'blackout poem'

███
█ b████ ████ █l████ a████ █████
██ █████ c███ ████████ ███ ████
███████ ███ ██████ ████ ███████ ██ ███ ████████
████ ████ ██████ ██ ██████ ███ ████████
████████████ ████ ████ ██ ██ ██████
████ █████ ██ ██ ███ █████ ██████
███ ███k ██ █o██ ██ ██u██ █████
███ ███ ███t █████ ████ ███ p██████
███ █o██ ██ ████ ██e ██m██ ████████
██ ██████ ████ ███████ ███ ███████

blackout poem
Excerpt of "The Bird Wounded By An Arrow." by Jean de La Fontaine
```

Running `blackout --help` or `blackout -h` will return the following help
message.

```text
Make a blackout poem with the given hidden message

Usage:
  blackout <message> [flags]

Flags:
  -h, --help             help for blackout
  -l, --max-length int   maximum poem length (default 400)
  -p, --print-original   print original poem before blacking out
  -V, --verbose          verbose output
  -v, --version          version for blackout
```

## Contributing

Contributions are welcome. If you find a bug, please report it through
Blackout's [Issues page](issues) on its GitHub repository.

## License

> Copyright © 2024 Vincent Mercator.
>
> Licensed under the Apache License, Version 2.0 (the "License"); you may not
> use this file except in compliance with the License. You may obtain a copy of
> the License at <http://www.apache.org/licenses/LICENSE-2.0>.

[regex-blackout]: (https://vm70.neocities.org/posts/2024-05-11-regex-blackout/)
[issues]: https://github.com/vm70/blackout/issues
