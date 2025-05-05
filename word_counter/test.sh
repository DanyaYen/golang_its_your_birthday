#!/bin/bash


echo -e "hello world\nGo is awesome\n123" > test.txt


read lines words bytes filename <<< $(wc test.txt)


actual=$(go run main.go test.txt)


expected="$lines $words $bytes test.txt"


if [ "$actual" == "$expected" ]; then
    echo "✅ Test passed"
else
    echo "❌ Test failed"
    echo "Expected: $expected"
    echo "Got     : $actual"
fi
