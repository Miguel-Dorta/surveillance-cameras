#!/usr/bin/bash

if [ "$#" -ne 2 ]; then
        echo "Usage:    $(basename $0) <root_path> <cam_name>"
        exit 1
fi

PATH_RAW="$1/raw/$2"
PATH_SORT="$1/sort/$2/$(date -u '+%Y/%m/%d')"

mkdir -p "$PATH_SORT"
chmod -R 0777 "$PATH_SORT"
rm "$PATH_RAW"
ln -s "$PATH_SORT" "$PATH_RAW"
