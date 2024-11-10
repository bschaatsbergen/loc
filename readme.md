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