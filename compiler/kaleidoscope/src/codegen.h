#ifndef CODEGEN_H
#define CODEGEN_H

#include "../include/KaleidoscopeJIT.h"
#include "ast.h"
#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/PassManager.h>
#include <llvm/IR/Value.h>
#include <llvm/Passes/PassBuilder.h>
#include <llvm/Passes/StandardInstrumentations.h>
#include <llvm/Support/Error.h>

using namespace llvm;
using namespace orc;

extern std::unique_ptr<LLVMContext> TheContext;
extern std::unique_ptr<Module> TheModule;
extern std::unique_ptr<IRBuilder<>> Builder;
extern std::map<std::string, Value *> NamedValues;
extern std::unique_ptr<KaleidoscopeJIT> TheJIT;
extern std::unique_ptr<FunctionPassManager> TheFPM;
extern std::unique_ptr<LoopAnalysisManager> TheLAM;
extern std::unique_ptr<FunctionAnalysisManager> TheFAM;
extern std::unique_ptr<CGSCCAnalysisManager> TheCGAM;
extern std::unique_ptr<ModuleAnalysisManager> TheMAM;
extern std::unique_ptr<PassInstrumentationCallbacks> ThePIC;
extern std::unique_ptr<StandardInstrumentations> TheSI;
extern std::map<std::string, std::unique_ptr<PrototypeAST>> FunctionProtos;
extern ExitOnError ExitOnErr;

#endif // CODEGEN_H
