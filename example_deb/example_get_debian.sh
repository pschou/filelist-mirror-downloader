#!/bin/bash

output=/dev/shm/debian_repo
mirrors=debian_mirrorlist.txt
repos=(dists/bookworm/{main,contrib,non-free,non-free-firmware}/binary-amd64)
file_list=$output/debian_filelist.txt

# An example Yum downloader for debian

for repo in "${repos[@]}"; do
  echo Getting latest index for $repo
  ../../deb-get-repomd/deb-get-repomd -mirrors $mirrors -keyring debian_keys -output $output -tree -repo "$repo"
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done

echo Creating a filelist for downloading
for repo in "${repos[@]}"; do
  ../../deb-packages-diff/deb-package-diff -repo "$repo" -new $output/"$repo" -showAdded -old ''
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done > $file_list

echo Download the files to a local folder
../filelist-mirror-downloader -mirrors $mirrors -list $file_list -output $output -keyring debian_keys/
if [ $? == 1 ]; then echo Error encountered; exit; fi
