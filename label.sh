#!/bin/bash

set -e

path="$1"

mkdir -p tmp
find $path -type f -name *.JPG > tmp/media

function resize () {
    ffmpeg -loglevel quiet -nostats -hide_banner -y -i "$1" -q:v 1 -vf scale="1920:trunc(ow/a/2)*2" "$2"
}

function caption () {
	tmp_img=/tmp/`basename "$1"`
    resize "$1" "$tmp_img"
    curl -s -X POST -F "image=@${tmp_img}" http://localhost:5000/upload | jq -r '.message.["<MORE_DETAILED_CAPTION>"]'
    rm "$tmp_img"
}

total=$(cat tmp/media | wc -l)
count=1

while IFS= read -r -u4 p; do
    echo "[+] Processing [$count:$total] $p"
    
    c=$(caption "$p")
    
    echo "$c"

    exiftool.exe -overwrite_original -description="$c" "`wslpath -w "$p"`"

    echo ""

    count=$[ $count +1 ]
done 4<tmp/media