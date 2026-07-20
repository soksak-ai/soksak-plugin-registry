#!/usr/bin/env bash
# registry.json 갱신 — soksak-ai org 의 PUBLIC soksak-plugin-* repo 의 plugin.json 을 집계한다.
#
# 단일 진실은 각 플러그인의 독립 repo(plugin.json)다. 이 스크립트는 그것을 fetch 해 카탈로그
# (registry.json)로 모을 뿐, 어떤 값도 새로 만들지 않는다. 기본 브랜치를 repo 별로 자동 감지하므로
# main/master 무관하게 동작한다(브랜치 하드코딩 금지). private·template 은 설치 대상이 아니므로 제외한다
# (공개 레지스트리는 설치 가능한 공개 플러그인만 싣는다). 새 plugin 은 org 에 공개 repo 가 생기면
# 다음 실행에서 자동 편입된다.
#
# 배포 게이트(2층, 외부 강제): 카탈로그에 들어오려면 두 게이트를 모두 통과해야 한다 — 우회 불가
# (저자 로컬 hook 과 달리 배포 경계에서 강제). 플러그인은 코어 계약에 conform 만 하면 된다:
#  1) 매니페스트 스키마(코어 packages/plugin-spec — 앱이 강제하는 단일진실): plugin.json 필드·타입·id 패턴·권한.
#  2) contract 무결성(soksak-plugin-doctor): 테마 계약·권한·명명·유령변수(코어 발행 contract.json 대조).
#
# 사용: ./update.sh  (gh 인증 + curl + jq + node + git + 코어 체크아웃[형제 ../core 또는 CORE=<경로>] 필요)
set -euo pipefail
cd "$(dirname "$0")"

# 게이트 1 준비 — 매니페스트 스키마 검증기는 코어 packages/plugin-spec(앱이 강제하는 그 단일진실)이다.
# npm 발행본에 기대지 않는다 — org 패키지는 퍼블릭 npm 에 없고, 동명의 낡은 패키지가 있으면 그 스키마가
# 현재 계약모델(consumes 등)을 몰라 유효 매니페스트를 전부 튕겨낸다(카탈로그 전멸). 형제 ../core 기본.
CORE="${CORE:-$(cd .. && pwd)/core}"
VALIDATE="$CORE/packages/plugin-spec/bin/validate.mjs"
[ -f "$VALIDATE" ] || { echo "게이트① 검증기 없음: $VALIDATE  (CORE=<코어 repo 경로> 로 지정)" >&2; exit 1; }

# 게이트 2 준비 — doctor 를 한 번 clone(코어 발행 contract.json 을 vendoring 한 게이트).
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
dep_edges=""  # 카탈로그된 플러그인의 의존 엣지(id<TAB>depId) — 발행 후 그래프 무결성 검증용
while IFS=$'\t' read -r id branch; do
  [ -z "$id" ] && continue
  pj=$(curl -fsSL "https://raw.githubusercontent.com/soksak-ai/$id/$branch/plugin.json" 2>/dev/null) \
    || { echo "  skip $id ($branch/plugin.json fetch 실패)" >&2; continue; }
  echo "$pj" | jq -e . >/dev/null 2>&1 || { echo "  skip $id (invalid json)" >&2; continue; }
  [ "$(echo "$pj" | jq -r '.template // false')" = "true" ] && { echo "  skip $id (template)" >&2; continue; }

  # 게이트 1 — 매니페스트 스키마(단일진실 = 코어 packages/plugin-spec). 검증기는 dirName(=폴더명)으로
  # id 일치도 보므로 $id 폴더에 임시 기록 후 검증한다. 스키마 미통과는 등재 거부.
  vtmp=$(mktemp -d); mkdir -p "$vtmp/$id"; printf '%s' "$pj" > "$vtmp/$id/plugin.json"
  if ! node "$VALIDATE" plugin "$vtmp/$id/plugin.json" >&2; then
    echo "  skip $id (매니페스트 스키마 검증 실패 — 등재 거부)" >&2; rm -rf "$vtmp"; continue
  fi
  rm -rf "$vtmp"

  # 게이트 2 — contract 무결성(doctor). 일반 플러그인은 entry(main.js) 산출물을 받아 검사한다.
  # 서비스 플러그인(entry:null + service 선언 — 웹뷰 진입 없이 로직을 사이드카 서비스가 소유,
  # PLUGIN-SERVICE.md)은 진입 파일 자체가 없다 — entry fetch 를 건너뛰고 매니페스트만 검사한다
  # (entry 부재는 빌드 누락이 아니라 부류의 정의다 — 모든 플러그인에 main.js 가정 금지).
  is_service=$(echo "$pj" | jq -r 'if (.entry == null) and (.service != null) then "yes" else "no" end')
  pdir="$DOCTOR_ROOT/$id"; mkdir -p "$pdir"
  printf '%s' "$pj" > "$pdir/plugin.json"
  if [ "$is_service" != "yes" ]; then
    entry_file=$(echo "$pj" | jq -r '.entry // "main.js"')
    curl -fsSL "https://raw.githubusercontent.com/soksak-ai/$id/$branch/$entry_file" > "$pdir/$entry_file" 2>/dev/null \
      || { echo "  skip $id ($entry_file fetch 실패 — 빌드 산출물 없음)" >&2; continue; }
  fi
  if ! node "$DOCTOR" "$pdir" >"$pdir/.doctor.log" 2>&1; then
    echo "  skip $id (doctor 무결성 미통과 — 카탈로그 제외):" >&2
    sed 's/^/    /' "$pdir/.doctor.log" >&2
    continue
  fi

  # branch 명시 — plugin.json 이 선언하면 그것, 아니면 감지된 기본 브랜치(설치 시 default 가정 금지).
  # commands = 매니페스트 선언 명령(이름+다국어 제목+위험 분류 — 설치 전 능력 조회용).
  # 값을 새로 만들지 않는다(투영만) — title 이 곧 사람용 설명(매니페스트 단일진실).
  entry=$(echo "$pj" | jq -c --arg br "$branch" '{id,name,version,description,author,repo, branch: (.branch // $br), commands: ([.contributes.commands[]? | {name, title} + (if .danger != null then {danger} else {} end)] | sort_by(.name))} | with_entries(select(.value!=null))')
  out=$(echo "$out" | jq -c ". + [$entry]")
  # 이 플러그인의 의존 대상(플러그인↔플러그인)을 엣지로 기록 — 발행 후 그래프 검증.
  while IFS= read -r depId; do
    [ -n "$depId" ] && dep_edges+="${id}"$'\t'"${depId}"$'\n'
  done < <(echo "$pj" | jq -r '.dependencies // {} | keys[]')
done <<< "$repos"

# 최종 카탈로그(아직 미기록) — 의존 그래프 검증을 통과해야만 registry.json 에 기록한다.
final=$(echo "$out" | jq -c 'sort_by(.id) | {spec: "soksak-registry@0.0.1", plugins: .}')

# 의존 그래프 무결성 — 카탈로그된 플러그인의 의존 대상이 카탈로그에 함께 있는가. 없으면 카탈로그가
# 동작 못 하는 플러그인을 광고하는 상태다. per-plugin 스킵(고립 결함)과 달리 이건 크로스-플러그인
# 무결성 위반이라 카탈로그 전체를 큰 소리로 실패시킨다(무음 캐스케이드=은폐 금지 — 대상을 함께
# 발행하거나 참조를 정정하게 강제). 버전(semver) 세부는 코어 dependency-graph-scan 이 본다.
catalog_ids=$(echo "$final" | jq -r '.plugins[].id')
graph_violations=""
while IFS=$'\t' read -r who dep; do
  [ -z "$who" ] && continue
  grep -qxF "$dep" <<< "$catalog_ids" || graph_violations+="  ✗ ${who} → ${dep} (배포 카탈로그에 없음)"$'\n'
done < <(printf '%b' "$dep_edges")
if [ -n "$graph_violations" ]; then
  echo "의존 그래프 위반 — 카탈로그된 플러그인이 미배포 대상을 의존한다(registry.json 미기록):" >&2
  printf '%s' "$graph_violations" >&2
  echo "대상을 함께 발행하거나(예: 라이브러리 플러그인) 의존 참조를 실제 id 로 정정하라." >&2
  exit 1
fi

echo "$final" | jq '.' > registry.json
echo "registry.json 갱신: $(jq '.plugins | length' registry.json)종 (PUBLIC, 스키마+doctor+의존그래프 게이트 통과분만)" >&2
