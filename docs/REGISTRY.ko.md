# Registry 기여와 게시

## Pull request

Plugin 작성자는 `plugins/<plugin-id>.json` 하나만 추가하거나 교체한다:

```json
{
  "id": "soksak-plugin-example",
  "version": "1.2.3"
}
```

Entry는 dependency 의도다. id와 version 외에는 없다. `url`, `size`, `sha256`을 담은 entry는
거부된다. 위치는 도출하고 digest는 게시 시점에 읽는다. 파일명과 entry id가 일치해야 하고, 도출한
위치 `https://github.com/soksak-ai/<id>/releases/download/v<version>/release.json`의 release 문서가
같은 id와 version을 선언해야 한다. Release와 전이 plugin·Sidecar dependency 전체를 각자 도출한
위치에서 byte 크기, SHA-256, kind, ID, version, manifest, artifact matrix, conformance evidence로
검증한다. Source checkout, branch, `latest`, package registry fallback, 저장소 topology 발견은
금지한다.

## 게시

main 게시 시각은 exact commit timestamp에서 파생한다. 같은 commit을 다시 실행하면 같은 sequence를
재사용하고 동일한 인증 byte를 다시 만든다. 새 commit은 가장 높은 immutable `registry-N` Release를
`N+1`로 올린다. 첫 plugins-only Registry는 기존 signed sequence 10을 이어 sequence 11이다.

`make build`는 모든 entry를 읽고, 도출한 위치에서 release 문서를 내려받고, 내려받은 byte로 인덱스
entry `{id, version, size, sha256}`를 쓴다. 인덱스는 release 문서의 다른 내용을 복사하지 않는다.
Core가 각 release 문서에서 `runtimeDependencies`를 걷는다. `make authenticate`가 인덱스를 서명한다.
Signing seed는 GitHub secret에만 있고 public key는 다운로드 Registry 밖의 Core에 내장된다.

`soksak-spec` package는 `package.json`이 정확한 version으로 선언하고 `pnpm-lock.yaml`이 integrity로
고정한다. Package registry는 make 인자 `REGISTRY`다. 저장소에는 생성 Registry 파일이나 별도
signer·parser·resolver 구현을 두지 않는다.
