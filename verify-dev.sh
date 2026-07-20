#!/usr/bin/env bash
# 로컬 개발중 플러그인이 발행 게이트를 통과하는지 발행 전에 미리 검증한다 — update.sh 가 GitHub 공개
# repo 를 카탈로그할 때 강제하는 바로 그 두 게이트를 로컬 dev checkout 에 돌린다(멱등, 재실행 가능):
#   게이트①  매니페스트 스키마    (@soksak-ai/plugin-spec 의 soksak-validate — 단일진실)
#   게이트②  contract 무결성       (soksak-plugin-doctor — 코어 발행 contract.json 대조)
#
# 게이트 정본은 update.sh 와 같다. 다만 로컬 dev 환경엔 org 패키지가 퍼블릭 npm 에 없으므로(npx 404)
# npx 대신 로컬 소스로 해소한다: 스키마는 코어 repo 의 plugin-spec bin, 무결성은 로컬 doctor checkout.
# 발행 전 이 스크립트로 통과를 확인해야 update.sh 단계에서 등재-거부로 걸리지 않는다.
#
# 사용: ./verify-dev.sh <plugin-dir>...
#   CORE=<coredir>    코어 repo (기본: 이 스크립트 옆의 ../core)
#   DOCTOR=<doctordir> soksak-plugin-doctor checkout (기본: ~/.soksak-dev/plugins/soksak-plugin-doctor)
set -uo pipefail
cd "$(dirname "$0")"

CORE="${CORE:-$(cd .. 2>/dev/null && pwd)/core}"
DOCTOR="${DOCTOR:-$HOME/.soksak-dev/plugins/soksak-plugin-doctor}"
VALIDATE="$CORE/packages/plugin-spec/bin/validate.mjs"
DOCTOR_BIN="$DOCTOR/bin/doctor.mjs"

[ -f "$VALIDATE" ] || { echo "게이트① validate 없음: $VALIDATE  (CORE=$CORE 로 코어 repo 지정)" >&2; exit 1; }
[ -f "$DOCTOR_BIN" ] || { echo "게이트② doctor 없음: $DOCTOR_BIN  (DOCTOR=<checkout> 로 지정)" >&2; exit 1; }
[ "$#" -gt 0 ] || { echo "사용: $0 <plugin-dir>...  (검증할 로컬 플러그인 폴더들)" >&2; exit 2; }

pass=0; fail=0
for d in "$@"; do
  d="${d%/}"
  id="$(basename "$d")"
  [ -f "$d/plugin.json" ] || { echo "  ✗ $id (plugin.json 없음: $d)"; fail=$((fail+1)); continue; }
  errs=""
  node "$VALIDATE" plugin "$d/plugin.json" >/dev/null 2>&1 || errs+="①schema "
  node "$DOCTOR_BIN" "$d" >/dev/null 2>&1 || errs+="②doctor "
  if [ -z "$errs" ]; then
    pass=$((pass+1)); echo "  ✓ $id"
  else
    fail=$((fail+1)); echo "  ✗ $id — 실패 게이트: $errs"
  fi
done

echo "verify-dev: PASS=$pass FAIL=$fail"
[ "$fail" -eq 0 ]
