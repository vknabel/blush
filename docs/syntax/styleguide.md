# Styleguide

To make Lithia code more readable and consistent, we have a few rules and conventions that we follow.
In general, abbreviations should be avoided and names should be as descriptive as possible.

## Modules

In Lithia we try to have many micro modules instead of a few large ones. Try to keep the size of modules small with only a few, small files. Consider hiding implementation details in a separate `internal` module.

Module names are written in `snake_case`, but should avoid using underscores if possible. The module name should be short, expressive, easy to read and ideally in plural form. The module name should harmonize with the names of the declarations it contains.

```lithia
// good
module http
module json
module strings
moodule reflect // reflect.typeOf

// bad
module http_client
module json_parser
module string_utils
module reflection // reflection.typeOf
```

If the filename or its members might exist in a different module, declare the `module` at the top of your file.

```lithia
module my_module
```

## Data types

Data types are written in `PascalCase`. In most cases, the name should be a noun. If the type represents a collection, the name might be plural, otherwise it should be singular. If the type is a witness type, it should describe the ability of the type. Witnesses may be postfixed with `able` or `ible`.
Avoid prefixing the type name with the module name.

### Good

```lithia
data Person
data People
data Printable // witness type
data Error // in module http
```

### Bad

```lithia
data person
data people
data Printer // witness type
data HttpError // in module http
```

## Functions

Function names are written in `camelCase`. The name should be a verb or a verb phrase. Ideally it should produce a sentence when combined with the name of the module it is defined in.

### Good

```lithia
print
printLine
reflect.typeOf
```

### Bad

```lithia
print_line
reflect.reflectType
```

## Variables

Variable names are written in `camelCase`. The name should be a noun or a noun phrase. The name of a variable should be as long as its scope. If the variable is only used in a small scope, it may have a short name. If the variable is used globally, it should have a longer name, that still plays nicely with the module name.

When implementing specific patterns, the variable should reuse common names.

```lithia
// good
let name = "John"
let person = Person("John", 42)
let err = http.Error("Not found")

func printPersonName(p) {
    print(p.name)
}

// bad
let n = "John"
let p = Person("John", 42)
let error = http.Error("Not found")

func printName(n) {
    print(n.name)
}
```

## Annotations

Annotations are written in `PascalCase`. The name should describe a property of the declaration it annotates.

```lithia
// good
annotation Returns
annotation HasKey
annotation IsOptional

// bad
annotation Return
annotation Key
annotation Optional
```

## Enums

Enum names are written in `PascalCase`. The name should be a noun or a noun phrase. If the enum represents a collection, the name might be plural, otherwise it should be singular. If the enum represents witness types, it should describe the ability of the type, consider a `Witness` postfix.

```lithia
// good
enum Stateful
enum JuristicPerson
enum FunctorWitness

// bad
enum StateOrStore
enum JuristicPersons
enum Functor
```

If the cases within the enum are more relevant than the enum, declare them outside. Otherwise inside.

```lithia
// good
enum Optional {
    data Some { value }
    data None
}

enum Maybe {
    Optional
    Any
}

// bad
enum Maybe {
    enum Optional {
        data Some { value }
        data None
    }
    Any
}
```
