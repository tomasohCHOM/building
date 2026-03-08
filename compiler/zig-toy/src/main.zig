const std = @import("std");
const Lexer = @import("lexer.zig").Lexer;

pub fn main() !void {
    const source =
        \\let x: int = 10;
        \\fn add(a: int, b: int) -> int {
        \\    return a + b;
        \\}
    ;
    var lexer = Lexer.init(source);
    while (true) {
        const tok = lexer.nextToken();
        std.debug.print("{any} \"{s}\"\n", .{ tok.kind, tok.lexeme });
        if (tok.kind == .eof) break;
    }
}
