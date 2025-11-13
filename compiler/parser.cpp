#include "parser.h"
#include <cstdio>

// LogError* - These are little helper functions for error handling.
std::unique_ptr<ExprAST> Parser::LogError(const char *Str) {
  fprintf(stderr, "Error: %s\n", Str);
  return nullptr;
}

std::unique_ptr<PrototypeAST> Parser::LogErrorP(const char *Str) {
  LogError(Str);
  return nullptr;
}

std::unique_ptr<ExprAST> ParseExpression();

// numberexpr ::= number
std::unique_ptr<ExprAST> Parser::ParseNumberExpr() {
  auto Result = std::make_unique<NumberExprAST>(Lex.getNumVal());
  getNextToken();
  return std::move(Result);
}

// parenexpr ::= '(' expression ')'
std::unique_ptr<ExprAST> Parser::ParseParenExpr() {
  getNextToken(); // eat (.
  auto V = ParseExpression();
  if (!V)
    return nullptr;

  if (CurTok != ')')
    return LogError("expected ')'");
  getNextToken(); // eat ).
  return V;
}

// identifierexpr
//   ::= identifier
//   ::= identifier '(' expression* ')'
std::unique_ptr<ExprAST> Parser::ParseIdentifierExpr() {
  std::string IdName = Lex.getIdentifierStr();

  getNextToken();  // eat identifier.

  if (CurTok != '(') // Simple variable ref.
    return std::make_unique<VariableExprAST>(IdName);

  // Call.
  getNextToken();  // eat (
  std::vector<std::unique_ptr<ExprAST>> Args;
  if (CurTok != ')') {
    while (true) {
      if (auto Arg = ParseExpression())
        Args.push_back(std::move(Arg));
      else
        return nullptr;

      if (CurTok == ')')
        break;

      if (CurTok != ',')
        return LogError("Expected ')' or ',' in argument list");
      getNextToken();
    }
  }

  // Eat the ')'.
  getNextToken();

  return std::make_unique<CallExprAST>(IdName, std::move(Args));
}
