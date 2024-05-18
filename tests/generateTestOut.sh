#!/bin/bash

script_dir=`dirname $0`
files=`ls $script_dir/*.txt`
out_dir="$script_dir/out"
mkdir $out_dir

for i in $files;
do
	filename=`basename $i .txt`
	go run $script_dir/../cmd/main.go $i > "$out_dir/$filename.out"
done;

