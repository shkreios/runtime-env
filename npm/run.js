#!/usr/bin/env node
const {
  Binary,
  GO_ARCH_MAPPING,
  GO_PLATFORM_MAPPING,
} = require("binary-cli-install");

const { join } = require("path");
const packageJson = require(join("..", "package.json"));

const bin = new Binary(
  packageJson,
  GO_ARCH_MAPPING,
  GO_PLATFORM_MAPPING,
  false
);
bin.run();
