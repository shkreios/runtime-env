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

# Example

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

## Usage

```js
console.log(window.__RUNTIME_CONFIG__.TEST);
```
