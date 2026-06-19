#!/usr/bin/env bash
# registry.json 갱신 — soksak-ai org 의 PUBLIC soksak-plugin-* repo 의 plugin.json 을 집계한다.
#
# 단일 진실은 각 플러그인의 독립 repo(plugin.json)다. 이 스크립트는 그것을 fetch 해 카탈로그
# (registry.json)로 모을 뿐, 어떤 값도 새로 만들지 않는다. 기본 브랜치를 repo 별로 자동 감지하므로
# main/master 무관하게 동작한다(브랜치 하드코딩 금지). private·template 은 설치 대상이 아니므로 제외한다
# (공개 레지스트리는 설치 가능한 공개 플러그인만 싣는다). 새 plugin 은 org 에 공개 repo 가 생기면
# 다음 실행에서 자동 편입된다.
#
# 사용: ./update.sh  (gh 인증 + curl + jq 필요)
set -euo pipefail
cd "$(dirname "$0")"

# 한 번의 호출로 name + visibility + 기본 브랜치 조회 — 브랜치를 하드코딩하지 않는다.
repos=$(gh repo list soksak-ai --limit 200 --json name,visibility,defaultBranchRef \
  --jq '.[]
        | select(.name | startswith("soksak-plugin-"))
        | select(.name != "soksak-plugin-registry")
        | select(.visibility == "PUBLIC")
        | "\(.name)\t\(.defaultBranchRef.name // "main")"')

out="[]"
while IFS=$'\t' read -r id branch; do
  [ -z "$id" ] && continue
  pj=$(curl -fsSL "https://raw.githubusercontent.com/soksak-ai/$id/$branch/plugin.json" 2>/dev/null) \
    || { echo "  skip $id ($branch/plugin.json fetch 실패)" >&2; continue; }
  echo "$pj" | jq -e . >/dev/null 2>&1 || { echo "  skip $id (invalid json)" >&2; continue; }
  [ "$(echo "$pj" | jq -r '.template // false')" = "true" ] && { echo "  skip $id (template)" >&2; continue; }
  # branch 명시 — plugin.json 이 선언하면 그것, 아니면 감지된 기본 브랜치(설치 시 default 가정 금지).
  entry=$(echo "$pj" | jq -c --arg br "$branch" '{id,name,version,description,author,repo, branch: (.branch // $br)} | with_entries(select(.value!=null))')
  out=$(echo "$out" | jq -c ". + [$entry]")
done <<< "$repos"

echo "$out" | jq 'sort_by(.id) | {spec: "soksak-registry@0.0.1", plugins: .}' > registry.json
echo "registry.json 갱신: $(jq '.plugins | length' registry.json)종 (PUBLIC, 기본 브랜치 자동 감지)" >&2
