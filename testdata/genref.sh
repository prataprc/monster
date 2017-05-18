#! /usr/bin/env bash

dopath() {
    for f in `find $1 -name "*.prod"`; do
        file=`basename $f`
        outfile=ref/$file.ref
        echo $outfile ...
        monster -bagdir ../bags -count 1 -seed 1 $f > $outfile;
    done
}

dopath `pwd`
