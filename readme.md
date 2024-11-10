## LOC
Simple program that counts lines of code in a directory or provided repository.

### Usage

#### Build
```
go build -o loc main.go
```

#### Count lines of code in a provided directory
```
./loc -dir /path/to/directory

#### Count lines of code in a provided repository
```
./loc -repo github.com/username/repo
```