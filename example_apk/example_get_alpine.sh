#!/bin/bash

output=/dev/shm/alpine_repo
mirrors=alpine_mirrorlist.txt

# An example Alpine downloader

for repo in latest-stable/{main,community}/x86_64 v3.15/main/x86_64; do
  echo Getting latest APKINDEX.tar.gz for $repo
  ../../alpine-get-repomd/alpine-get-repomd -mirrors $mirrors -output $output/$repo/ -repo $repo -keysDir ../../alpine-get-repomd/keys/
done

echo Creating a filelist for downloading
for repo in latest-stable/{main,community}/x86_64 v3.15/main/x86_64; do
  ../../alpine-packages-diff/alpine-package-diff -repo $repo -new $output/$repo/APKINDEX.tar.gz -old "" -showAdded
done > alpine_filelist.txt

echo Download the files to a local folder
../filelist-mirror-downloader -mirrors $mirrors -list alpine_filelist.txt -output $output -log alpine_filelist.log 
