# Blackout - A Blackout / Erasure Poem Generator

Blackout is a command-line application that automates the process of making
simple blackout poems, where characters and words of an original text source are
removed to create an entirely new piece. It combs through a database of
public-domain poetry to find one with characters that match a given message,
then prints the resulting blacked-out poem to standard output.

Blackout is based on the article
["Using Regular Expressions to Make Blackout Poetry"](regex-blackout) by Vincent
Mercator and is written in Go.

## Installation

You can install Blackout through Go itself using the `go install` command.

```bash
# install the latest stable version
go install github.com/vm70/blackout@latest
# install a specific tagged version / branch
go install github.com/vm70/blackout@v0.3.0-alpha.0
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

```text
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
  -p, --allow-profanities   Allow blacking out poems with profanities
  -h, --help                help for blackout
  -l, --max-length int      maximum poem length (default 400)
  -o, --print-original      print original poem before blacking out
  -V, --verbose             verbose output
  -v, --version             version for blackout

```

## Contributing

Contributions are welcome. If you find a bug, please report it through
Blackout's [Issues page](issues) on its GitHub repository.

## Special Thanks

- HuggingFace user [`DanFosing`](DanFosing) and the
  [`public-domain-poetry`](public-domain-poetry) dataset
- [Puttock International](https://pi01.net/), the owner(s) of the
  [Public Domain Poetry](pdp) website

## License

The poems downloaded and stored by this program are in the public domain, and
can be viewed either at the [`public-domain-poetry`](public-domain-poetry)
dataset page on HuggingFace or the [Public Domain Poetry](pdp) website.

The code in this repository uses the Apache 2.0 license. For more information,
see [LICENSE](https://github.com/vm70/blackout/blob/main/LICENSE).

> Copyright © 2024 Vincent Mercator.
>
> Licensed under the Apache License, Version 2.0 (the "License"); you may not
> use this file except in compliance with the License. You may obtain a copy of
> the License at <http://www.apache.org/licenses/LICENSE-2.0>.

[pdp]: https://www.public-domain-poetry.com/
[DanFosing]: https://huggingface.co/DanFosing
[issues]: https://github.com/vm70/blackout/issues
[public-domain-poetry]:
  https://huggingface.co/datasets/DanFosing/public-domain-poetry
[regex-blackout]: https://vm70.neocities.org/posts/2024-05-11-regex-blackout/
