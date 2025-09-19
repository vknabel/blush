# Copilot Instructions for Blush

## Project Overview

Blush is an experimental programming language implemented in Go with a bytecode compiler and virtual machine. The repository contains:

- **Language Type**: Programming language implementation (lexer, parser, compiler, VM)
- **Primary Language**: Go 1.23+ 
- **Size**: Medium-sized project (~18 test files, 17 Go packages)
- **Target**: Educational/experimental programming language with modern features

### Core Components

- **Lexer & Parser**: Tokenize and parse Blush source code into AST
- **Compiler**: Compile AST to bytecode for virtual machine
- **Virtual Machine**: Execute bytecode with stack-based architecture  
- **Runtime**: Built-in types (Array, Bool, Int, String) and standard library
- **Package Management**: Cavefile-based dependency management with registries

## Build & Validation Instructions

### Prerequisites

- Go 1.23.0 or newer (current environment has Go 1.24.7)
- All dependencies are managed via go.mod - no additional tools required

### Essential Commands

**ALWAYS run these commands from the repository root `/home/runner/work/blush/blush`:**

#### Build (Required before testing)
```bash
go build -v ./...
```
- **Time**: ~30-60 seconds on first run, ~5-10 seconds on subsequent runs
- **Purpose**: Compile all packages to verify syntax and dependencies
- **Always succeed**: Should exit with code 0, no output on success

#### Test (Critical validation)
```bash
go test ./...
```
- **Time**: ~2-5 seconds  
- **Purpose**: Run all unit tests across 18 test files
- **Expected**: All tests pass, ~10-15 packages tested
- **Note**: Some packages show "[no test files]" - this is normal

#### Test with Race Detection (Recommended for concurrency changes)
```bash
go test -race ./...
```
- **Time**: ~10-15 seconds
- **Purpose**: Detect race conditions in concurrent code
- **Use when**: Making changes to VM, runtime, or concurrent operations

#### Format Check (Always required)
```bash
gofmt -l .
```
- **Known Issues**: Files with formatting problems will be listed
- **Fix with**: `gofmt -w <filename>` for each file  
- **Critical**: Must fix before committing - GitHub Actions will fail otherwise

#### Format Fix (When gofmt -l shows files)
```bash
gofmt -w .
```
- **Purpose**: Fix all indentation issues (spaces vs tabs) across the codebase
- **Always run**: After making changes, before committing
- **Note**: Some files historically had formatting issues, this fixes all of them

#### Vet (Recommended)
```bash
go vet ./...
```
- **Time**: ~2-3 seconds
- **Purpose**: Check for suspicious constructs
- **Should pass**: Clean exit with no output

#### Module Maintenance (If dependency issues)
```bash
go mod tidy
```
- **Use when**: Adding/removing dependencies or getting module errors
- **Note**: May download additional test dependencies

### Command Sequences That Work

#### Clean Development Build
```bash
go clean -cache && go build -v ./... && go test ./...
```

#### Full Validation Sequence  
```bash
gofmt -l . && go vet ./... && go build -v ./... && go test ./...
```

#### Fix Formatting + Test
```bash
gofmt -w . && go test ./...
```

### Continuous Integration

GitHub Actions workflow (`.github/workflows/go.yml`):
1. **Go Setup**: Uses Go 1.23
2. **Build**: `go build -v ./...`  
3. **Test**: `go test -v ./...`

**Critical**: Both build and test must pass for PR approval. The workflow runs on every push and PR to main branch.

## Project Layout & Architecture

### Repository Structure

```
/
├── .github/workflows/go.yml    # CI/CD pipeline
├── README.md                   # Project documentation  
├── LICENSE                     # Mozilla Public License 2.0
├── go.mod                      # Go module definition
├── grammar.ebnf                # Formal grammar specification
├── docs/                       # Language documentation
│   ├── compiler-and-vm.md      # Bytecode & VM architecture
│   └── syntax/                 # Language syntax docs
├── examples/project/           # Example Blush project
├── stdlib/prelude/             # Standard library (.blush files)
├── ast/                        # Abstract Syntax Tree definitions
├── lexer/                      # Tokenization
├── parser/                     # Parse tokens to AST  
├── compiler/                   # Compile AST to bytecode
├── op/                         # Bytecode operation definitions
├── vm/                         # Virtual machine execution
├── runtime/                    # Built-in types and runtime system
├── cavefile/                   # Package management structures
├── registry/                   # Module/package registries
├── syncheck/                   # Syntax validation
├── token/                      # Token definitions
├── version/                    # Semantic versioning
├── world/                      # OS interaction abstractions
└── pkgmanager/                 # Package management logic
```

### Key Files for Code Changes

- **`ast/`**: Add new AST nodes when extending language syntax
- **`lexer/lexer.go`**: Modify for new token types
- **`parser/parser.go`**: Update for new language constructs  
- **`compiler/compiler.go`**: Add compilation logic for new features
- **`vm/vm.go`**: Extend VM for new bytecode operations
- **`runtime/prelude-*.go`**: Built-in type implementations
- **`op/defs.go`**: Define new bytecode operations

### Testing Patterns

- **Unit tests**: `*_test.go` files alongside source
- **Test data**: Inline test cases in table-driven tests
- **Compiler tests**: `compiler/compiler_test.go` - bytecode validation
- **Parser tests**: `parser/*_test.go` - AST validation  
- **Integration**: Through VM execution tests

### Dependencies (Not Obvious)

- **go-git**: Used for Git-based package registries
- **go-billy**: Virtual filesystem for package management
- **google/go-cmp**: Test comparison utilities
- **Native Go**: No external build tools, pure Go implementation

### Standards & Style

- **Code Format**: Use `gofmt` - tabs for indentation, Go standard formatting
- **Testing**: Table-driven tests, helper functions for common setup
- **Naming**: Go conventions - exported/unexported based on capital letters
- **Documentation**: Godoc comments for public APIs

### Validation Checklist for Changes

1. **Format**: `gofmt -l .` should return no files
2. **Build**: `go build -v ./...` must succeed  
3. **Test**: `go test ./...` must pass all tests
4. **Vet**: `go vet ./...` should report no issues
5. **Race**: `go test -race ./...` for concurrency-related changes

### Common Gotchas & Workarounds

- **Formatting**: Several files have tab/space mix - always use `gofmt -w` to fix
- **Module cache**: If weird dependency errors, try `go clean -cache && go mod tidy`
- **Test timing**: `gitreg` tests can be slow (~800ms) due to Git operations
- **No main package**: This is a library, not an executable - no `main.go`

## Instructions for Coding Agents

**Trust these instructions** - they are comprehensive and tested. Only search/explore if:
1. These instructions contradict current code structure  
2. New build failures occur that aren't covered here
3. Instructions appear incomplete for your specific task

**Always run the validation checklist before finalizing changes.** The GitHub Actions CI will reject PRs that fail formatting, building, or testing.