import assert from "node:assert/strict";
import fs from "node:fs";
import test from "node:test";

const read = (path) => fs.readFileSync(path, "utf8");

test("Registry owns one dependency and Make command graph", () => {
  const makefile = read("Makefile");
  const manifest = JSON.parse(read("package.json"));
  assert.equal(read(".node-version").trim(), manifest.engines.node);
  assert.match(manifest.packageManager, /^pnpm@\d+[.]\d+[.]\d+$/);
  assert.deepEqual(Object.keys(manifest.devDependencies), ["@soksak-ai/plugin-spec"]);
  for (const target of ["preflight", "prepare", "build", "verify", "authenticate"]) {
    assert.match(makefile, new RegExp(`^${target}:`, "m"));
  }
  for (const duplicate of ["NODE_VERSION :=", "PNPM_VERSION :=", "SPEC_VERSION :="]) {
    assert.equal(makefile.includes(duplicate), false);
  }
});

test("Registry workflows inject owners and do not reinstall the spec manually", () => {
  for (const path of [".github/workflows/verify.yml", ".github/workflows/publish.yml"]) {
    const workflow = read(path);
    assert.match(workflow, /node-version-file: [.]node-version/);
    assert.match(workflow, /package_json_file: package[.]json/);
    assert.match(workflow, /make (?:verify|build|authenticate)/);
    for (const bypass of ["curl -fsSL https://github.com/soksak-ai/soksak-spec", ".dependency/spec/bin/validate.mjs"]) {
      assert.equal(workflow.includes(bypass), false, `${path} bypasses package ownership through ${bypass}`);
    }
  }
});
