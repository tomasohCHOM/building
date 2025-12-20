#ifndef CODEGEN_H
#define CODEGEN_H

#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/Value.h>

using namespace llvm;

extern std::unique_ptr<LLVMContext> TheContext;
extern std::unique_ptr<Module> TheModule;
extern std::unique_ptr<IRBuilder<>> Builder;

#endif // CODEGEN_H
