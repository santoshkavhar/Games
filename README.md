# Sudoku-Game

### To run
```
go run ./(name of the program)
```

### To build
```
go build
./(name of the folder)
```


### To Test
```
go test -v
```

# Test Coverage Commands
```
go test -v -coverprofile cover.out 
go tool cover -html=cover.out -o cover.html
open cover.html
```

gridGeneration.go 
Contains Logic for Grid Generation

validation.go
Contains only validation logic

Thanks Mayur Deshmukh for Frontend
