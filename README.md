# Fragment

Fragment is a Go package designed to split text into fragments. It provides a
convenient way to extract meaningful units of text from a larger body of text,
while also supporting special rules to ignore delimiters within specific
contexts.

## Installation

install the package using the following command:

```bash
go get github.com/fljdin/fragment
```

## Usage

Import package

```go
import "github.com/fljdin/fragment"
```

### Defining Language

Every text input follow predefined language's delimiter and rules.

```go
type Language struct {
	Delimiters []string
	Rules      []RangeRule
}
```

The `Language` struct is used to define the delimiters and rules for text
splitting. Each `RangeRule` consists of a start and end condition that defines a
context in which delimiters should be ignored.

```go
// define the source language
file := fragment.Language{
    Delimiters: []string{"\n"},
    Rules: []RangeRule{
        {Start: "\\", End: "\n"},
    }
}
```

### Splitting Text

You can use the `Split` method of the `Language` struct to split text into
fragments based on the defined rules and delimiters.

```go
// split the source text into fragments
fragments := file.Split(`
    Line 1
    Line 2 \
      on multiple lines
    Line 3
`)

// will print "Found 3 fragments"
fmt.Println("Found", len(fragments), "fragments")
```

## Contributing

Contributions are welcome! Feel free to fork the repository, make changes, and
submit pull requests. If you find any bugs or have suggestions for improvements,
please create an issue.

## License

This project is licensed under the [MIT] License - see the [LICENSE] file for
details.

[MIT]: https://choosealicense.com/licenses/mit
