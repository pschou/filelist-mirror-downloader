#!/bin/bash

output=./alpine_repo/
mirrors=alpine_mirrorlist.txt

# An example Alpine downloader!

echo Getting latest APKINDEX.tar.gz
../alpine-get-repomd/alpine-get-repomd -mirrors $mirrors -output $output -debug -keysDir ../alpine-get-repomd/keys/

echo Creating a filelist for downloading
../alpine-packages-diff/alpine-package-diff -new alpine_repo/APKINDEX.tar.gz -old "" -showAdded -output alpine_filelist.txt

echo Download the files to a local folder
./filelist-mirror-downloader -mirrors $mirrors -list alpine_filelist.txt -output /dev/shm/alpine_test
