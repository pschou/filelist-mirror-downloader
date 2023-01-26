#!/bin/bash

output=/dev/shm/theforeman_repo
mirrors=theforeman_mirrorlist.txt
repos=latest/el8/x86_64
file_list=$output/foreman_filelist.txt
user=myself
pass=mypass

# An example Yum downloader

for repo in $repos; do
  echo Getting latest index for $repo
  ../../yum-get-repomd/yum-get-repomd --insecure -mirrors $mirrors -output $output/$repo/ -repo $repo -keyring theforeman_keys/ -client-user "$user" -client-pass "$pass"
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done

echo Creating a filelist for downloading
for repo in $repos; do
  ../../yum-packages-diff/yum-package-diff -repo $repo -new $output/$repo/ -old "" -showAdded
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done > $file_list

echo Download the files to a local folder
../filelist-mirror-downloader -mirrors $mirrors -list $file_list -output $output -client-user "$user" -client-pass "$pass" -keyring theforeman_keys/
if [ $? == 1 ]; then echo Error encountered; exit; fi
