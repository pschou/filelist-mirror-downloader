#!/bin/bash

./filelist-mirror-downloader -mirrors yum_mirrorlist.txt -list yum_filelist.txt -threads 2 -output yumrepo_test2 -log logfile.txt
#./filelist-mirror-downloader -mirrors yum_mirrorlist.txt -list yum_filelist.txt -threads 2 -output yumrepo_test2 -log logfile.txt -after 2022-01-01

while read line; do
  kind=${line%: *};
  file=${line#*: };
  case $kind in 
    OnDisk)
      echo found file on disk: $file
      ;;
    OnDiskSkip)
      echo found file on disk outside date range: $file
      ;;
    Skipped)
      echo date range excluded: $file
      ;;
    Fetched)
      echo file downloaded: $file
      ;;
  esac
done < logfile.txt
