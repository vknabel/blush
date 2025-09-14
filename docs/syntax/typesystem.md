# Typesystem

Blush is a dynamically, but strongly typed language. In the future it may add type inference, but for now all types can be assigned to any variable, field or function parameter.

In the end, we try to keept the type system simple and easy to understand. Blush supports the following classes of types:

- `data` types are the most common types. They are used to store data and can be easily created by calling the type name as a function.
- `enum` types are used to express that their values can be one of a group of types. In other languages they are also called union types.
- `annotation` types are used to annotate declarations with metadata. They can only be constructed at compile time.
- `extern` types are built-in types that are implemented in the runtime like `Func`, `String` or `Int`.

## Data types

Data types are the most common types in Blush and store data. In most other languages they are called classes or structs. They are defined by the `data` keyword followed by the type name and a list of fields.

```blush
data Person {
    name
    age
}
```

To create a new instance of a data type, simply call the type name as a function with the field values as arguments.

```blush
let person = Person("John", 42)
```

### Fields

Fields are the building blocks of data types. They are defined by their name and optionally annotations.
To increase the expressiveness, function fields can be defined by adding a function signature after the field name.

```blush
data Greetable {
    greeting(ofValue) // function field
}
```

### Field annotations

Fields can be annotated with metadata. These annotations will be processed at compile time and can be accessed at runtime.

```blush
data Person {
    @String
    @json.HasKey("name")
    name

    @Int
    @json.HasKey("age")
    age
}
```

## Enum types

Enum types are used to express that their values can be one of a group of types. In other languages they are also called union types. They are defined by the `enum` keyword followed by the type name and a list of types.
As convenience, you can even declare types within the enum declaration. These will still be available outside of the enum.

```blush
enum JuristicPerson {
    Person
    data Company {
        name
        corporateForm
    }
}
```

In this example every `Person` and every `Company` is a `JuristicPerson`.

The only way to check wether a given value is of an enum type, ist to tuse the `type`-expression.
It requires you to list all types of the enum type. It returns a function which takes a valid enum type.

```
import strings

func nameOf(juristic) {
    type juristic = JuristicPerson {
        Person: { person => person.name },
        Company: { company =>
                strings.concat [
                company.name, " ", company.corporateForm
            ]
        }
    }
}

nameOf you
```

> _**Attention:** If the given value is not valid, your program will crash. If you might have arbitrary values, you can add an `Any` case. As it matches all values, make sure it is always the last value._

## Annotation types

Annotation types are used to annotate declarations with metadata. They can only be constructed at compile time.
They are defined by the `annotation` keyword followed by the type name and a list of fields.

```blush
module json

annotation HasKey {
    name
}
```

To annotate a declaration, start with the `@` symbol followed by the annotation type name and a list of field values.

```blush
import json

data Person {
    @Type(String)
    @json.HasKey("name")
    name

    @Int // shorthand for @Type(Int)
    @json.HasKey("age")
    age
}
```

### Accessing annotations

Annotations can be accessed at runtime by using the `reflect`-module.

```blush
import json
import reflect

let person = Person("John", 42)

let nameAnnotation = reflect.typeOf(person).
    field("name").
    annotation(json.HasKey)

```

## Extern types

Extern types are built-in types that are implemented in the runtime like `Func`, `String` or `Int`. They are defined by the `extern` keyword followed by the type name. Optionally you can add a list of fields.

```blush
extern Int

extern String {
    length
}
```

Each extern type behaves slightly different in terms of how it is created and accessed.
Many types like `String`, `Int`, `Float` and `Dict` will be created by literals, types like `Func` and `Module` by declarations. `Any` on the other hand is more like an `enum` containing all types.

> _**Note:**_ The `extern` keyword is also used to declare functions provided by the compiler like `extern print(str)`.
