# Blush

Blush is an experimental programming language with a reference implementation in Go. The project is in an early stage and offers the essential building blocks of a modern language:

- Lexer, parser, and AST
- Bytecode compiler and virtual machine
- Standard library with prelude types such as `Array`, `Bool`, `Int`, and `String`
- Package management via `Cavefile` and registries
- Documentation on syntax, types, and style in [`docs/`](docs)

## Language overview

A quick glimpse at core features:

### Variables, functions, and data types

```blush
let answer = 42
let items = [1, 2, 3]

func greet(name) {
    "Hello, " + name
}

data Person {
    name
    age
}

let alice = Person("Alice", 30)
print(greet(alice.name))
```

### Control flow
Blush provides both expression and statement variants of `if` and `for`.
The examples below use the expression forms, which yield values and can be
nested inside other expressions.

```blush
let message = if answer == 42 { "yes" } else { "no" }

for item <- items {
    print(item)
}
```

When only side effects are required, use the statement forms, which run for
their effects and produce no value:

```blush
if answer == 42 {
    print("yes")
} else {
    print("no")
}

for item -> items {
    print(item)
}
```

An empty `for { }` forms an infinite loop, useful for servers or event
processors. It runs forever and only terminates when a `break` statement is
encountered.

### Enums

```blush
enum Result {
    data Ok { value }
    data Err { message }
}

let r = Ok(42)
```

### Annotations

Blush has no interfaces; instead, annotations attach behaviour and metadata to
types and functions. They enable generic code to rely on declared capabilities
without a formal interface system:

```blush
annotation Countable { length(value) }

@Countable({ v -> v.length })
data Bag {
    items
    length
}
```

Here `Countable` supplies a `length` implementation, allowing tools to treat
`Bag` like any other countable collection.

## Prerequisites

- [Go](https://go.dev/) 1.23 or newer

## Getting started

Clone the repository and run the test suite to see the implementation in action:

```bash
go test ./...
```

For more examples and guides, explore the [`docs/`](docs) and [`stdlib/`](stdlib) directories.

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).
