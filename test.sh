#!/usr/bin/env bash

try() {
  expected="$1"
  input="$2"

  ./9cc "$input" > tmp.s
  gcc -static -o tmp tmp.s test/*.o
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

try   0 'int main(){ 0; }'
try  42 'int main(){ 42; }'
try  21 'int main(){ 5+20-4; }'
try  41 'int main(){  12 + 34 - 5 ; }'
try  47 'int main(){ 5+6*7; }'
try  50 'int main(){ 38+3*8/2; }'
try  15 'int main(){ 5*(9-6); }'
try   4 'int main(){ (3+5)/2; }'
try   4 'int main(){ -(3+5) / -2; }'
try   5 'int main(){ -3 * +5 - -20; }'
try  10 'int main(){ -10 + 20; }'
try   1 'int main(){ 1+2==2*4-5; }'
try   0 'int main(){ 1+2!=2*4-5; }'
try   0 'int main(){ 3*4==-3*4; }'
try   1 'int main(){ 3*4!=-3*4; }'
try   1 'int main(){ 36/3<36/2; }'
try   1 'int main(){ 36/3<=36/2; }'
try   0 'int main(){ 5*6<90/3; }'
try   1 'int main(){ 5*6<=90/3; }'
try   1 'int main(){ 48/3>48/4; }'
try   1 'int main(){ 48/3>=48/4; }'
try   0 'int main(){ 3*8>4*6; }'
try   1 'int main(){ 3*8>=4*6; }'
try   5 'int main(){ int a; a=5; a; }'
try   8 'int main(){ int b; b=4; b*2; }'
try   8 'int main(){ int c; c=4*2; c; }'
try  14 'int main(){ int a; int b; a=3; b=5*6-8; a+b/2; }'
try  23 'int main(){ int foo; int bar; foo=3; bar=4*5; foo=foo+bar; foo; }'
try   2 'int main(){ int a1; int a2; a1=10; a2=20; a2/a1; }'
try   5 'int main(){ return 5; return 8; }'
try  10 'int main(){ int return1; return1=7; return1+3; }'
try   5 'int main(){ int a; a=1; if(3>a) a=5; a; }'
try   1 'int main(){ int a; a=1; if(a>3) a=5; a; }'
try   5 'int main(){ int a; a=1; if(3>a) a=5; else a=7; a; }'
try   7 'int main(){ int a; a=1; if(a>3) a=5; else a=7; a; }'
try  16 'int main(){ int a; a=1; while(a<10) a=a*2; a; }'
try   1 'int main(){ int a; a=1; while(a>10) a=a*2; a; }'
try   8 'int main(){ int a; a=1; for(; a<5; ) a=a*2; a; }'
try  21 'int main(){ int a; int b; b=1; for(a=63; a>10; a=a/3) b=b+1; a*b; }'
try 135 'int main(){ int a; for (a=5; a<100; a=a*3) {} a; }'
try   8 'int main(){ int a; int b; a=1; b=2; if(a<2) {a=a+1; b=b+2;} a*b; }'
try   2 'int main(){ 1+func0(); }'
try   3 'int main(){ 1+func1(1); }'
try   5 'int main(){ 1+func2(1, 2); }'
try   8 'int main(){ 1+func3(1, 2, 3); }'
try  12 'int main(){ 1+func4(1, 2, 3, 4); }'
try  17 'int main(){ 1+func5(1, 2, 3, 4, 5); }'
try  23 'int main(){ 1+func6(1, 2, 3, 4, 5, 6); }'
try   4 'int main(){ 1+add(1, 2); } int add(int a, int b){ a + b; }'
try   7 'int mul(int a, int b){ a * b; } int main(){ 1+mul(2, 3); }'
try   0 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(0); }'
try   1 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(1); }'
try   1 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(2); }'
try   2 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(3); }'
try   3 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(4); }'
try   5 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(5); }'
try   8 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(6); }'
try  13 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(7); }'
try  21 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(8); }'
try  34 'int fib(int n){ if(n==0) return 0; if(n==1) return 1; fib(n-2)+fib(n-1); } int main(){ fib(9); }'
try   3 'int main(){ int x; int *y; x=3; y=&x; *y; }'
try   3 'int main(){ int x; int y; int *z; x=3; y=5; z=&y+1; *z; }'
try   3 'int main(){ int x; int *y; y=&x; *y=3; x; }'
try   4 'int main(){ int *p; int *q; alloc4(&p, 1, 2, 4, 8); q=p+2; *q; }'
try   8 'int main(){ int *p; int *q; alloc4(&p, 1, 2, 4, 8); q=3+p; *q; }'
try   4 'int main(){ int x; sizeof(x); }'
try   8 'int main(){ int *x; sizeof(x); }'
try   4 'int main(){ int x; sizeof(x + 1); }'
try   8 'int main(){ int *x; sizeof(x + 2); }'
try   4 'int main(){ int *x; sizeof(*x); }'
try   4 'int main(){ sizeof(1); }'
try   4 'int main(){ sizeof(sizeof(1)); }'
try  16 'int main(){ int a[4]; sizeof(a); }'
try   4 'int main(){ int a[4]; sizeof(*a); }'
try  32 'int main(){ int *a[4]; sizeof(a); }'
try   8 'int main(){ int *a[4]; sizeof(*a); }'
try   3 'int main(){ int a[2]; *a = 1; *(a+1) = 2; int *p; p = a; *p + *(p+1); }'
try   6 'int main(){ int a[3]; a[0]=1; a[1]=2; a[2]=3; a[0]+a[1]+a[2]; }'
try  15 'int main(){ int a[3]; 0[a]=4; a[1]=5; 2[a]=6; a[0]+1[a]+a[2]; }'
try   5 'int a; int main(){ a=5; a; }'
try   4 'int a[3]; int main(){ a[0]=1; a[2]=3; a[0]+a[1]+a[2]; }'
try   3 'char a[3]; int main(){ a[0]=-1; a[2]=2; int b; b=4; a[0]+b; }'
try   7 'int add(char a, char b){ a + b; } int main(){ add(-3, 10); }'

echo OK
