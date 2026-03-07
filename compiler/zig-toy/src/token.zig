const std = @import("std");

pub const TokenKind = enum {
    identifier,
    int_literal,

    keyword_let,
    keyword_fn,
    keyword_return,
    keyword_if,
    keyword_else,
    keyword_true,
    keyword_false,

    colon,
    semicolon,
    comma,

    l_paren,
    r_paren,
    l_brace,
    r_brace,

    arrow,

    plus,
    minus,
    star,
    slash,

    equal,
    equal_equal,
    bang_equal,

    less,
    less_equal,
    greater,
    greater_equal,

    and_and,
    or_or,

    eof,
};

pub const Token = struct {
    kind: TokenKind,
    lexeme: []const u8,
    line: usize,
    column: usize,
};
