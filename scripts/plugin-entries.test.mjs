import assert from "node:assert/strict";
import { mkdtempSync, readdirSync, readFileSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { parseDependencyIntent } from "@soksak-ai/plugin-spec";

const root = join(import.meta.dirname, "..");

// A registry entry is a dependency intent { id, version } named by its id. The release reference
// { id, version, size, sha256 } is derived by the publish workflow from the release.json at the
// location derived from the id and version; an entry that carries any other key is refused.
const entryErrors = (directory) =>
  readdirSync(directory).sort().flatMap((name) => {
    const intent = parseDependencyIntent(JSON.parse(readFileSync(join(directory, name), "utf8")), `plugins/${name}`);
    if (!intent.ok) return intent.errors;
    return name === `${intent.value.id}.json` ? [] : [`plugins/${name}: file name must be ${intent.value.id}.json`];
  });

test("every plugins/<id>.json is a dependency intent named by its id", () => {
  assert.deepEqual(entryErrors(join(root, "plugins")), []);
});

test("an entry carrying url, size, and sha256 is refused", () => {
  const directory = mkdtempSync(join(tmpdir(), "registry-entries-"));
  writeFileSync(join(directory, "soksak-plugin-example.json"), `${JSON.stringify({
    id: "soksak-plugin-example",
    version: "1.2.3",
    url: "https://github.com/soksak-ai/soksak-plugin-example/releases/download/v1.2.3/release.json",
    size: 1234,
    sha256: "a".repeat(64),
  }, null, 2)}\n`);
  const errors = entryErrors(directory);
  assert.equal(errors.length, 3, errors.join("; "));
  for (const key of ["url", "size", "sha256"]) assert.ok(errors.some((error) => error.includes(`"${key}"`)), `${key}: ${errors.join("; ")}`);
});

test("an entry whose file name differs from its id is refused", () => {
  const directory = mkdtempSync(join(tmpdir(), "registry-entries-"));
  writeFileSync(join(directory, "example.json"), `${JSON.stringify({ id: "soksak-plugin-example", version: "1.2.3" })}\n`);
  assert.deepEqual(entryErrors(directory), ["plugins/example.json: file name must be soksak-plugin-example.json"]);
});
