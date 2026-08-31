import assert from "node:assert/strict";
import fs from "node:fs";
import test from "node:test";

const read = (path) => fs.readFileSync(path, "utf8");

test("Registry owns one dependency and Make command graph", () => {
  const makefile = read("Makefile");
  const manifest = JSON.parse(read("package.json"));
  assert.equal(read(".node-version").trim(), manifest.engines.node);
  assert.match(manifest.packageManager, /^pnpm@\d+[.]\d+[.]\d+$/);
  assert.deepEqual(Object.keys(manifest.devDependencies), ["@soksak/soksak-spec"]);
  for (const target of ["preflight", "prepare", "build", "verify", "authenticate"]) {
    assert.match(makefile, new RegExp(`^${target}:`, "m"));
  }
  for (const duplicate of ["NODE_VERSION :=", "PNPM_VERSION :=", "SPEC_VERSION :="]) {
    assert.equal(makefile.includes(duplicate), false);
  }
});

test("preflight judges the repository-selected pnpm", () => {
  const preflight = read("scripts/check-build-environment.sh");
  assert.equal(/pnpm_executable|pnpmExecutable/.test(preflight), false);
  assert.match(preflight, /cd "\$root" && pnpm --version/);
});

test("Registry workflows use the repository-declared immutable validator and do not reinstall it", () => {
  for (const path of [".github/workflows/verify.yml", ".github/workflows/publish.yml"]) {
    const workflow = read(path);
    assert.match(workflow, /node-version-file: [.]node-version/);
    assert.match(workflow, /package_json_file: package[.]json/);
    assert.match(workflow, /make (?:verify|build|authenticate)/);
    for (const bypass of ["curl -fsSL https://github.com/soksak-ai/soksak-spec", ".dependency/spec/bin/validate.mjs", "PATH="]) {
      assert.equal(workflow.includes(bypass), false, `${path} bypasses package ownership through ${bypass}`);
    }
  }
});

test("Registry publication is explicit and a main push only verifies", () => {
  const publish = read(".github/workflows/publish.yml");
  const publishTrigger = publish.slice(publish.indexOf("on:\n") + "on:\n".length, publish.indexOf("concurrency:\n"));
  assert.equal(publishTrigger.trim(), "workflow_dispatch:");

  const verify = read(".github/workflows/verify.yml");
  const verifyTrigger = verify.slice(verify.indexOf("on:\n") + "on:\n".length, verify.indexOf("permissions:\n"));
  assert.match(verifyTrigger, /^  pull_request:\n    paths: \["plugins\/\*\*"\]$/m);
  assert.match(verifyTrigger, /^  push:\n    branches: \[main\]\n    paths: \["plugins\/\*\*"\]$/m);
  assert.match(verify, /- name: Require one plugin entry change\n        if: github[.]event_name == 'pull_request'/);
  for (const publicationSurface of ["actions/create-github-app-token", "registry-publish", "SOKSAK_REGISTRY_ED25519_SEED"]) {
    assert.equal(verify.includes(publicationSurface), false, `verify workflow exposes ${publicationSurface}`);
  }
});
