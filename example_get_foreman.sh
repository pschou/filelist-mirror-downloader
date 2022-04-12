#!/bin/bash

output=/dev/shm/theforeman_repo
mirrors=theforeman_mirrorlist.txt
repos=latest/el7/x86_64
file_list=foreman_filelist.txt

# An example Yum downloader

for repo in $repos; do
  echo Getting latest index for $repo
  echo ../yum-get-repomd/yum-get-repomd --insecure -mirrors $mirrors -output $output/$repo/ -repo $repo -keyring theforeman_keys/
  ../yum-get-repomd/yum-get-repomd --insecure -mirrors $mirrors -output $output/$repo/ -repo $repo -keyring theforeman_keys/
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done

echo Creating a filelist for downloading
for repo in $repos; do
  ../yum-packages-diff/yum-package-diff -repo $repo -new $output/$repo/ -old "" -showAdded
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done > $file_list

echo Download the files to a local folder
./filelist-mirror-downloader -mirrors $mirrors -list $file_list -output $output 
if [ $? == 1 ]; then echo Error encountered; exit; fi
