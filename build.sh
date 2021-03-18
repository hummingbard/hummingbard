#!/bin/bash

set -e
cd "$(dirname "$0")"

rm -rf static/css/*.css
echo "Minifying CSS files..."
cd ui/css
for i in $( find . -type f -name '*\.css' ); do
    name=`basename $i | cut -f 1 -d '.'`
    id=`cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1`
    uglifycss $i > ../../static/css/$name.$id.css
done

