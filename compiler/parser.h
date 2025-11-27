#pragma once
#include "lexer.h"
#include "ast.h"
#include <memory>

class Parser {
public:
  explicit Parser(Lexer &lexer);
  std::unique_ptr<ExprAST> ParseExpression();

private:
  // CurTok is the current token the parser is looking at.
  int CurTok;
  Lexer &Lex;

  int getNextToken();
  int GetTokPrecedence();

  std::unique_ptr<ExprAST> ParsePrimary();

  std::unique_ptr<ExprAST> ParseNumberExpr();
  std::unique_ptr<ExprAST> ParseParenExpr();
  std::unique_ptr<ExprAST> ParseIdentifierExpr();
  std::unique_ptr<ExprAST> ParseBinOpRHS(int ExprPrec, std::unique_ptr<ExprAST> LHS);
  std::unique_ptr<PrototypeAST> ParsePrototype();
  std::unique_ptr<FunctionAST> ParseDefinition();
  std::unique_ptr<PrototypeAST> ParseExtern();
  std::unique_ptr<FunctionAST> ParseTopLevelExpr();

  std::unique_ptr<ExprAST> LogError(const char *Str);
  std::unique_ptr<PrototypeAST> LogErrorP(const char *Str);
};

