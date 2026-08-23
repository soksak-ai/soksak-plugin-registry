# Registry 기여와 게시

## Pull request

Plugin 작성자는 `plugins/<plugin-id>.json` 하나만 추가하거나 교체한다. 파일은 공통
`{id, version, url, size, sha256}` release reference다. 파일명, entry identity, 내려받은 plugin
release identity가 모두 일치해야 한다. Release와 전이 plugin·Sidecar dependency 전체를 URL, byte
크기, SHA-256, kind, ID, version, manifest, artifact matrix, conformance evidence로 검증한다. Source
checkout, branch, `latest`, package registry fallback, 저장소 topology 발견은 금지한다.

## 게시

main 게시 시각은 exact commit timestamp에서 파생한다. 같은 commit을 다시 실행하면 같은 sequence를
재사용하고 동일한 인증 byte를 다시 만든다. 새 commit은 가장 높은 immutable `registry-N` Release를
`N+1`로 올린다. 첫 plugins-only Registry는 기존 signed sequence 10을 이어 sequence 11이다.

Workflow는 정확한 `soksak-spec` tarball 하나를 내려받아 digest를 검증하고 Registry를 생성·인증한
뒤 같은 package의 immutable publisher를 사용한다. Signing seed는 GitHub secret에만 있고 public key는
다운로드 Registry 밖의 Core에 내장된다. 저장소에는 생성 Registry 파일이나 별도 signer·parser 구현을
두지 않는다.
