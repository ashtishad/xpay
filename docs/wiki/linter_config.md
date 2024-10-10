# Linter Configuration Guide

golangci-lint is a fast, parallel runner for Go linters. It's used in the XPay project to ensure code quality and consistency. This configuration is used both locally and in the GitHub Actions CI pipeline.

## Configuration File: .golangci.yaml

### Run Settings

```yaml
run:
  timeout: 3m
  tests: false
```

- `timeout`: Sets a 3-minute timeout for linter execution.
- `tests`: Disables linting of test files.

**Purpose**: These settings optimize the linting process, focusing on production code and ensuring the linter doesn't run indefinitely.

### Linter Selection

```yaml
linters:
  disable-all: true
  enable:
    # bugs/error
    - staticcheck
    - gosec
    - errcheck

    # performance
    - prealloc

    # style, formatting
    - gofmt
    - goconst
    - unconvert
    - misspell
    - unparam
    - nakedret
    - tagliatelle
    - dupl
```

This section explicitly enables specific linters while disabling all others. The enabled linters are categorized by their primary focus:

1. **Bugs/Error Detection**:
   - `staticcheck`: Comprehensive static analyzer
   - `gosec`: Security-focused linter
   - `errcheck`: Checks for unchecked errors

2. **Performance**:
   - `prealloc`: Suggests slice preallocation

3. **Style and Formatting**:
   - `gofmt`: Standard Go formatter
   - `goconst`: Finds repeated strings that could be constants
   - `unconvert`: Removes unnecessary type conversions
   - `misspell`: Finds commonly misspelled English words
   - `unparam`: Reports unused function parameters
   - `nakedret`: Finds naked returns in functions
   - `tagliatelle`: Checks struct tags
   - `dupl`: Finds code clones

**Purpose**: This curated selection of linters helps catch common errors, improve performance, and maintain consistent code style.

### Linter Settings

```yaml
linters-settings:
  gofmt:
    rewrite-rules:
      - pattern: 'interface{}'
        replacement: 'any'
      - pattern: 'a[b:len(a)]'
        replacement: 'a[b:]'

  misspell:
    locale: US

  errcheck:
    check-type-assertions: true

  dupl:
    threshold: 150
```

- `gofmt`: Configures rewrite rules for modern Go syntax.
- `misspell`: Sets the spell-checking locale to US English.
- `errcheck`: Enables checking of type assertions.
- `dupl`: Sets a threshold of 150 tokens for considering code as duplicated.

**Purpose**: These settings fine-tune the behavior of specific linters to match project conventions and catch more relevant issues.

### Issue Configuration

```yaml
issues:
  max-same-issues: 0
  max-issues-per-linter: 0
  exclude-use-default: false
  exclude:
    - G104
```

- `max-same-issues`: No limit on the number of same issues reported.
- `max-issues-per-linter`: No limit on issues reported per linter.
- `exclude-use-default`: Disables default exclusions.
- `exclude`: Specifically excludes the G104 (Unhandled errors) issue from gosec.

**Purpose**: This configuration ensures all issues are reported, giving developers a comprehensive view of potential problems in the codebase.

## Usage

### Local Usage

Run golangci-lint locally using:

```bash
golangci-lint run
```

Or use the Makefile command:

```bash
make lint
```

### GitHub Actions Integration

In the GitHub Actions workflow (`test.yaml`), golangci-lint is run as part of the CI process:

```yaml
- name: Golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: ${{ env.GOLANGCI_LINT_VERSION }}
```

This step uses the golangci-lint GitHub Action to run the linter with the project's configuration.

## Best Practices

1. Run linters locally before pushing changes to catch issues early.
2. Regularly update the golangci-lint version in both local development and CI environments.
3. Periodically review and update the linter configuration to align with evolving project needs and Go best practices.
4. Use inline comments (`//nolint:lintername`) sparingly to disable linters for specific lines when absolutely necessary.

## Troubleshooting

- If linting is too slow, consider using `golangci-lint run --fast` for quicker checks during development.
- For memory issues, try increasing the `GOGC` environment variable: `GOGC=200 golangci-lint run`.
- If certain linters are causing problems, you can temporarily disable them using the `--disable` flag.

## Further Reading

- [golangci-lint Documentation](https://golangci-lint.run/)
- [Effective Go - Code Style](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [CoachroachDB Team's Go Coding Guidelines](https://cockroachlabs.atlassian.net/wiki/spaces/CRDB/pages/181371303/Go+Golang+coding+guidelines)
