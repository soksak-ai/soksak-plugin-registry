#!/usr/bin/env bash
# registry.json 갱신 — soksak-ai org 의 PUBLIC soksak-plugin-* repo 의 plugin.json 을 집계한다.
#
# 단일 진실은 각 플러그인의 독립 repo(plugin.json)다. 이 스크립트는 그것을 fetch 해 카탈로그
# (registry.json)로 모을 뿐, 어떤 값도 새로 만들지 않는다. 기본 브랜치를 repo 별로 자동 감지하므로
# main/master 무관하게 동작한다(브랜치 하드코딩 금지). private·template 은 설치 대상이 아니므로 제외한다
# (공개 레지스트리는 설치 가능한 공개 플러그인만 싣는다). 새 plugin 은 org 에 공개 repo 가 생기면
# 다음 실행에서 자동 편입된다.
#
# 배포 게이트(인터페이스 기반): 각 후보 플러그인을 soksak-plugin-doctor 로 검사한다. 무결성(테마 계약·
# 권한·명명) 미통과 플러그인은 카탈로그에서 제외 — 설치 불가 = 사실상 미배포. 플러그인은 코어 계약에
# conform 만 하면 되고, 검증은 배포 경계(여기)에서 외부로 수행한다(플러그인은 doctor 에 종속되지 않는다).
#
# 사용: ./update.sh  (gh 인증 + curl + jq + node + git 필요)
set -euo pipefail
cd "$(dirname "$0")"

# 배포 게이트 준비 — doctor 를 한 번 clone(코어 발행 contract.json 을 vendoring 한 게이트).
DOCTOR_ROOT=$(mktemp -d)
trap 'rm -rf "$DOCTOR_ROOT"' EXIT
git clone -q --depth 1 https://github.com/soksak-ai/soksak-plugin-doctor "$DOCTOR_ROOT/doctor" \
  || { echo "doctor clone 실패 — 게이트 없이 카탈로그 발행 금지" >&2; exit 1; }
DOCTOR="$DOCTOR_ROOT/doctor/bin/doctor.mjs"

# 한 번의 호출로 name + visibility + 기본 브랜치 조회 — 브랜치를 하드코딩하지 않는다.
repos=$(gh repo list soksak-ai --limit 200 --json name,visibility,defaultBranchRef \
  --jq '.[]
        | select(.name | startswith("soksak-plugin-"))
        | select(.name != "soksak-plugin-registry")
        | select(.name != "soksak-plugin-doctor")
        | select(.visibility == "PUBLIC")
        | "\(.name)\t\(.defaultBranchRef.name // "main")"')

out="[]"
while IFS=$'\t' read -r id branch; do
  [ -z "$id" ] && continue
  pj=$(curl -fsSL "https://raw.githubusercontent.com/soksak-ai/$id/$branch/plugin.json" 2>/dev/null) \
    || { echo "  skip $id ($branch/plugin.json fetch 실패)" >&2; continue; }
  echo "$pj" | jq -e . >/dev/null 2>&1 || { echo "  skip $id (invalid json)" >&2; continue; }
  [ "$(echo "$pj" | jq -r '.template // false')" = "true" ] && { echo "  skip $id (template)" >&2; continue; }

  # 배포 게이트 — doctor 무결성 검사. main.js(entry) 를 받아 임시 디렉토리(이름=id, 명명 검사용)에서 검사.
  entry_file=$(echo "$pj" | jq -r '.entry // "main.js"')
  pdir="$DOCTOR_ROOT/$id"; mkdir -p "$pdir"
  printf '%s' "$pj" > "$pdir/plugin.json"
  curl -fsSL "https://raw.githubusercontent.com/soksak-ai/$id/$branch/$entry_file" > "$pdir/$entry_file" 2>/dev/null \
    || { echo "  skip $id ($entry_file fetch 실패 — 빌드 산출물 없음)" >&2; continue; }
  if ! node "$DOCTOR" "$pdir" >"$pdir/.doctor.log" 2>&1; then
    echo "  skip $id (doctor 무결성 미통과 — 카탈로그 제외):" >&2
    sed 's/^/    /' "$pdir/.doctor.log" >&2
    continue
  fi

  # branch 명시 — plugin.json 이 선언하면 그것, 아니면 감지된 기본 브랜치(설치 시 default 가정 금지).
  entry=$(echo "$pj" | jq -c --arg br "$branch" '{id,name,version,description,author,repo, branch: (.branch // $br)} | with_entries(select(.value!=null))')
  out=$(echo "$out" | jq -c ". + [$entry]")
done <<< "$repos"

echo "$out" | jq 'sort_by(.id) | {spec: "soksak-registry@1", plugins: .}' > registry.json
echo "registry.json 갱신: $(jq '.plugins | length' registry.json)종 (PUBLIC, doctor 게이트 통과분만)" >&2
