#!/bin/sh

input1=./sample.d/input1.gz
input2=./sample.d/input2.gz

geninput(){
	echo generating input files...

	echo hw | gzip > "${input1}"
	echo wl | gzip > "${input2}"
}

test -f "${input1}" || geninput

ls ./sample.d/input*.gz |
	./gzstat |
	jq -c
