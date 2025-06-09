## LOC
Simple program that counts lines of code in a directory or provided repository.

### Usage

#### Build
```bash
go build -o loc main.go
```

#### Count lines of code in a provided directory
```bash
./loc -dir /path/to/directory
```

#### Count lines of code in a provided repository
```bash
./loc -repo github.com/username/repo
```

#### Exclude files and directories using regex patterns
```bash
# Exclude test files
./loc -dir /path/to/directory --exclude ".*_test\.go$" --exclude ".*\.spec\.ts$"

# Exclude multiple directories
./loc -dir /path/to/directory --exclude "node_modules/" --exclude "\.git/" --exclude "build/"

# Exclude generated files and test directories
./loc -dir /path/to/directory --exclude "gen/.*" --exclude "/tests/" --exclude ".*\.generated\."

# Multiple exclusions can be combined
./loc -dir /path/to/directory \
--exclude ".*_test\.go$" \
--exclude "node_modules/" \
--exclude "\.git/" \
--exclude "build/" \
--exclude "dist/"
```

### Supported Languages
- Go
- Python
- C
- Java
- Ruby
- Rust
- C#
- JavaScript
- TypeScript
- PHP
- HTML
- CSS
- Shell
- Kotlin
- Swift
- Scala
- Perl
- R
- Lua
- Haskell
- Objective-C
- Groovy
- Dart
- Elixir
- Erlang
- Fortran
- Pascal
- MATLAB
- Julia
- SQL
- JSON
- XML
- T-SQL
- VHDL
- COBOL
- Assembly
- ActionScript
- VimL
- Bash
- Ada
- Delphi
- Smalltalk
- Scheme
- Clojure
- F#
- OCaml
- Nim
- Racket
- C++