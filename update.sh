#!/usr/bin/env bash
# registry.json 갱신 — soksak-ai org 의 soksak-plugin-* repo 의 plugin.json 을 집계한다.
#
# 단일 진실은 각 플러그인의 독립 repo(plugin.json)다. 이 스크립트는 그것을 fetch 해
# 카탈로그(registry.json)로 모을 뿐, 어떤 값도 새로 만들지 않는다. template 플러그인
# (개발 템플릿)은 설치 대상이 아니므로 제외한다. 새 플러그인은 org 에 repo 가 생기면
# 다음 실행에서 자동 편입된다.
#
# 사용: ./update.sh  (gh 인증 + curl + jq 필요)
set -euo pipefail
cd "$(dirname "$0")"

ids=$(gh repo list soksak-ai --limit 200 --json name -q '.[].name' \
  | grep '^soksak-plugin-' | grep -v '^soksak-plugin-registry$' | sort)

out="[]"
for id in $ids; do
  pj=$(curl -fsSL "https://raw.githubusercontent.com/soksak-ai/$id/main/plugin.json" 2>/dev/null) \
    || { echo "  skip $id (plugin.json fetch 실패)" >&2; continue; }
  echo "$pj" | jq -e . >/dev/null 2>&1 || { echo "  skip $id (invalid json)" >&2; continue; }
  [ "$(echo "$pj" | jq -r '.template // false')" = "true" ] && { echo "  skip $id (template)" >&2; continue; }
  entry=$(echo "$pj" | jq -c '{id,name,version,description,author,repo} | with_entries(select(.value!=null))')
  out=$(echo "$out" | jq -c ". + [$entry]")
done

echo "$out" | jq 'sort_by(.id) | {spec: "soksak-registry@0.0.1", plugins: .}' > registry.json
echo "registry.json 갱신: $(jq '.plugins | length' registry.json)종" >&2
