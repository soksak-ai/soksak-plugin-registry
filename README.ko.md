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

현재 source catalogue는 아직 `registry-signed.json`으로 서명·게시되지 않았습니다. 필요한 trust
root 절차는 [SIGNING.ko.md](SIGNING.ko.md)에 있습니다. 이 절차가 완료되기 전에는 official signed
registry 설치 경로를 운영 상태로 취급하지 않습니다. Unsigned fallback은 허용하지 않습니다.
