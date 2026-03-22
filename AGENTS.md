# AGENTS.md

## Project Definition: Versifyr
Versifyr is a CLI tool designed to manage project versions across multiple files of different types (Go, YAML, XML, Java, etc.) using templates and configuration. It ensures consistency by replacing version strings based on a central source of truth or command-line arguments.

## Mandatory Agent Workflow

### 1. Bug Fixes
- **Reproduction First**: Before implementing any fix, you MUST create a test case or a reproduction script that systematically demonstrates the bug.
- **Verification**: A fix is only considered complete when the previously failing test case passes and no regressions are introduced.

### 2. Version Management
- **Tool Usage**: Never manually edit version strings in the codebase. You MUST use `versifyr` itself to advance the version.
- **Timing**: The code's internal version must be advanced BEFORE committing and tagging a new release.

### 3. Commit Standards
- **Conventional Commits**: Use the Conventional Commits specification (e.g., `feat:`, `fix:`, `docs:`, `chore:`).
- **Contextual Detail**: Commits should be contextual, explaining the "why" and referencing specific issues or architectural reasons when applicable.

### 4. Release Process
- **Automation**: Be aware of CI/CD workflows (e.g., `.github/workflows/go2.yml`). Tagging a release triggers automated builds and artifact generation.
- **Consistency**: Ensure the internal version matches the Git tag exactly.
