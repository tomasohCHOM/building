#ifndef LEXER_H
#define LEXER_H

#include <string>

/// gettok - Return the next token from standard input
int gettok();

extern std::string IdentifierStr; // Filled in if tok_identifier
extern double NumVal;             // Filled in if tok_number

#endif // LEXER_H
