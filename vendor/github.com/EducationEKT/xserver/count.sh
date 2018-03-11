#!bin/bash
find ./ -name "*.go" | xargs cat | wc -l
