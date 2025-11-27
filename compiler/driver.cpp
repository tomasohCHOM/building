#include "lexer.h"
#include "parser.h"
#include <iostream>

Lexer lexer;
Parser parser(lexer);

static void HandleDefinition() {
  if (parser.ParseDefinition()) {
    fprintf(stderr, "Parsed a function definition.\n");
  } else {
    // Skip token for error recovery.
    lexer.getNextToken();
  }
}

static void HandleExtern() {
  if (parser.ParseExtern()) {
    fprintf(stderr, "Parsed an extern\n");
  } else {
    // Skip token for error recovery.
    lexer.getNextToken();
  }
}

static void HandleTopLevelExpression() {
  // Evaluate a top-level expression into an anonymous function.
  if (parser.ParseTopLevelExpr()) {
    fprintf(stderr, "Parsed a top-level expr\n");
  } else {
    // Skip token for error recovery.
    lexer.getNextToken();
  }
}

/// top ::= definition | external | expression | ';'
static void MainLoop() {
  while (true) {
    fprintf(stderr, "ready> ");
    switch (parser.CurTok) {
    case tok_eof:
      return;
    case ';': // ignore top-level semicolons.
      lexer.getNextToken();
      break;
    case tok_def:
      HandleDefinition();
      break;
    case tok_extern:
      HandleExtern();
      break;
    default:
      HandleTopLevelExpression();
      break;
    }
  }
}

int main() {
  Lexer lexer;
  Parser parser(lexer);

  std::cout << "ready> ";
  while (true) {
    int tok = lexer.getNextToken();
    if (tok == tok_eof) break;
    std::cout << "Token: " << tok << std::endl;
  }
}

