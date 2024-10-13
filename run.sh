#!/bin/bash

rm readline || true
# rm /Users/patrik/sandbox/dnanexus/hw-read-lines/data/data.txt || true
# rm /Users/patrik/sandbox/dnanexus/hw-read-lines/data/large.txt || true
# rm /Users/patrik/sandbox/dnanexus/hw-read-lines/data/small.txt || true

go build .

# ./readline generate --path /Users/patrik/sandbox/dnanexus/hw-read-lines/data/data.txt --lines 1000000 --wordsPerLine 26
# ./readline generate --path /Users/patrik/sandbox/dnanexus/hw-read-lines/data/large.txt --lines 100000000 --wordsPerLine 30
# ./readline generate --path /Users/patrik/sandbox/dnanexus/hw-read-lines/data/small.txt --lines 10000 --wordsPerLine 30

./readline rand --verbose \
    --indexPath /Users/patrik/sandbox/dnanexus/hw-read-lines/index.idx \
    --path /Users/patrik/sandbox/dnanexus/hw-read-lines/data/large.txt
