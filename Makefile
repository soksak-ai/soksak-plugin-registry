SHELL := /bin/sh
.PHONY: preflight guard prepare verify build authenticate
registry_flags = --@soksak:registry=$(REGISTRY) --@soksak-ai:registry=$(REGISTRY) --config.minimum-release-age=0
# REGISTRY is accepted from the make command line only ($(origin) must be "command line").
# GNU make's own environment channels (MAKEFLAGS, GNUMAKEFLAGS, MAKEFILES, -e) are outside this
# Makefile's control and are not refused; setting them is a deliberate act of the caller.
preflight:
	@scripts/check-build-environment.sh
# A package that depends on @soksak/* or @soksak-ai/* requires REGISTRY for every install, the public registry included.
guard:
	@case "$(origin REGISTRY)" in undefined|"command line") ;; *) echo 'REGISTRY from the $(origin REGISTRY) is refused: make verify REGISTRY=http://host:port/' >&2; exit 64 ;; esac
	@case "$(origin REGISTRY):$(REGISTRY)" in undefined:|"command line:http://"*|"command line:https://"*) ;; *) echo 'REGISTRY must be an absolute URL: make verify REGISTRY=http://host:port/' >&2; exit 64 ;; esac
	@dependency=$$(node -p 'const p=require("$(CURDIR)/package.json");Object.keys({...p.dependencies,...p.devDependencies,...p.peerDependencies}).find((name)=>/^@soksak(-ai)?\//.test(name))??""') || exit $$?; test -z "$$dependency" || test "$(origin REGISTRY)" = "command line" || { echo "REGISTRY required: this package depends on $$dependency: make verify REGISTRY=http://host:port/" >&2; exit 64; }
prepare: guard preflight
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm install --frozen-lockfile $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags))
# pnpm 11 compares the settings recorded by the install before every exec and reinstalls without
# the registry flags on any difference (CI toggles enableGlobalVirtualStore, the flags set
# minimumReleaseAge). The exec commands repeat the install environment and flags; the flags precede
# exec so pnpm consumes them instead of the spec command.
# plugins/<id>.json is an intent { id, version }. registry-verify and registry-build derive each
# release.json location from the id and version, read it, and fill size and sha256 from the bytes.
verify: prepare
	@node --test scripts/*.test.mjs
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) exec soksak-validate registry-verify plugins
build: prepare
	@test -n '$(REGISTRY_SEQUENCE)' && test -n '$(REGISTRY_ISSUED_AT)' && test -n '$(REGISTRY_EXPIRES_AT)' && test -n '$(REGISTRY_OUT)' || { echo 'REGISTRY_SEQUENCE, REGISTRY_ISSUED_AT, REGISTRY_EXPIRES_AT and REGISTRY_OUT are required' >&2; exit 2; }
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) exec soksak-validate registry-build plugins --id official --sequence '$(REGISTRY_SEQUENCE)' --issued-at '$(REGISTRY_ISSUED_AT)' --expires-at '$(REGISTRY_EXPIRES_AT)' --out '$(REGISTRY_OUT)'
authenticate: prepare
	@test -n "$$SOKSAK_REGISTRY_ED25519_SEED" && test -n '$(REGISTRY_IN)' && test -n '$(REGISTRY_OUT)' || { echo 'SOKSAK_REGISTRY_ED25519_SEED, REGISTRY_IN and REGISTRY_OUT are required' >&2; exit 2; }
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) exec soksak-validate registry-authenticate '$(REGISTRY_IN)' --out '$(REGISTRY_OUT)'
