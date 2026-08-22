# Registry 서명 상태

`registry-source.json`은 검증된 catalogue 원본입니다. 현재 이 저장소에는 서명 workflow,
Ed25519 private key secret, 게시된 `registry-signed.json`이 없습니다. 따라서 Core는 이
catalogue를 작동하는 공식 registry로 신뢰할 수 없습니다.

서명 활성화는 하나의 명시적인 trust-root 작업입니다.

1. Core trust root 교체를 승인합니다.
2. 통제된 환경에서 Ed25519 key pair를 생성합니다.
3. private key는 GitHub Actions secret에만 저장합니다.
4. 원본 검증, 서명, signed document 게시 workflow를 추가합니다.
5. 새 public key와 key ID를 Core에 고정합니다.
6. 잘못된 서명 거부, sequence rollback 거부, key rotation 복구를 테스트합니다.

개발 fallback, unsigned 허용 경로, private key commit은 허용하지 않습니다.
