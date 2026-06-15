# soksak-plugin-registry

soksak 공식 플러그인 레지스트리 — 앱의 "설치 가능" 목록의 단일 진실.

`registry.json` 은 각 플러그인의 표시용 메타 + git 레포 URL(설치 source)을 담는다. soksak 앱은
빌드 스냅샷으로 즉시 표시하고, 세션 1회/새로고침 시 이 파일을 fetch 해 최신화한다.

## 스키마

```json
{
  "spec": "soksak-registry@1",
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

`name`/`description` 은 문자열 또는 `{ "ko": …, "en": … }` 다국어 객체. `repo` 는 임의 git URL
(github/gitlab/self-host). 실제 설치는 repo clone 후 앱이 manifest 를 엄격 재검증한다.

## 등재

새 플러그인 등재는 `registry.json` 에 엔트리를 추가하는 PR 로. 각 플러그인은 자기 git 레포에서
독립적으로 개발·버전관리된다(여러 저자).
