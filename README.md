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

Every text input follows predefined language's delimiter and rules.

```go
type Language struct {
	Delimiters []string
	Rules      []Rule
}
```

The `Language` struct is used to define the delimiters and rules for text
splitting. Each `Rule` consists of a start and end condition that defines a
context in which delimiters should be ignored.

### StringRule

The `StringRule` struct defines simple string-based rules to detect the start
and end of fragments.

```go
// define a basic escape newline rule
escape := fragment.StringRule{
    Start: "\\",
    End:   "\n",
}

// define a new language to split lines from a text
file := fragment.Language{
    Delimiters: []string{"\n"},
    Rules:      []fragment.Rule{escape},
}
```

### RegexRule

The `RegexRule` struct allows use of regular expressions to define rules for
detecting the start and end of fragments. You can also use capture groups in the
`End` regex to dynamically replace placeholders in the matched text.

To handle capture group, the `RegexRule` must be passed by pointer:

```go
// define a postgresql dollar-quoted expression rule
dollar := &fragment.RegexRule{
    Start: `\$([a-zA-Z0-9_]*)\$`,
    End:   `\$\1\$`,
}

// define a new language to split queries from a text
pgsql := fragment.Language{
    Delimiters: []string{";"},
    Rules:      []fragment.Rule{dollar}
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
