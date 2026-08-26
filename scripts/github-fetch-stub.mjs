// Test preload (node --import): replaces the global fetch with a reader of a release directory
// mirror under the working directory. GET https://github.com/soksak-ai/<id>/releases/download/v<version>/<file>
// is served from <cwd>/github/<id>/<version>/<file>; a missing file answers 404. Every requested
// URL is appended to <cwd>/fetch.log so a test can assert which locations the tool derived.
import { appendFileSync, readFileSync } from "node:fs";
import { join } from "node:path";

const RELEASE_URL_RE = /^https:\/\/github\.com\/soksak-ai\/([^/]+)\/releases\/download\/v([^/]+)\/([^/]+)$/;

globalThis.fetch = async (input) => {
  const url = String(input);
  appendFileSync(join(process.cwd(), "fetch.log"), `${url}\n`);
  const match = RELEASE_URL_RE.exec(url);
  if (!match) throw new Error(`fetch stub: unexpected URL ${url}`);
  let bytes;
  try {
    bytes = readFileSync(join(process.cwd(), "github", match[1], match[2], match[3]));
  } catch {
    return new Response(null, { status: 404 });
  }
  return new Response(bytes, { status: 200, headers: { "content-length": String(bytes.length) } });
};
