#ifndef PARSER_H
#define PARSER_H

#include "ast.h"
#include <map>
#include <memory>

extern int CurTok;
extern std::map<char, int> BinopPrecedence;

int getNextToken();

std::unique_ptr<FunctionAST> ParseDefinition();
std::unique_ptr<FunctionAST> ParseTopLevelExpr();
std::unique_ptr<PrototypeAST> ParseExtern();

#endif // PARSER_H
