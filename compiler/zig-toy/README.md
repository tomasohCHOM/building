# Toy Language in Zig

My attempt at creating a toy programming language purely in Zig.

## Language Specification

This toy language will support the following features:

- Top-level function definitions and `let` declarations
- `int` and `bool` as the only built-in types
- Expressions for literals, binary ops, parentheses, function calls, and variable references
- Statements with the following keywords: `let`, `return`, and `if`
- Function definitions (no closures, no nested functions, and no overloading)

Example:

```rust
let x: int = 10;

fn add(a: int, b: int) -> int {
    return a + b;
}

fn main() -> int {
    let y: int = add(x, 20);
    return y;
}
```
