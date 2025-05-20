# Testing and Coverage Guide

This document explains how to run tests and check test coverage for the Sutram project.

## Running Tests

To run all tests:

```bash
make test
```

## Checking Test Coverage

To run tests with coverage analysis:

```bash
make test-coverage
```

This will:
1. Run all tests and generate a coverage profile
2. Display a summary of coverage per package and function

## GitHub Actions Test Workflow

Each pull request to the `main` branch automatically triggers a GitHub Actions workflow that:
1. Runs all tests
2. Identifies which Go files were changed in the PR
3. Checks test coverage specifically for those changed files
4. Ensures that new code has 100% test coverage

The workflow configuration is located in `.github/workflows/pr-test.yml`.

## Writing Tests

When adding new code, ensure it is accompanied by tests that provide 100% coverage. This helps maintain code quality and prevents regressions.

## Coverage Requirements

- All new or modified Go files must have 100% test coverage
- This requirement is enforced by the GitHub Actions workflow on every PR
- Existing files without tests are not required to have 100% coverage immediately, but any changes to those files will require test coverage