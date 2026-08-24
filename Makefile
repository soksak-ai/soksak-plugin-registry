SHELL := /bin/sh

.PHONY: preflight prepare verify build authenticate

preflight:
	@scripts/check-build-environment.sh

prepare: preflight
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm install --frozen-lockfile

verify: prepare
	@node --test scripts/*.test.mjs
	@pnpm exec soksak-validate registry-verify plugins

build: prepare
	@test -n '$(REGISTRY_SEQUENCE)' && test -n '$(REGISTRY_ISSUED_AT)' && test -n '$(REGISTRY_EXPIRES_AT)' && test -n '$(REGISTRY_OUT)' || { echo 'REGISTRY_SEQUENCE, REGISTRY_ISSUED_AT, REGISTRY_EXPIRES_AT and REGISTRY_OUT are required' >&2; exit 2; }
	@pnpm exec soksak-validate registry-build plugins --id official --sequence '$(REGISTRY_SEQUENCE)' --issued-at '$(REGISTRY_ISSUED_AT)' --expires-at '$(REGISTRY_EXPIRES_AT)' --out '$(REGISTRY_OUT)'

authenticate: prepare
	@test -n "$$SOKSAK_REGISTRY_ED25519_SEED" && test -n '$(REGISTRY_IN)' && test -n '$(REGISTRY_OUT)' || { echo 'SOKSAK_REGISTRY_ED25519_SEED, REGISTRY_IN and REGISTRY_OUT are required' >&2; exit 2; }
	@pnpm exec soksak-validate registry-authenticate '$(REGISTRY_IN)' --out '$(REGISTRY_OUT)'
