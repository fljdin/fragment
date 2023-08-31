# Fragment

[![tests](https://github.com/fljdin/fragment/actions/workflows/tests.yml/badge.svg)](https://github.com/fljdin/fragment/actions/workflows/tests.yml)
[![go report](https://goreportcard.com/badge/github.com/fljdin/fragment)](https://goreportcard.com/report/github.com/fljdin/fragment)

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

* The `Language` struct is used to define the **delimiters** and **rules** for
  text splitting.
* A new fragment is built as soon as one of the **delimiters** is detected when
  reading the text.
* Each `Rule` consists of a start and stop condition that defines a context in
  which **delimiters** should be ignored.

### StringRule

The `StringRule` struct defines simple string-based rules to detect the start
and stop of fragments. Here's a concise example of newline separated fragments
with an exception rule when escape character preceed a newline (inpired from
shell syntax).

```go
// define a basic escape newline rule
escape := fragment.StringRule{
    Start: "\\",
    Stop:  "\n",
}

// define a new language to split lines from a text
text := fragment.Language{
    Delimiters: []string{"\n"},
    Rules:      []fragment.Rule{escape},
}
```

### RegexRule

The `RegexRule` struct allows use of regular expressions to define rules for
detecting the start and stop of fragments. Capture groups are supported in the
`Stop` regex to dynamically replace placeholders by positional values found
during `Start` regex's call.

The following example use a regular expression to find the start of a
[dollar-quoted string], in which we should ignore delimiters. To handle capture
group, the `RegexRule` must be passed by pointer.

[dollar-quoted string]: https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-DOLLAR-QUOTING

```go
// define a postgresql dollar-quoted expression rule
dollar := &fragment.RegexRule{
    Start: `\$([a-zA-Z0-9_]*)\$`,
    Stop:  `\$\1\$`,
}

// define a new language to split queries from a text
pgsql := fragment.Language{
    Delimiters: []string{";"},
    Rules:      []fragment.Rule{dollar}
}
```

To perform case-insensitive searches, use the `(?i)` flag within your regular
expression pattern. This flag indicates that the pattern matching should be done
without considering letter case. Here's an example of using the case-insensitive
flag to create a rule for XML markup tags:

```go
markup := &fragment.RegexRule{
    Start: `(?i)<(\w+)>`,
    Stop:  `(?i)</\1>`,
}

xml := fragment.Language{
    Delimiters: []string{`\n`},
    Rules:      []fragment.Rule{markup},
}
```

### Splitting Text

You can use the `Split` method of the `Language` struct to split text into
fragments based on the defined rules and delimiters. All leading and trailing
white space for each fragment are removed.

```go
// split the source text into fragments
fragments := text.Split(`
    Line 1
    Line 2 \
      on multiple lines
    Line 3
`)

for _, fragment := range fragments {
    fmt.Print("---- ")
    fmt.Println(fragment)
}
```

Will print:

```
---- Line 1
---- Line 2 \
      on multiple lines
---- Line 3
```

## Testing

Unit tests are provided under `tests` package.

```bash
go test ./tests
```

## Contributing

Contributions are welcome! Feel free to fork the repository, make changes, and
submit pull requests. If you find any bugs or have suggestions for improvements,
please create an issue.

## License

This project is licensed under the [MIT License](LICENSE).
