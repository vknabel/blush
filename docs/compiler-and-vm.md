# Compiler and VM

- All constants will be stored inside a constant pool
- big endian



## OpCodes

| Mnemonic      | Width | Description                                    | Comments |
| ------------- | ----- | ---------------------------------------------- | -------- |
| const         | 2     | Push constant from constant pool               |          |
| consttrue     | 0     | Push boolean `true`                            |          |
| constfalse    | 0     | Push boolean `false`                           |          |
| pop           | 0     | Discard top of stack                           |          |
| array         | 0     | Build array from preceding values             | length on stack |
| dict          | 0     | Build dictionary from preceding key/value pairs | length on stack |
| asserttype    | 2     | Assert top value has given type ID             |          |
| jump          | 2     | Unconditional jump to address                  |          |
| jumptrue      | 2     | Jump if top value is truthy                    |          |
| jumpfalse     | 2     | Jump if top value is `false`                   |          |
| negate        | 0     | Numeric negation                               |          |
| invert        | 0     | Boolean NOT                                    |          |
| add           | 0     | Add two numbers                                |          |
| sub           | 0     | Subtract two numbers                           |          |
| mul           | 0     | Multiply two numbers                           |          |
| div           | 0     | Divide two numbers                             |          |
| mod           | 0     | Remainder of integer division                  |          |
| eq            | 0     | Compare for equality                           |          |
| neq           | 0     | Compare for inequality                         |          |
| gt            | 0     | Compare greater-than                           |          |
| gte           | 0     | Compare greater-than-or-equal                  |          |
| lt            | 0     | Compare less-than                              |          |
| lte           | 0     | Compare less-than-or-equal                     |          |
| debug         | 0     | Optional breakpoint instruction                | omitted in release builds |
