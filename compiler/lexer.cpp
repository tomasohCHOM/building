#include "lexer.h"
#include "token.h"
#include <cctype>
#include <cstdio>
#include <cstdlib>

std::string IdentifierStr; // Filled in if tok_identifier
double NumVal;             // Filled in if tok_number

int gettok() {
  int LastChar = ' ';

  // Skip any whitespace.
  while (std::isspace(LastChar))
    LastChar = getchar();

  if (std::isalpha(LastChar)) { // identifier: [a-zA-Z][a-zA-Z0-9]*
    IdentifierStr = LastChar;
    while (std::isalnum((LastChar = std::getchar())))
      IdentifierStr += LastChar;
    if (IdentifierStr == "def")
      return tok_def;
    if (IdentifierStr == "extern")
      return tok_extern;
    return tok_identifier;
  }

  if (std::isdigit(LastChar) || LastChar == '.') { // Number: [0-9.]+
    std::string NumStr;
    do {
      NumStr += LastChar;
      LastChar = getchar();
    } while (std::isdigit(LastChar) || LastChar == '.');
    NumVal = std::strtod(NumStr.c_str(), 0);
    return tok_number;
  }
  if (LastChar == '#') {
    // Comment until end of line
    do
      LastChar = std::getchar();
    while (LastChar != EOF && LastChar != '\n' && LastChar != '\r');

    if (LastChar != EOF)
      return gettok();
  }

  if (LastChar == EOF)
    return tok_eof;

  int ThisChar = LastChar;
  LastChar = std::getchar();
  return ThisChar;
}
