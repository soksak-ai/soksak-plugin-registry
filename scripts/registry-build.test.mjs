import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { createHash } from "node:crypto";
import { existsSync, mkdirSync, mkdtempSync, readFileSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { pathToFileURL } from "node:url";

const root = join(import.meta.dirname, "..");
const validate = join(root, "node_modules/@soksak/soksak-spec/bin/validate.mjs");
const stub = pathToFileURL(join(root, "scripts/github-fetch-stub.mjs")).href;

const sha256 = (bytes) => createHash("sha256").update(bytes).digest("hex");
const document = (value) => Buffer.from(`${JSON.stringify(value, null, 2)}\n`);
const integrity = (file) => ({ file, size: file.length, sha256: sha256(file) });
const commit = createHash("sha1").update("fixture").digest("hex");

// The sidecar release.json is written first: the plugin release references its bytes.
const sidecar = document({
  kind: "sidecar",
  id: "soksak-sidecar-example",
  version: "0.1.0",
  manifest: integrity("sidecar.json"),
  source: { repository: "https://github.com/soksak-ai/soksak-sidecar-example", commit },
  artifacts: [{ target: "aarch64-apple-darwin", format: "tar.gz", manifest: "sidecar.json", ...integrity("soksak-sidecar-example-0.1.0-aarch64-apple-darwin.tar.gz") }],
  evidence: [integrity("conformance-sidecar.json")],
});
const plugin = document({
  kind: "plugin",
  id: "soksak-plugin-example",
  version: "1.2.3",
  manifest: integrity("plugin.json"),
  source: { repository: "https://github.com/soksak-ai/soksak-plugin-example", commit },
  artifacts: [{ target: "any", format: "tgz", manifest: "plugin.json", ...integrity("soksak-plugin-example-1.2.3-any.tgz") }],
  runtimeDependencies: { sidecars: [{ id: "soksak-sidecar-example", version: "0.1.0", size: sidecar.length, sha256: sha256(sidecar) }] },
  evidence: [integrity("conformance-plugin.json")],
});

const releaseURL = (id, version) => `https://github.com/soksak-ai/${id}/releases/download/v${version}/release.json`;

const fixture = () => {
  const cwd = mkdtempSync(join(tmpdir(), "registry-build-"));
  for (const [id, version, bytes] of [["soksak-plugin-example", "1.2.3", plugin], ["soksak-sidecar-example", "0.1.0", sidecar]]) {
    mkdirSync(join(cwd, "github", id, version), { recursive: true });
    writeFileSync(join(cwd, "github", id, version, "release.json"), bytes);
  }
  mkdirSync(join(cwd, "plugins"));
  return cwd;
};
const build = (cwd) =>
  spawnSync(process.execPath, ["--import", stub, validate, "registry-build", "plugins", "--id", "official", "--sequence", "11", "--issued-at", "2026-08-26T00:00:00Z", "--expires-at", "2026-11-24T00:00:00Z", "--out", "unsigned.json"], { cwd, encoding: "utf8" });

test("registry-build derives the release location from {id, version} and fills size and sha256 from the fetched release.json", () => {
  const cwd = fixture();
  writeFileSync(join(cwd, "plugins/soksak-plugin-example.json"), document({ id: "soksak-plugin-example", version: "1.2.3" }));
  const result = build(cwd);
  assert.equal(result.status, 0, result.stderr);
  const registry = JSON.parse(readFileSync(join(cwd, "unsigned.json"), "utf8"));
  assert.deepEqual(Object.keys(registry).sort(), ["expiresAt", "id", "issuedAt", "plugins", "sequence"]);
  assert.deepEqual(registry.plugins, [{ id: "soksak-plugin-example", version: "1.2.3", size: plugin.length, sha256: sha256(plugin) }]);
  assert.deepEqual(readFileSync(join(cwd, "fetch.log"), "utf8").trim().split("\n"), [
    releaseURL("soksak-plugin-example", "1.2.3"),
    releaseURL("soksak-sidecar-example", "0.1.0"),
  ]);
});

test("registry-build refuses an entry that carries url, size, and sha256 before any fetch", () => {
  const cwd = fixture();
  writeFileSync(join(cwd, "plugins/soksak-plugin-example.json"), document({ id: "soksak-plugin-example", version: "1.2.3", url: releaseURL("soksak-plugin-example", "1.2.3"), size: plugin.length, sha256: sha256(plugin) }));
  const result = build(cwd);
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /registry entry soksak-plugin-example\.json: .*"url"/);
  assert.equal(existsSync(join(cwd, "fetch.log")), false);
  assert.equal(existsSync(join(cwd, "unsigned.json")), false);
});

test("registry-build refuses an entry whose release.json is absent at the derived location", () => {
  const cwd = fixture();
  writeFileSync(join(cwd, "plugins/soksak-plugin-example.json"), document({ id: "soksak-plugin-example", version: "9.9.9" }));
  const result = build(cwd);
  assert.notEqual(result.status, 0);
  assert.match(result.stderr, /unresolved release soksak-plugin-example@9\.9\.9: https:\/\/github\.com\/soksak-ai\/soksak-plugin-example\/releases\/download\/v9\.9\.9\/release\.json/);
  assert.equal(existsSync(join(cwd, "unsigned.json")), false);
});
