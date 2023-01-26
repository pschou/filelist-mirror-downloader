#!/bin/bash

output=/dev/shm/rancher_repo
mirrors=rancher_mirrorlist.txt
repos=(rke2/stable/common/centos/{7,8}/noarch rke2/stable/1.18/centos/{7,8}/x86_64)
file_list=$output/rancher_filelist.txt

# An example Yum downloader for rancher

for repo in "${repos[@]}"; do
  echo Getting latest index for $repo
  ../../yum-get-repomd/yum-get-repomd -insecure -mirrors $mirrors -output $output/$repo/ -repo $repo -keyring rancher_keys/
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done

echo Creating a filelist for downloading
for repo in "${repos[@]}"; do
  ../../yum-packages-diff/yum-package-diff -repo $repo -new $output/$repo/ -old "" -showAdded
  if [ $? == 1 ]; then echo Error encountered; exit; fi
done > $file_list

echo Download the files to a local folder
../filelist-mirror-downloader -mirrors $mirrors -list $file_list -output $output -keyring rancher_keys/
if [ $? == 1 ]; then echo Error encountered; exit; fi
