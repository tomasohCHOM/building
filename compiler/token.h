#pragma once

// The lexer returns tokens [0-255] if it is an unknown character, otherwise
// one of these for known things.
enum Token {
  tok_eof = -1,
  tok_def = -2,
  tok_extern = -3,
  tok_identifier = -4,
  tok_number = -5
};
