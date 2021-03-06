#! /usr/bin/env bash

for f in `ls -1 *.prod`; do
    file=`basename $f`
    reffile=ref/$file.ref
    outfile=monster.out
    echo ".........", $file
    monster -bagdir ../bags -count 1 -seed 1 $f 1 > $outfile 2>&1;
    diff $outfile $reffile
    echo
done
exit 0
