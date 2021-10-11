const { readFileSync } = require("fs");
const { extract } = require("tar");
const { dirname, resolve } = require("path");
const { spawnSync } = require("child_process");
const gunzip = require("gunzip-maybe");
const fetch = (...args) =>
  import("node-fetch").then(({ default: fetch }) => fetch(...args));

const errorFn = (msg) => {
  console.error(msg);
  process.exit(1);
};

const validateConfiguration = (packageJson) => {
  if (!packageJson.version) return "'version' property must be specified";

  if (!packageJson.binary || typeof packageJson.binary !== "object")
    return "'binary' property must be defined and be an object";

  if (!packageJson.binary.name) return "'name' property is necessary";

  if (!packageJson.binary.url) return "'url' property is required";

  return undefined;
};

const ARCH_MAPPING = {
  ia32: "386",
  x64: "amd64",
  arm: "arm",
};

const PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
  freebsd: "freebsd",
};

function getPlatformMetadata() {
  if (!(process.arch in ARCH_MAPPING))
    errorFn(
      "Installation is not supported for this architecture: " + process.arch
    );

  if (!(process.platform in PLATFORM_MAPPING))
    errorFn(
      "Installation is not supported for this platform: " + process.platform
    );

  const packageJson = JSON.parse(
    readFileSync(resolve(__dirname, "..", "package.json"))
  );
  const error = validateConfiguration(packageJson);
  if (error && error.length > 0) errorFn("Invalid package.json: " + error);

  // We have validated the config. It exists in all its glory

  let { name, url } = packageJson.binary;
  let { version } = packageJson;
  if (version[0] === "v") version = version.substr(1); // strip the 'v' if necessary v0.0.1 => 0.0.1

  // Binary name on Windows has .exe suffix
  if (process.platform === "win32") name += ".exe";

  // Interpolate variables in URL, if necessary
  url = url.replace(/{{arch}}/g, ARCH_MAPPING[process.arch]);
  url = url.replace(/{{platform}}/g, PLATFORM_MAPPING[process.platform]);
  url = url.replace(/{{version}}/g, version);
  url = url.replace(/{{bin_name}}/g, name);

  return {
    name,
    url,
  };
}

const download = async (url, installDirectory) => {
  const res = await fetch(url);

  return new Promise((resolve, reject) => {
    res.body
      .pipe(gunzip())
      .pipe(extract({ cwd: installDirectory }))
      .on("end", resolve)
      .on("error", reject);
  });
};

const run = async () => {
  const { name, url } = getPlatformMetadata();
  const [, binLocation, ...args] = process.argv;

  const binFolder = dirname(binLocation);
  await download(url, binFolder);

  const options = { cwd: process.cwd(), stdio: "inherit" };
  const result = spawnSync(resolve(binFolder, name), args, options);

  if (result.error) error(result.error);

  process.exit(result.status);
};

module.exports = run;
