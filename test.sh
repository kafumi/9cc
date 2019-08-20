#!/usr/bin/env bash

try() {
  expected="$1"
  input="$2"

  ./9cc "$input" > tmp.s
  gcc -o tmp tmp.s
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

try 0 0
try 42 42
try 21 "5+20-4"
try 41 " 12 + 34 - 5 "
try 47 '5+6*7'
try 50 '38+3*8/2'
try 15 '5*(9-6)'
try 4 '(3+5)/2'
try 4 '-(3+5) / -2'
try 5 '-3 * +5 - -20'
try 10 '-10 + 20'

echo OK
