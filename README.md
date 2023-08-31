# Fragment

[![tests](https://github.com/fljdin/fragment/actions/workflows/tests.yml/badge.svg)](https://github.com/fljdin/fragment/actions/workflows/tests.yml)
[![go report](https://goreportcard.com/badge/github.com/fljdin/fragment)](https://goreportcard.com/report/github.com/fljdin/fragment)

Fragment is a Go package designed to split text into fragments. It provides a
convenient way to extract meaningful units of text from a larger body of text,
while also supporting special rules to ignore delimiters or other rules within
specific contexts.

## Installation

Install package using the following command:

```bash
go get github.com/fljdin/fragment
```

## Predefined languages

The `fragment/languages` package offers several predefined languages to help you
split text into fragments based on language-specific rules.

* The `PgSQL` defines delimiters for SQL statements and special PostgreSQL
  commands (like `\g` or `\gexec`). It also specifies various rules to handle
  comments, single-quoted and double-quoted string literals, transactions
  BEGIN-END blocks, and dollar-quoted strings.

```go
package main

import (
    "fmt"
    "github.com/fljdin/fragment/languages"
)

func main() {
    // Define a PostgreSQL query containing multiple 
    // SQL statements
    queries := `
        SELECT * FROM employees;
        -- This is a comment
        INSERT INTO products VALUES (1, 'Product 1');
        BEGIN;
        UPDATE orders SET status = 'Shipped' 
         WHERE order_id = 100;
        COMMIT;
        $custom_tag$
        This is a custom tag.
        $custom_tag$
    `

    // Split an SQL file script into fragments using 
    // the PgSQL language configuration
    fragments := languages.PgSQL.Split(queries)

    // Print the extracted fragments
    for i, fragment := range fragments {
        fmt.Printf("Fragment %d:\n%s\n\n", i+1, fragment)
    }
}
```

* `Shell` defines a newline delimiter, allowing for the accurate splitting of Shell
  code into meaningful fragments. The configuration also includes rules for
  handling single-quoted and double-quoted string literals, line continuation
  with backslashes, Shell-style comments, and Here-Documents (Heredocs).

* `XML` defines a newline delimiter to split documents. The primary rule in this
  configuration identifies XML tags, both opening and closing tags, in a
  case-insensitive manner.

These predefined languages make it easy to split text into fragments based on
language-specific rules, which can be useful in various text processing
applications.

## Usage

```go
import "github.com/fljdin/fragment"
```

Every text input follows predefined language's delimiters and rules.

* The `Language` struct defines the **delimiters** and **rules** required by
  text splitting.
* A new fragment is built as soon as a `Delimiter` is detected when reading the
  text.
* Each `Rule` consists of a start and stop condition that defines a context in
  which **delimiters** and other **rules** should be ignored.

### Delimiter

The `Delimiter` struct determines whether the fragment must be built when it
encounters a simple string or when it matches a regular expression. This
distinction is mainly made for performance reasons.

```go
// define a newline delimiter
newline := fragment.Delimiter{
    String: "\n",
}

// define a psql's meta-command delimiter
command := fragment.Delimiter{
    Regex: `\\g.*\n`,
}
```

### StringRule

The `StringRule` struct defines simple string-based rules to detect start and
stop of fragments. Here's a concise example of newline separated fragments with
an exception rule when an escape character preceeds a newline (inpired from
shell syntax).

```go
// define a basic escape newline rule
escape := fragment.StringRule{
    Start: "\\",
    Stop:  "\n",
}

// define a new language to split lines from a text
text := fragment.Language{
    Delimiters: []fragment.Delimiter{newline},
    Rules:      []fragment.Rule{escape},
}
```

In some cases, delimiters may not be ignored by a rule. The following example
shows how to define a comment that ends with a newline delimiter. The field
`Stop` is replaced by `StopAtDelim`.

```go
// define a one-line comment rule, inspired from shell 
comment := fragment.StringRule{
    Start:       "#",
    StopAtDelim: true,
}

// define a new language to split commands from a script
shell := fragment.Language{
    Delimiters: []fragment.Delimiter{newline},
    Rules:      []fragment.Rule{comment}
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
    Delimiters: []fragment.Delimiter{{String: ";"}},
    Rules:      []fragment.Rule{dollar}
}
```

To perform case-insensitive searches, use the `(?i)` flag within your regular
expression pattern. This flag indicates that the pattern matching should be done
without considering letter case. Here's an example of using the case-insensitive
flag to create a rule for XML markup tags:

```go
// define a markup rule with capture group and placeholder
markup := &fragment.RegexRule{
    Start: `(?i)<(\w+)>`,
    Stop:  `(?i)</\1>`,
}

// define a new language to split XML documents from a file
xml := fragment.Language{
    Delimiters: []fragment.Delimiter{newline},
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

Unit tests are provided under `languages` package.

```bash
go test ./languages
```

## Contributing

Contributions are welcome! Feel free to fork the repository, make changes, and
submit pull requests. If you find any bugs or have suggestions for improvements,
please create an issue.

## License

This project is licensed under the [MIT License](LICENSE).
