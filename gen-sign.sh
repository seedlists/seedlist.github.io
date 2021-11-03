#!/usr/bin/env bash

#cd ./sha && make && cd -
rm *.csv*
./sha/sha1-mac html
Sum5=`shasum -a 256 sha1-each-file-in-html-dir.csv | awk '{print $1}'`
mv sha1-each-file-in-html-dir.csv sha1-each-file-in-html-dir-${Sum5}.csv
gpg --detach-sign  --armor sha1-each-file-in-html-dir-${Sum5}.csv
