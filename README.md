# AI 종목 분석 리포트 웹서비스

## 🧭 개요

AI Agent가 분석한 종목별 주식 리포트를 활용하여, 사용자가 웹을 통해 쉽게 리포트를 열람할 수 있는 서비스.  
리포트는 날짜/종목명/의사결정/언어별로 Markdown 파일로 저장되며, 이를 자동으로 수집, 정리하여 웹에서 보기 쉽게 제공.

## 📁 시스템 구조

```
├── main.go          # Go 서버 메인 파일
├── config/          # 서버 설정
│   ├── config.go    # CORS 등 설정 관리
│   └── config.json  # 설정 파일
├── public/          # 정적 파일
│   ├── index.html   # 메인 페이지
│   ├── css/         # 스타일시트
│   ├── js/          # 자바스크립트
│   └── assets/      # 이미지 등 리소스
└── reports
    └── 2025-06-14
        ├── BUY
        │   └── TSLA
        │       ├── EN
        │       │   ├── complete.md
        │       │   ├── final_trade_decision.md
        │       │   ├── investment_plan.md
        │       │   ├── market_report.md
        │       │   └── trader_investment_plan.md
        │       └── KO
        │           ├── complete.md
        │           ├── final_trade_decision.md
        │           ├── investment_plan.md
        │           ├── market_report.md
        │           └── trader_investment_plan.md
        ├── HOLD
        │   └── NVDA
        │       └── EN
        │           ├── complete.md
        │           ├── final_trade_decision.md
        │           ├── investment_plan.md
        │           ├── market_report.md
        │           └── trader_investment_plan.md
        └── SELL
            └── NVDA
                ├── EN
                │   ├── complete.md
                │   ├── final_trade_decision.md
                │   ├── investment_plan.md
                │   ├── market_report.md
                │   └── trader_investment_plan.md
                └── KO
                    ├── complete.md
                    ├── final_trade_decision.md
                    ├── investment_plan.md
                    ├── market_report.md
                    └── trader_investment_plan.md
```

## ✅ 주요 기능 정의

### 1. 리포트 자동 탐색 및 불러오기 (백엔드)
- 폴더 구조: `reports/YYYY-MM-DD/DECISION/종목명/LOCALE/타입.md`
- 파일 내용 읽어와서 구조화된 JSON 배열로 반환
- API: `GET /api/reports`
- 반환 형태:
```json
[
  {
    "date": "2025-06-14",
    "name": "NVDA",
    "decision": "HOLD",
    "locale": "EN",
    "type": "complete",
    "content": "AI 분석 결과 내용..."
  }
]
```

### 2. 리포트 목록 페이지 (프론트엔드)
- 순수 HTML/CSS/JavaScript로 구현
- API 호출 후 데이터 받아 렌더링
- 리스트/그리드 뷰 전환 가능
- 언어 선택 (EN/KO)
- 리포트 타입 선택
- 카드 형태로 날짜, 종목명, 분석 요약 출력
- 기본 리스트는 최신 날짜부터 정렬
- 마크다운 렌더링 (marked.js 사용)

## 📡 API 명세

### 1. 리포트 조회 API

#### 1.1 모든 리포트 조회
```
GET /api/reports
```

**쿼리 파라미터:**
- `locale` (선택): 언어 선택 (EN/KO)
  - 기본값: EN
  - 예시: `?locale=KO`
- `type` (선택): 리포트 타입
  - 기본값: complete
  - 가능한 값: complete, final_decision, investment_plan, market_report, trader_plan
  - 예시: `?type=investment_plan`

**응답 예시:**
```json
[
  {
    "date": "2025-06-14",
    "name": "TSLA",
    "decision": "BUY",
    "locale": "EN",
    "type": "complete",
    "content": "마크다운 형식의 리포트 내용..."
  }
]
```

#### 1.2 날짜별 리포트 조회
```
GET /api/reports/date
```

**쿼리 파라미터:**
- `date` (필수): 조회할 날짜 (YYYY-MM-DD 형식)
  - 예시: `?date=2025-06-14`
- `locale` (선택): 언어 선택 (EN/KO)
  - 기본값: EN
  - 예시: `?locale=KO`
- `type` (선택): 리포트 타입
  - 기본값: complete
  - 예시: `?type=investment_plan`

#### 1.3 투자 의사결정별 리포트 조회
```
GET /api/reports/decision
```

**쿼리 파라미터:**
- `decision` (필수): 투자 의사결정 (BUY/HOLD/SELL)
  - 예시: `?decision=BUY`
- `locale` (선택): 언어 선택 (EN/KO)
  - 기본값: EN
  - 예시: `?locale=KO`
- `type` (선택): 리포트 타입
  - 기본값: complete
  - 예시: `?type=investment_plan`

### 2. 응답 필드 설명

| 필드      | 타입     | 설명                    |
|---------|--------|-----------------------|
| date    | string | 리포트 작성일 (YYYY-MM-DD) |
| name    | string | 종목명                  |
| decision| string | 투자 의사결정 (BUY/HOLD/SELL) |
| locale  | string | 언어 (EN/KO)           |
| type    | string | 리포트 타입 (complete/final_decision/investment_plan/market_report/trader_plan) |
| content | string | 마크다운 형식의 리포트 내용     |

### 3. 에러 응답

모든 API는 에러 발생 시 다음과 같은 형식으로 응답합니다:

```json
{
  "error": "에러 메시지"
}
```

**HTTP 상태 코드:**
- 200: 성공
- 400: 잘못된 요청 (필수 파라미터 누락 등)
- 500: 서버 내부 오류

## 📌 요구사항 정의서

### 기능 요구사항

| 구분 | 항목         | 설명                              |
| -- | ---------- | ------------------------------- |
| 필수 | 리포트 목록 API | `/api/reports` API를 통해 모든 분석 리포트 제공 |
| 필수 | 날짜/종목 추출   | 파일 경로 및 이름에서 날짜와 종목명 자동 추출      |
| 필수 | 리포트 보기 UI  | 프론트에서 모든 리포트 리스트를 보여줌           |
| 필수 | 언어 선택      | EN/KO 언어 선택 기능 제공              |
| 필수 | 리포트 타입 선택  | 다양한 리포트 타입 선택 기능 제공           |
| 선택 | 마크다운 렌더링   | 텍스트를 HTML로 변환      |
| 선택 | 날짜/종목명 필터링 | 검색 또는 드롭다운 필터 (확장 가능)           |

### 비기능 요구사항

| 항목   | 설명                             |
| ---- | ------------------------------ |
| 확장성  | 리포트 파일이 많아질 경우에도 빠르게 탐색 가능해야 함 |
| 단순성  | MVP는 로그인이나 DB 없이 동작            |
| 보안   | 리포트 파일 외 외부 파일 접근 제한           |
| 유지보수 | 리포트 폴더 구조가 변경되면 쉽게 반영 가능해야 함   |

## ⚙️ 상세 구현 사항

### 📦 서버 (`Go`)

* 정적 파일 서빙
  * `public` 디렉토리의 파일들을 웹 루트로 서빙
  * API 엔드포인트는 `/api` 경로로 구분

* `/api/reports` 핸들러
  * `filepath.Walk` 사용해 `./reports` 하위 모든 `.md` 파일 탐색
  * 경로 예: `reports/2025-06-14/BUY/TSLA/EN/complete.md`
  * 파일명 → 종목명, 상위 폴더명 → 날짜/의사결정/언어
  * `Report` 구조체로 변환 후 JSON 반환

* 보안
  * `.md` 외 파일 필터링
  * 디렉토리 이동 제한 (`./reports`만 접근)
  * CORS 설정으로 허용된 도메인만 접근 가능

### 🖥️ 프론트엔드 (`HTML/CSS/JavaScript`)

* 순수 JavaScript로 API 호출
* 받은 데이터 리스트로 mapping
* CSS로 UI 구성
* marked.js로 마크다운 렌더링
* 언어 선택 기능
* 리포트 타입 선택 기능
* 리스트/그리드 뷰 전환 기능

## ⏱️ 향후 확장 계획 (Roadmap)

| 단계  | 기능             |
| --- | -------------- |
| 1단계 | 날짜/종목 필터링 UI   |
| 2단계 | 개별 리포트 페이지 라우팅 |
| 3단계 | 마크다운 HTML 렌더링  |
| 4단계 | 사용자 인증 및 구독 기능 |
| 5단계 | 백테스트, 추천 종목 기능 |

## 👤 타겟 사용자

* 개인 투자자 (초보\~중급)
* 투자 콘텐츠 생산자 (리포트를 활용할 유튜버, 블로거)
* 리서치에 기반한 투자 전략 수립을 원하는 사용자

## 🧪 테스트 항목

| 테스트 항목     | 방법            | 기대 결과         |
| ---------- | ------------- | ------------- |
| 리포트 API 응답 | `/api/reports` 호출 | JSON 배열 반환    |
| 파일 누락 처리   | 존재하지 않는 경로    | 오류 없이 무시      |
| UI 출력      | 다양한 종목 이름 표시  | 카드 형태로 정확히 출력 |
| 리포트 수 증가   | 100건 이상 리포트   | UI 느려짐 없이 출력  |
| 언어 전환      | EN/KO 전환     | 해당 언어 리포트만 표시  |
| 타입 필터링     | 타입별 선택      | 선택한 타입 리포트만 표시  |

## 📎 참고

* 마크다운 렌더링: [marked.js](https://marked.js.org/)
* CSS 프레임워크: [Tailwind CSS](https://tailwindcss.com/)
* Go `filepath.Walk` 성능 문제 발생 시: goroutine + 채널 구조로 개선 가능

