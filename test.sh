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

try 0 '0;'
try 42 '42;'
try 21 '5+20-4;'
try 41 ' 12 + 34 - 5 ;'
try 47 '5+6*7;'
try 50 '38+3*8/2;'
try 15 '5*(9-6);'
try 4 '(3+5)/2;'
try 4 '-(3+5) / -2;'
try 5 '-3 * +5 - -20;'
try 10 '-10 + 20;'
try 1 '1+2==2*4-5;'
try 0 '1+2!=2*4-5;'
try 0 '3*4==-3*4;'
try 1 '3*4!=-3*4;'
try 1 '36/3<36/2;'
try 1 '36/3<=36/2;'
try 0 '5*6<90/3;'
try 1 '5*6<=90/3;'
try 1 '48/3>48/4;'
try 1 '48/3>=48/4;'
try 0 '3*8>4*6;'
try 1 '3*8>=4*6;'
try 5 'a=5; a;'
try 8 'b=4; b*2;'
try 8 'c=4*2; c;'
try 14 'a=3; b=5*6-8; a+b/2;'
try 23 'foo=3; bar=4*5; foo=foo+bar; foo;'
try 2 'a1=10; a2=20; a2/a1;'
try 5 'return 5; return 8;'
try 10 'return1=7; return1+3;'
try 5 'a=1; if(3>a) a=5; a;'
try 1 'a=1; if(a>3) a=5; a;'
try 5 'a=1; if(3>a) a=5; else a=7; a;'
try 7 'a=1; if(a>3) a=5; else a=7; a;'
try 16 'a=1; while(a<10) a=a*2; a;'
try 1 'a=1; while(a>10) a=a*2; a;'
try 8 'a=1; for(; a<5; ) a=a*2; a;'
try 21 'b=1; for(a=63; a>10; a=a/3) b=b+1; a*b;'
try 135 'for (a=5; a<100; a=a*3) {} a;'
try 8 'a=1; b=2; if(a<2) {a=a+1; b=b+2;} a*b;'
try 2 '1+func0();'

echo OK
