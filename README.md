# soksak-plugin-registry

The official soksak plugin registry — the single source of truth for the app's "installable" list.

`registry.json` holds each plugin's display metadata plus its git repo URL (the install source). The
soksak app shows plugins instantly from a build snapshot, and refreshes from this file once per session
or on demand.

## Schema

```json
{
  "spec": "soksak-registry@0.0.1",
  "plugins": [
    {
      "id": "soksak-plugin-shark",
      "name": "shork shark",
      "version": "1.0.2",
      "description": "...",
      "repo": "https://github.com/soksak-ai/soksak-plugin-shark.git"
    }
  ]
}
```

`name`/`description` is a string or a `{ "ko": …, "en": … }` localized object. `repo` is any git URL
(github/gitlab/self-host). Actual install clones the repo and the app strictly re-validates the manifest.

## Registering

Register a new plugin with a PR that adds an entry to `registry.json`. Each plugin is developed and
versioned independently in its own git repo (multiple authors).
