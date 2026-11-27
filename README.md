# explain-loop

> 설명할 때 쓰는 **패턴**을 관리하고, 어떤 패턴으로 설명했고 **얼마나 이해됐는지**를 쌓아서,
> 다음 설명을 더 잘 이해되게 만드는 도구.

```
요청 ──▶ 컨텍스트 수집 ──▶ 패턴 선택(+근거) ──▶ 설명 생성 ──▶ 이해도 점수
  ▲                                                              │
  └───────────────────── 다음 선택에 반영 ◀──────────────────────┘
```

## 문서 지도

| 계층 | 문서 | 내용 |
|------|------|------|
| 가치 | [explain-loop-values.md](docs/explain-loop-values.md) | 최상위 판단 기준 (V1–V5) |
| PRD-1 | [explain-loop-prd-pattern-registry.md](docs/explain-loop-prd-pattern-registry.md) | 설명 패턴의 동적 관리 |
| PRD-2 | [explain-loop-prd-context-capture.md](docs/explain-loop-prd-context-capture.md) | 요청·컨텍스트 수집과 저장 |
| PRD-3 | [explain-loop-prd-pattern-selection.md](docs/explain-loop-prd-pattern-selection.md) | 패턴 선택·근거 저장·설명 생성 |
| PRD-4 | [explain-loop-prd-understanding-feedback.md](docs/explain-loop-prd-understanding-feedback.md) | 이해도 점수화와 다음 설명 반영 |
| 테스트 | `docs/explain-loop-test-*.md` | PRD별 AC 검증 시나리오 (T1–T4) |
| 상태 | [explain-loop-doc-tracker/](docs/explain-loop-doc-tracker/index.md) | 가치–PRD–AC–테스트 연결 상태 (월별) |

## 아직 정하지 않은 것 (구현 전 결정 필요)

PRD는 구현 형태와 무관하게 썼다. 아래는 구현을 시작하기 전에 정해야 한다.

1. **제품 형태** — Claude 스킬 / MCP 서버 / CLI / 웹앱 중 무엇으로 시작할지.
   설명 생성은 LLM이 하고, 이 도구는 패턴·기록·점수를 관리하는 쪽이 자연스럽다
   (예: MCP 서버가 `select_patterns`, `record_explanation`, `submit_score` 도구를 노출).
2. **저장소** — git으로 관리하는 파일(JSON/YAML)인지, DB인지. 패턴 레지스트리는
   파일+git이 버전 관리에 유리하고, 요청·점수 기록은 양이 늘면 DB가 유리하다.
3. **사용자 범위** — 단일 사용자(나 자신)용인지, 여러 학습자를 구분하는지.
   PRD는 학습자 구분이 가능하도록 썼다.
4. **컨텍스트 수집 범위** — 대화 전체를 남길지 발췌만 남길지, 민감 정보는 어떻게 다룰지.
