# Changelog

## 0.3.1

### Fixed
- Code folding missing or ending too early for nested `if` / `while` / `thread` blocks
- Code folding not shown when branch bodies contain only comments
- Document formatting and indentation for `if(` and `while(` without a space after the keyword

## 0.3.0

### Added
- Per-parameter documentation in signature help (from `@param` tags on builtins)
- `halt` as a keyword (highlighting, completion, and no false “not defined” diagnostics)

### Fixed
- Signature help race where stale buffer content could hide help after edits
- Signature help on multiline function calls
- Signature help documentation formatting (markdown aligned with hover)

## 0.2.0

### Added
- Default parameter values in function definitions
- Multiline string support
- Multiline function calls and list literals
- Diagnostic warning when `for` loops are used (URScript only has `while`)

### Fixed
- Crash on empty symbol names in document symbols
- Crash on nil pointer in inlay hints
- Docstring `### ... ###` blocks now correctly strip leading `#` markers
- Standardized I/O builtin documentation wording

### Changed
- Consolidated internal packages (24 → 6) for maintainability
- Removed `for` loop support (not part of URScript)
- Version is now derived from the git tag at build time

## 0.1.0

Initial release.

- Diagnostics (undefined functions, wrong argument counts, duplicate definitions)
- Completions (built-in functions, user symbols, keywords)
- Hover documentation
- Signature help
- Go to definition (cross-file)
- Find references (cross-file)
- Rename (cross-file)
- Formatting
- Document symbols
- Semantic highlighting
- Code actions (typo quick fixes)
- Inlay hints (parameter names)
- Code lens (reference counts)
- Folding ranges
- Call hierarchy
