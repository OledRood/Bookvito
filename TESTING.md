# Bookvito Testing Model

## Layers

- `unit`: fast local tests without Docker.
- `integration`: backend endpoint tests with real router + Postgres + MinIO from `back/docker-compose.test.yml`.
- `e2e`: browser/API end-to-end tests via Playwright in `web/tests/e2e`.

## Commands

Backend (`/back`):

- `npm run test:unit`
- `npm run test:integration`
- `npm run test:all`

Frontend (`/web`):

- `npm run test:unit`
- `npm run test:e2e`
- `npm run test:all`

## Naming And Structure

- Backend unit: `Test<Subject>_<Condition>_<Expected>`.
- Backend integration: `back/tests/integration/*_test.go`.
- Frontend unit/component: `<name>.test.ts(x)` near runtime files.
- Frontend E2E: `<flow>.spec.ts` under `web/tests/e2e`.
- Shared fixtures/builders are centralized in:
  - backend: `back/tests/integration/helpers_test.go`, `back/cmd/testseed`
  - frontend: `web/tests/unit/*`, `web/tests/e2e/helpers.ts`

## Quality Gates

- Frontend unit coverage thresholds are enforced in `web/vitest.config.ts`:
  - lines/statements/functions: `>= 70%`
  - branches: `>= 60%`
- Backend unit coverage gate logic is implemented in `back/scripts/test-unit.sh`:
  - overall backend target: `>= 65%`
  - `internal/usecase` target: `>= 80%`
  - if current local Go toolchain has no `covdata`, script runs tests without coverage gate and prints a warning.

## Isolation And Cleanup

- Backend integration cleanup: `TRUNCATE ... CASCADE` + MinIO bucket object cleanup between tests.
- Playwright auth bootstrap: `web/tests/e2e/global-setup.ts` generates role storage states (`user/moder/admin`) from `back/cmd/testseed`.
- External API for E2E is isolated with local stub server:
  - `web/tests/e2e/google-books-stub.cjs`
  - wired in `web/playwright.config.ts` via `GOOGLE_BOOKS_BASE_URL`.
