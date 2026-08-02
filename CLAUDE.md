# Project Rules

## Git

- **Never `git push` without explicit instruction.** Wait for the user to say "push" before pushing anything.
- When pushing, always write the PR title and description in English.
- Commits must be granular — one logical change per commit, so it is clear what each commit does.

## Branch Strategy — Read This Before Cutting Any Branch

- **One parent branch per feature/domain unit** (e.g. `feat/pins`, `feat/note`, `feat/plan`), cut from `main`.
- **Before cutting any sub-branch for a domain, confirm that domain's own parent branch already exists.** If it doesn't exist yet, create it from `main` first — do not skip this step.
- **Every sub-branch (layer branch: types/hooks/mapper/container/UI, or Domain/Infra/Application/Presentation) for a domain must be cut from that same domain's parent branch — never from another domain's branch.** e.g. Note sub-branches must be cut from `feat/note`, never from `feat/pins`, even if `feat/pins` happens to be the current branch or has related work in progress.
- If a sub-branch was cut from the wrong parent, stop and fix the branch structure (rebase / cherry-pick onto the correct parent) before adding more commits on top of it.
- When in doubt about which parent branch a new piece of work belongs under, ask before cutting the branch — don't guess.

## Backend Branch Strategy

- Cut feature branches from `main` by domain unit (e.g. `feat/user`, `feat/purchase`).
- From each feature branch, cut sub-branches for each DDD layer as needed:
  - `Domain`, `Infrastructure`, `Application`, `Presentation`
- Tests are required for each layer **except UseCase**.
- Always be mindful of Codecov coverage.
- Presentation layer: write **integration tests only** — no unit tests.
- API URL prefix: `api/v1`

## Frontend Branch Strategy

- Cut feature branches from `main` by domain unit (e.g. `feat/pins`, `feat/plans`).
- From the feature branch, cut sub-branches in this order:
  1. `types`
  2. `hooks`
  3. `mapper`
  4. `container` — one container per screen/view
  5. `UI` — one branch per screen/view
- All types/hooks/mapper work must be complete before moving to container.
- Container and UI branches are scoped per screen, not per feature as a whole.

## UI Development

- Do not implement UI components without prior discussion with the user.
- Always agree on the design before writing code.
