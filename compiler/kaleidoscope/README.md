# LLVM Language Frontend

Going through the [LLVM First Language Frontend Tutorial](https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/index.html), learning the basics of creating a programming language by building the lexer, parser and AST nodes, and then hooking it up with the [LLVM Backend](https://llvm.org/) to achieve the basic functionality for the language.

Requires [LLVM](https://releases.llvm.org/) to be installed (I used version 17.0.6 to follow along). You can run the REPL using the `./run.sh` script.

## Key Learnings

Some things I learned through following along with the tutorial:

- Lexer reads the source code into tokens, exposed via a `gettok` function to read the next token from standard input
- The Abstract Syntax Tree (AST) serves as a way to create nodes that define the language and how different constructs are represented within it (derived from a base `ExprAST` class to express all nodes)
- Parser can be implemented using recursive descent to parse all kinds of expressions and return parsed statements as AST nodes
  - Not all methods created for parsing have (or should be) exposed as a public interface, but rather they are called as helpers within those user-exposed methods (like `ParseDefinition`, `ParseTopLevelExpr`, `ParseExtern`, etc.)
- With the frontend pieces set up, we can utilize the to actually get the hard parts done lol. It turns the frontend output into IR that can then be translated to produce target-specific machine code, handling code optimizations and everthing else along the way
  - It works within this tutorial by implementing `codegen` functions for each expression node, which turns into IR. We use the `llvm` C++ library for this
  - We can declare several optimization passes that will simplify the "instruction set" for a given architecture
  - LLVM simplifies ways that we can implement basic language features, such as control flow, user-defined operators, and mutable variables

In summary, this tutorial teaches how to reason about building the frontend of a language (lexer, parser, AST), and then using the LLVM IR to handle the rest (code optimizations, generating code for different architectures, and more).

NGL, it took me some time to go through it and I still don't understand everything, so I think I will continue learning through projects like [Daisy](https://github.com/daisylanguage/daisy) and putting my current knowledge into practice. It would be cool to learn how to build a compiler backend down the line too...
