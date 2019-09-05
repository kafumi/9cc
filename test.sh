#!/usr/bin/env bash

try() {
  expected="$1"
  input="$2"

  ./9cc "$input" > tmp.s
  gcc -o tmp tmp.s test/*.o
  ./tmp
  actual="$?"
  rm -f tmp

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

try   0 'main(){ 0; }'
try  42 'main(){ 42; }'
try  21 'main(){ 5+20-4; }'
try  41 'main(){  12 + 34 - 5 ; }'
try  47 'main(){ 5+6*7; }'
try  50 'main(){ 38+3*8/2; }'
try  15 'main(){ 5*(9-6); }'
try   4 'main(){ (3+5)/2; }'
try   4 'main(){ -(3+5) / -2; }'
try   5 'main(){ -3 * +5 - -20; }'
try  10 'main(){ -10 + 20; }'
try   1 'main(){ 1+2==2*4-5; }'
try   0 'main(){ 1+2!=2*4-5; }'
try   0 'main(){ 3*4==-3*4; }'
try   1 'main(){ 3*4!=-3*4; }'
try   1 'main(){ 36/3<36/2; }'
try   1 'main(){ 36/3<=36/2; }'
try   0 'main(){ 5*6<90/3; }'
try   1 'main(){ 5*6<=90/3; }'
try   1 'main(){ 48/3>48/4; }'
try   1 'main(){ 48/3>=48/4; }'
try   0 'main(){ 3*8>4*6; }'
try   1 'main(){ 3*8>=4*6; }'
try   5 'main(){ int a; a=5; a; }'
try   8 'main(){ int b; b=4; b*2; }'
try   8 'main(){ int c; c=4*2; c; }'
try  14 'main(){ int a; int b; a=3; b=5*6-8; a+b/2; }'
try  23 'main(){ int foo; int bar; foo=3; bar=4*5; foo=foo+bar; foo; }'
try   2 'main(){ int a1; int a2; a1=10; a2=20; a2/a1; }'
try   5 'main(){ return 5; return 8; }'
try  10 'main(){ int return1; return1=7; return1+3; }'
try   5 'main(){ int a; a=1; if(3>a) a=5; a; }'
try   1 'main(){ int a; a=1; if(a>3) a=5; a; }'
try   5 'main(){ int a; a=1; if(3>a) a=5; else a=7; a; }'
try   7 'main(){ int a; a=1; if(a>3) a=5; else a=7; a; }'
try  16 'main(){ int a; a=1; while(a<10) a=a*2; a; }'
try   1 'main(){ int a; a=1; while(a>10) a=a*2; a; }'
try   8 'main(){ int a; a=1; for(; a<5; ) a=a*2; a; }'
try  21 'main(){ int a; int b; b=1; for(a=63; a>10; a=a/3) b=b+1; a*b; }'
try 135 'main(){ int a; for (a=5; a<100; a=a*3) {} a; }'
try   8 'main(){ int a; int b; a=1; b=2; if(a<2) {a=a+1; b=b+2;} a*b; }'
try   2 'main(){ 1+func0(); }'
try   3 'main(){ 1+func1(1); }'
try   5 'main(){ 1+func2(1, 2); }'
try   8 'main(){ 1+func3(1, 2, 3); }'
try  12 'main(){ 1+func4(1, 2, 3, 4); }'
try  17 'main(){ 1+func5(1, 2, 3, 4, 5); }'
try  23 'main(){ 1+func6(1, 2, 3, 4, 5, 6); }'
try   4 'main(){ 1+add(1, 2); } add(a, b){ a + b; }'
try   7 'mul(a, b){ a * b; } main(){ 1+mul(2, 3); }'
try   0 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(0); }'
try   1 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(1); }'
try   1 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(2); }'
try   2 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(3); }'
try   3 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(4); }'
try   5 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(5); }'
try   8 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(6); }'
try  13 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(7); }'
try  21 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(8); }'
try  34 'fib(n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } main(){ fib(9); }'
try   3 'main(){ int x; int y; x=3; y=&x; *y; }'
try   3 'main(){ int x; int y; int z; x=3; y=5; z=&y+8; *z; }'

echo OK
