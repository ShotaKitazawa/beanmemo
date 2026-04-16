# CLAUDE.md

## Architecture

### API Design
`openapi.yaml` is the single source of truth. All backend (ogen) and frontend (openapi-typescript) code is generated from it. Any API change must start with modifying `openapi.yaml`.

### Backend
Responsibilities are separated into three layers: Handler → UseCase → Repository. Dependencies flow in one direction only; upper layers depend on interfaces of lower layers.

### Frontend
Page-level components and custom hooks are kept separate. All API communication is encapsulated in custom hooks — components must not call fetch directly.

## Development Cycle

All commands (tests, code generation, linting, builds, etc.) must be run via `mise run <task>`. Available tasks are defined in `mise.toml`. Never invoke `go test`, `pnpm test`, or other tool-specific commands directly — always use the corresponding `mise run` task.

### Exception: commands that modify dependency manifests

The following commands modify `package.json` and must **never** be run as mise tasks or invoked automatically. Run them manually with explicit human intent only:

- `pnpm run sync:pin` (from `frontend/`) — detects packages violating `minimumReleaseAge` and automatically pins them to an older compliant version via `pnpm.overrides` in `package.json`.

- Changes that require code generation (OpenAPI, SQL) must pass `mise run pre-merge`, which detects missing re-generation via `git diff --exit-code`.
- Generated code is committed to the repository; CI assumes generated files are up to date.
- `mise run pre-merge` must pass before merging any PR.

## Test-Driven Development

All feature development and bug fixes must follow the RED → GREEN → REFACTOR cycle:

1. **RED** — Write a failing test that specifies the desired behavior. Confirm the test fails before writing any implementation code.
2. **GREEN** — Write the minimum implementation code needed to make the test pass. Do not over-engineer at this stage.
3. **REFACTOR** — Clean up the implementation and tests while keeping all tests green. Eliminate duplication and improve clarity without changing behavior.

### Backend (Go)

- Place test files alongside source files (`*_test.go`).
- Use table-driven tests for UseCase and Repository layers.
- Mock lower-layer interfaces when testing upper layers (e.g., mock Repository when testing UseCase).
- Run tests with `mise run test-backend`.

### Frontend (TypeScript / React)

- Place test files alongside source files (`*.test.ts` / `*.test.tsx`).
- Test custom hooks with `@testing-library/react` (`renderHook`).
- Test components with `@testing-library/react` focusing on user-visible behavior, not implementation details.
- Run tests with `mise run ci-test-frontend`.
