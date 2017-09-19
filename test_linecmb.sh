#!/bin/bash

set -e
executefile=./linecmb

echo "Trivial test"
${executefile} 2> /dev/null

#--------------------------------
echo "Basic test"
SIZE=$((1024*63))
perl -se '
print "0"x$shell_var, "\n";
print "1"x$shell_var, "\n";
print "2"x$shell_var, "\n";
print "3"x$shell_var, "\n";
print "4"x$shell_var, "\n";
print "5"x$shell_var, "\n";
print "6"x$shell_var, "\n";
print "7"x$shell_var, "\n";
print "8"x$shell_var, "\n";
print "9"x$shell_var, "\n";
' -- -shell_var=${SIZE} > test_sample.out
for((i=0; i<20; ++i)) {
  printf "# iteration $i of 20";

  time ${executefile}  6 5 7 8 9 10 11 12 13 14 15 16 17 \
    5<  <( sleep 0; perl -se 'print "0"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    6<  <( sleep 0; perl -se 'print "1"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    7<  <( sleep 0; perl -se 'print "2"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    8<  <( sleep 0; perl -se 'print "3"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    9<  <( sleep 0; perl -se 'print "4"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    10< <( sleep 0; perl -se 'print "5"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    11< <( sleep 0; perl -se 'print "6"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    12< <( sleep 0; perl -se 'print "7"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    13< <( sleep 0; perl -se 'print "8"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    14< <( sleep 0; perl -se 'print "9"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    > test.out;
  sort -n test.out | cmp test_sample.out -
  
}

echo OK!

#--------------------------------
echo Interleaved test
SIZE=$((1024*64))
perl -se 'print "0\n"x$shell_var, "1\n"x$shell_var, "2\n"x$shell_var' -- -shell_var=${SIZE} > test_sample.out
for((i=0; i<10; ++i)) {
  printf "# iteration $i of 10";

  time ${executefile}  6 5 7 \
    5< <( sleep 0; perl -se 'print "0\n"x$shell_var' -- -shell_var=${SIZE}) \
    6< <( sleep 0; perl -se 'print "1\n"x$shell_var' -- -shell_var=${SIZE}) \
    7< <( sleep 0; perl -se 'print "2\n"x$shell_var' -- -shell_var=${SIZE}) \
    > test.out;

  sort test.out | cmp test_sample.out -
}
echo OK!



#--------------------------------
echo "Non-ASCII & Random FD Num test"
SIZE=$((1024*8))
perl -se '
print "あ"x$shell_var, "\n";
print "0"x$shell_var, "\n";
print "10"x$shell_var, "\n";
' -- -shell_var=${SIZE} > test_sample.out
for((i=0; i<20; ++i)) {
  printf "# iteration $i of 20";

  time ${executefile}  6 9 12 \
    9<  <( sleep 0; perl -se 'print "あ"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    6<  <( sleep 0; perl -se 'print "0"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    12< <( sleep 0; perl -se 'print "10"x$shell_var, "\n"' -- -shell_var=${SIZE}) \
    > test.out;

  cmp <(sort -n test_sample.out) <(sort -n test.out)
}
echo OK!


#--------------------------------
echo "Long line test"
echo " this case will take long time"

cat <( sleep 0; perl -e 'print "0\n" foreach 1..(35*1024)' ) > test_sample.out
cat <( sleep 0; perl -e 'print "1"x(1024*1024) foreach 1..31'; echo ) >> test_sample.out
cat <( sleep 0; perl -e 'print "2\n" foreach 1..(35*1024)' ) >> test_sample.out

for((i=1; i<=2; ++i)) {
  echo "    iteration $i of 2";

  (ulimit -v 10000000;time ${executefile}  6 5 7) \
    5< <( sleep 0; perl -e 'print "0\n" foreach 1..(35*1024)' ) \
    6< <( sleep 0; perl -e 'print "1"x(1024*1024) foreach 1..31'; echo ) \
    7< <( sleep 0; perl -e 'print "2\n" foreach 1..(35*1024)' ) \
    > test.tmp.out
  cat test.tmp.out | grep 0 > test.out
  cat test.tmp.out | grep 1 >> test.out
  cat test.tmp.out | grep 2 >> test.out

  cmp test_sample.out test.out

}

echo OK!

rm test_sample.out test.out test.tmp.out
