# Registry 서명 상태

`registry-source.json`은 검증된 catalogue 원본입니다. 수동 서명 workflow는 원본을 검증하고
서명한 뒤 self-verification을 수행하며 `registry-signed.json`만 게시합니다. Production trust root,
Ed25519 private key secret, signed index는 아직 없으므로 Core는 이 catalogue를 작동하는 공식
registry로 신뢰할 수 없습니다.

서명 활성화는 하나의 명시적인 trust-root 작업입니다.

1. Core trust root 교체를 승인합니다.
2. 통제된 환경에서 Ed25519 key pair를 생성합니다.
3. private key는 GitHub Actions secret에만 저장합니다.
4. 새 public key와 key ID를 Core에 고정합니다.
5. 정확한 expiry를 지정해 수동 workflow를 실행합니다.
6. Installed acceptance에서 잘못된 서명, sequence rollback, key continuity를 검증합니다.

개발 fallback, unsigned 허용 경로, private key commit은 허용하지 않습니다.

## 만료 갱신

`issuedAt`과 `expiresAt`은 서명되는 byte입니다. 같은 sequence에 새 timestamp를 게시하면
equivocation이 되므로 signer는 게시된 index보다 증가하지 않은 sequence를 거부합니다. 만료 전에 다음을
수행합니다.

1. `go run ./cmd/renew`를 실행합니다.
2. Sequence만 변경된 `registry-source.json`을 검증하고 commit합니다.
3. 더 늦은 정확한 expiry를 지정해 수동 signing workflow를 실행합니다.

일반 catalogue 등록은 이미 sequence를 증가시키므로 별도 renew가 필요하지 않습니다.
