# Runtime-env

Runtime-env is cli written in go to parse environment variables & generate a javascript file which adds these to the browsers window object.

[![Release](https://img.shields.io/github/release/shkreios/runtime-env.svg?style=for-the-badge)](https://github.com/shkreios/runtime-env/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/shkreios/runtime-env/build?style=for-the-badge)](https://github.com/shkreios/runtime-env/actions?workflow=publish)

[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)

# Install

## Script install

```bash
curl -sfL https://raw.githubusercontent.com/shkreios/runtime-env/main/install.sh | sh
```

## NPM install

### NPM

```bash
npm install --dev-save runtime-env
```

### YARN

```bash
yarn add -D runtime-env
```

# Usage

```
NAME:
   runtime-env - runtime envs for SPAs

USAGE:
   runtime-env [global options]

VERSION:
   v1.0.0

AUTHOR:
   Simon Hessel <simon.hessel@kreios.lu>

GLOBAL OPTIONS:
   --env-file value, -f value                   The .env file to be parsed
   --prefix value, -p value                     The env prefix to matched
   --output value, -o value                     Output file path (default: "./env.js")
   --type-declarations-file value, --dts value  Output file path for the typescript declaration file
   --global-key value, --key value              Customize the key on which the envs will be set on window object
   --remove-prefix                              Remove the prefix from the env (default: false)
   --no-envs                                    Only read envs from file not from environment variables (default: false)
   --disable-logs, --no-logs                    Disable logging output (default: false)
   --help, -h                                   show help (default: false)
   --version, -v                                print the version (default: false)

COPYRIGHT:
   Copyright © 2020 Simon Hessel
```

## Input

```sh
# .env input file
TEST=Test
```

## Command

```sh
$ runtime-env -f .env --no-envs
```

## Output

```js
window.__RUNTIME_CONFIG__ = { TEST: "Test" };
```

## Code

```js
console.log(window.__RUNTIME_CONFIG__.TEST);
```

## Docker

```sh
#!/bin/sh

...

# any other entrypoint operations here
....
runtime-env # your runtime-env flags here

# any other entrypoint operations here
...

# call CMD
exec $@
```

```dockerfile
...

# install runtime-env via install.sh script
RUN wget -O - https://raw.githubusercontent.com/shkreios/runtime-env/main/install.sh | sh

# important to be set after install.sh as otherwise binary will placed under /app/bin/runtime-env
WORKDIR /app/

# Copy your entrypoint script
COPY entrypoint.sh entrypoint.sh

...

#
ENTRYPOINT ["./entrypoint.sh"]

# optionaly set command
# CMD ["nginx"]

...
```

## In CI/CD

```sh

# use either wget, curl or whatevery you ci env supports

curl -sfL https://raw.githubusercontent.com/shkreios/runtime-env/main/install.sh | sh

# wget -O - https://raw.githubusercontent.com/shkreios/runtime-env/main/install.sh | sh

runtime-env # your runtime-env flags here
```

### [Github Action](https://github.com/shkreios/runtime-env-action)

#### Read from secrets

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@master
      - uses: shkreios/runtime-env-action@v1
        with:
          version: v1.1.0
          prefix: RUN_ENV_
          output: ./public/env.js
          removePrefix: "true"
        env:
          RUN_ENV_EXAMPLE: ${{ secrets.EXAMPLE }}
```

```js
// resulting env.js in public folder
window.__RUNTIME_CONFIG__ = { EXAMPLE: "SECRET_VALUE" };
```

# Why go / Why should i use this package?

There many solutions similar to this package `react-env`, `runtime-env-cra` or `sh` scripts to name a few. They all fall in one of 2 categories:

1. they are either dependent on a runtime/interpreter to be preinstalled like nodejs or python or
2. they are not platform-agnostic

Go supports building binaries for multiple platforms and arches, and the binary itself includes everything to be executed. Therefore, no matter if development with npm or in docker CI/CD you can expect it to be easily installed and as light as it can be.

# Why have a npm wrapper/placeholder script?

As this tool will mostly be used in npm codebases where you expect everything, including a tool that generates runtime-envs, to be installed via npm/yarn it's the best option.
