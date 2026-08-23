# soksak-plugin-registry

Soksak 공식 release catalogue입니다. 공개 형식과 서명 규칙은 `soksak-spec`이 소유하고
`soksak-contract-registry`가 검증합니다.

Catalogue는 plugin, sidecar, kit, contract, spec의 직접 release를 게시합니다. 각 owner
repository는 자신의 immutable archive와 conformance report를 build하고 검증합니다. Runtime
requirement는 owner manifest가 소유하며 사용자의 sidecar 역할 binding은 `environment.json`에
있습니다. 이 repository는 owner source tree를 읽거나 build하지 않습니다.

## 등록

`cmd/register`는 owner가 게시한 `release.json`을 읽어 같은 exact identity를 교체하고 종류별로
정렬합니다. 전체 catalogue를 검증한 뒤에만 `registry-source.json`을 atomic하게 갱신합니다.

```sh
go run ./cmd/register https://github.com/soksak-ai/example/releases/download/v0.0.1/release.json
go test ./...
go vet ./...
go run ./cmd/verify
```

`registry-signed.json`은 현재 운영 중인 signed index입니다. Core는 대응하는 public trust root를
고정하며 [SIGNING.ko.md](SIGNING.ko.md)는 만료 갱신과 sequence continuity를 정의합니다. Unsigned
catalogue byte는 설치 출처로 사용하지 않습니다.
