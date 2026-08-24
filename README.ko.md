# Soksak 플러그인 Registry

공식 Registry는 데이터 저장소다. 각 `plugins/<id>.json`은 공통
`{id, version, url, size, sha256}` 형식으로 현재 불변 plugin release 하나를 가리킨다. Plugin
작성자는 자기 파일 하나만 변경한다. Sidecar는 owner release 내부의 runtime dependency이며 독립
상품으로 등록하지 않는다.

PR은 digest가 고정된 공개 `soksak-spec` package로 변경 entry와 전체 전이 release chain을
검증한다. main 병합 후에는 모든 entry를 결정적으로 합성하고 `registry.json`을 인증하여
owner-enforced immutable `registry-<sequence>` GitHub Release의 유일한 asset으로 게시한다. 생성된
Registry byte, 서명 키, parser, publisher 구현은 이 저장소에 두지 않는다. Core가 public trust
root를 내장한다.

[docs/REGISTRY.ko.md](docs/REGISTRY.ko.md)를 참고한다.

## 명령과 dependency 소유권

```sh
make verify
```

정확한 환경과 immutable spec validation package는 `.node-version`,
`package.json#packageManager`, `pnpm-lock.yaml`이 소유한다. `make build`는 sequence·시간·output을
명시적으로 받고, `make authenticate`는 입력·출력과 호출자가 주입한 signing seed를 받는다.
GitHub Actions는 tag 발견, secret, credential, publication만 소유한다. Registry는 release train의
마지막 consumer다. 변경된 spec이 immutable해진 뒤 HTTPS dependency를 갱신하고 clean install과
전체 release chain을 검증한 다음, plugin·Sidecar·Kit·spec source를 다시 빌드하지 않고 인증된
Registry byte를 공개한다.
