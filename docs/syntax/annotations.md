# Annotations

Lithia supports annotating a declaration with metadata. These annotations will be processed at compile time and can be accessed at runtime.

## Syntax

Annotations are written as a list of `@` followed by the annotation name and a list of arguments. The arguments are separated by commas and can be either an identifier, a string or a number. More complex arguments need to be referenced by their name.

```lithia
@AnnotationName("argument", 123)
data MyData
```

Before being able to use an annotation, it must be declared and imported:

```lithia
module my

annotation AnnotationName {
    field
    field2
}

// other module
import my

@my.AnnotationName
```

### Syntactic sugar

If the annotation has no arguments, the parentheses can be omitted.

```lithia
@AnnotationName
data MyData
```

If the referenced type is not an annotation, it will implicitly be converted to an `@Type` annotation.

```lithia
data MyData {
    @String // equivalent to @Type(String)
    field
}
```

## Built-in annotations

There are some built-in annotations that can be used to modify the behavior of the compiler or the runtime.

### `@Type`

The `@Type` annotation can be used to specify the type of a declaration. This is useful when the type cannot be inferred by the compiler. In case a declaration has been annotated with a non-annotation type like an `extern`, `enum` or `data`, the `@Type` annotation is required.

```lithia
data MyData {
    @Type(String)
    field
}

// or

data MyData {
    @Type("String")
    field
}
```
