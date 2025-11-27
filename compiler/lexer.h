#pragma once
#include <string>
#include "token.h"

class Lexer {
public:
  Lexer();

  // getNextToken reads another token from the lexer and returns its results.
  int getNextToken();
  const std::string &getIdentifierStr() const { return IdentifierStr; }
  double getNumVal() const { return NumVal; }

private:
  // gettok - Return the next token from standard input
  int gettok();

  int LastChar;
  std::string IdentifierStr; // Filled in if tok_identifier
  double NumVal; // Filled in if tok_number
};

