#include "lexer.h"
#include "parser.h"
#include <iostream>

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

