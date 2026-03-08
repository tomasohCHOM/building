const std = @import("std");
const token = @import("token.zig");

const TokenKind = token.TokenKind;
const Token = token.Token;

pub const Lexer = struct {
    const Self = @This();

    source: []const u8,
    pos: usize,
    line: usize,
    column: usize,

    pub fn init(source: []const u8) Self {
        return .{
            .source = source,
            .pos = 0,
            .line = 1,
            .column = 1,
        };
    }

    fn peek(self: *const Self) ?u8 {
        if (self.pos >= self.source.len) return null;
        return self.source[self.pos];
    }

    fn advance(self: *Self) ?u8 {
        if (self.pos >= self.source.len) return null;

        const c = self.source[self.pos];
        self.pos += 1;

        if (c == '\n') {
            self.line += 1;
            self.column = 1;
        } else {
            self.column += 1;
        }

        return c;
    }

    fn match(self: *Self, expected: u8) bool {
        if (self.peek() != expected) return false;
        _ = self.advance();
        return true;
    }

    pub fn nextToken(self: *Self) Token {
        while (true) {
            const start = self.pos;
            const line = self.line;
            const column = self.column;

            const c = self.advance() orelse {
                return .{
                    .kind = .eof,
                    .lexeme = "",
                    .line = line,
                    .column = column,
                };
            };

            switch (c) {
                ' ', '\t', '\r', '\n' => continue,

                ':' => return self.makeToken(.colon, start, line, column),
                ';' => return self.makeToken(.semicolon, start, line, column),

                '(' => return self.makeToken(.l_paren, start, line, column),
                ')' => return self.makeToken(.r_paren, start, line, column),
                '{' => return self.makeToken(.l_brace, start, line, column),
                '}' => return self.makeToken(.r_brace, start, line, column),

                '+' => return self.makeToken(.plus, start, line, column),
                '-' => {
                    if (self.match('>'))
                        return self.makeToken(.arrow, start, line, column);
                    return self.makeToken(.minus, start, line, column);
                },
                '*' => return self.makeToken(.star, start, line, column),
                '/' => return self.makeToken(.slash, start, line, column),

                '=' => {
                    if (self.match('='))
                        return self.makeToken(.equal_equal, start, line, column);
                    return self.makeToken(.equal, start, line, column);
                },

                '!' => {
                    if (self.match('='))
                        return self.makeToken(.bang_equal, start, line, column);
                },

                '<' => {
                    if (self.match('='))
                        return self.makeToken(.less_equal, start, line, column);
                    return self.makeToken(.less, start, line, column);
                },

                '>' => {
                    if (self.match('='))
                        return self.makeToken(.greater_equal, start, line, column);
                    return self.makeToken(.greater, start, line, column);
                },

                else => {
                    if (std.ascii.isDigit(c))
                        return self.number(start, line, column);

                    if (std.ascii.isAlphabetic(c) or c == '_') {
                        while (true) {
                            const next_c = self.peek() orelse break;
                            if (!std.ascii.isAlphanumeric(next_c) and next_c != '_')
                                break;
                            _ = self.advance();
                        }

                        const text = self.source[start..self.pos];
                        const kind = tokenKind(text);

                        return .{
                            .kind = kind,
                            .lexeme = text,
                            .line = line,
                            .column = column,
                        };
                    }

                    return .{
                        .kind = .err,
                        .lexeme = self.source[start..self.pos],
                        .line = line,
                        .column = column,
                    };
                },
            }
        }
    }

    fn makeToken(
        self: Self,
        kind: TokenKind,
        start: usize,
        line: usize,
        column: usize,
    ) Token {
        return .{
            .kind = kind,
            .lexeme = self.source[start..self.pos],
            .line = line,
            .column = column,
        };
    }

    fn number(self: *Self, start: usize, line: usize, column: usize) Token {
        while (true) {
            const d = self.peek() orelse break;
            if (!std.ascii.isDigit(d)) break;
            _ = self.advance();
        }

        return .{
            .kind = .int_literal,
            .lexeme = self.source[start..self.pos],
            .line = line,
            .column = column,
        };
    }
};

fn tokenKind(text: []const u8) TokenKind {
    if (std.mem.eql(u8, text, "let")) return .keyword_let;
    if (std.mem.eql(u8, text, "fn")) return .keyword_fn;
    if (std.mem.eql(u8, text, "return")) return .keyword_return;
    if (std.mem.eql(u8, text, "if")) return .keyword_if;
    if (std.mem.eql(u8, text, "else")) return .keyword_else;
    if (std.mem.eql(u8, text, "int")) return .keyword_int;
    if (std.mem.eql(u8, text, "true")) return .keyword_true;
    if (std.mem.eql(u8, text, "false")) return .keyword_false;

    return .identifier;
}

const expect = std.testing.expect;
test "Lexer" {
    const Helper = struct {
        fn verifyTokenWithKind(
            curr_token: Token,
            expected_token: []const u8,
            expected_kind: TokenKind,
        ) bool {
            return std.mem.eql(u8, curr_token.lexeme, expected_token) and
                curr_token.kind == expected_kind;
        }
    };

    const source = "let x: int = 0;";
    const expected_tokens = [_]struct { []const u8, TokenKind }{
        .{ "let", .keyword_let },
        .{ "x", .identifier },
        .{ ":", .colon },
        .{ "int", .keyword_int },
        .{ "=", .equal },
        .{ "0", .int_literal },
        .{ ";", .semicolon },
    };

    var lexer = Lexer.init(source);
    var curr_token: Token = undefined;

    for (expected_tokens) |expected_token| {
        curr_token = lexer.nextToken();
        const lexeme, const token_kind = expected_token;
        try expect(Helper.verifyTokenWithKind(curr_token, lexeme, token_kind));
    }
}
