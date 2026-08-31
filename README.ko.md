# Soksak 플러그인 Registry

공식 Registry는 데이터 저장소다. 각 `plugins/<id>.json`은 현재 불변 plugin release 하나를 의도
`{id, version}`으로 가리킨다. Plugin 작성자는 자기 파일 하나만 변경한다. Sidecar는 owner release
내부의 runtime dependency이며 독립 상품으로 등록하지 않는다.

Release 위치는 도출하며 기록하지 않는다. `{id, version}`의 release 문서는
`https://github.com/soksak-ai/<id>/releases/download/v<version>/release.json`이다. PR은 변경 entry를
검증하고, 도출한 위치에서 release 문서를 읽고, 정확한 `soksak-spec` package로 전체 전이 release
chain을 검증한다. main push는 전체 검증을 반복할 뿐 게시하지 않는다. 검증된 commit을 대상으로
publication workflow를 명시적으로 dispatch해야 모든 entry를 결정적으로 서명 인덱스
`registry.json`으로 합성한다. 인덱스의 plugin reference `{id, version, size, sha256}`는 내려받은
release 문서의 크기와 SHA-256만 담는다. 인덱스는 `runtimeDependencies`를 복사하지 않는다. Core가
각 release 문서에서 closure를 걷기 때문이다. 인덱스를 인증하여 owner-enforced immutable
`registry-<sequence>` GitHub Release의 유일한 asset으로 게시한다. 생성된 Registry byte, 서명 키,
parser, publisher 구현은 이 저장소에 두지 않는다. Core가 public trust root를 내장한다.

[docs/REGISTRY.ko.md](docs/REGISTRY.ko.md)를 참고한다.

## 명령과 dependency 소유권

이 패키지는 `@soksak-ai/plugin-spec`에 의존하므로, install을 수행하는 모든 `make` 호출은 make
명령줄의 `REGISTRY`를 요구한다. 패키지가 `https://registry.npmjs.org`에 게시된 뒤에도 같다. 환경
변수로 전달된 값은 거부된다. Makefile은 `package.json`에서 이 요구를 읽고, 없으면
`REGISTRY required: this package depends on @soksak-ai/plugin-spec`으로 거부한다.

빌드 입력의 정체성은 `REGISTRY`가 아니라 `pnpm-lock.yaml`의 integrity다. pnpm은 content-addressable
store에 없는 integrity의 패키지만 `REGISTRY`에서 받으므로, 같은 기계에서 같은 lockfile을 다시
install하면 store를 읽고 `REGISTRY`에 접속하지 않는다.

```sh
make verify REGISTRY=http://host:port/
```

정확한 환경과 immutable spec validation package는 `.node-version`,
`package.json#packageManager`, `pnpm-lock.yaml`이 소유한다. `make build`는 sequence·시간·output을
명시적으로 받고, `make authenticate`는 입력·출력과 호출자가 주입한 signing seed를 받는다.
GitHub Actions는 tag 발견, secret, credential, publication만 소유한다. Publication workflow는
수동 호출 전용이다. Push workflow는 그 credential을 받을 수 없고 Registry tag·Release·asset을
변경할 수 없다. Registry는 release train의 마지막 consumer다. 변경된 spec이 immutable해진 뒤
정확한 version을 갱신하고 clean install과 전체 release chain을 검증한 다음,
plugin·Sidecar·Kit·spec source를 다시 빌드하지 않고 인증된 Registry byte를 공개한다.
