# 9cc: Yet another C compiler

[低レイヤを知りたい人のためのCコンパイラ作成入門](https://www.sigbus.info/compilerbook)

## EBNF

```ebnf
program    = funct*
funct      = ident "(" (ident ("," ident)*)? ")" "{" stmt* "}"
stmt       = expr ";"
           | "{" stmt* "}"
           | "if" "(" expr ")" stmt ("else" stmt)?
           | "while" "(" expr ")" stmt
           | "for" "(" expr? ";" expr? ";" expr? ")" stmt
           | "return" expr ";"
           | "int" ident ";"
expr       = assign
assign     = equality ("=" assign)?
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = "+"? primary
           | "-"? primary
           | "&" unary
           | "*" unary
primary    = num
           | ident ("(" (expr ("," expr)*)? ")")?
           | "(" expr ")"
```