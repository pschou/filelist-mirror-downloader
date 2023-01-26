#!/bin/bash

output=/dev/shm/ubuntu_repo
mirrors=ubuntu_mirrorlist.txt
repos=(dists/xenial/{main,restricted,universe,multiverse}/binary-amd64 dists/bionic/{main,restricted,universe,multiverse}/binary-amd64)
file_list=$output/ubuntu_filelist.txt

# An example Yum downloader for ubuntu

for repo in "${repos[@]}"; do
  echo Getting latest index for $repo
  ../../deb-get-repomd/deb-get-repomd -mirrors $mirrors -keyring ubuntu_keys -output $output -tree -repo "$repo"
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done

echo Creating a filelist for downloading
for repo in "${repos[@]}"; do
  ../../deb-packages-diff/deb-package-diff -repo "$repo" -new $output/"$repo" -showAdded -old ''
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done > $file_list

echo Download the files to a local folder
../filelist-mirror-downloader -mirrors $mirrors -list $file_list -output $output -keyring ubuntu_keys/
if [ $? == 1 ]; then echo Error encountered; exit; fi
