# Filelist Mirror Downloader

This shim takes a file list and a mirror list and downloads the files in the file list
while doing checksums and parallel downloads for speed.

# Example usage:
CentOS / Yum repo
```bash
./filelist-mirror-downloader -mirrors yum_mirrorlist.txt -list yum_filelist.txt -output yumrepo_test
```

Alpine repo
```bash
./filelist-mirror-downloader -mirrors alpine_mirrorlist.txt -list alpine_filelist.txt -output alpine_test -debug -attempts 15
```

Threaded download:
```bash
./filelist-mirror-downloader -mirrors yum_mirrorlist.txt -list yum_filelist.txt -threads 2 -output yumrepo_test
```

Just test the validity of the downloads:
```bash
./filelist-mirror-downloader -list alpine_filelist.txt -output alpine_test -test ; echo $?
```

# Usage help:
```bash
$ ./filelist-mirror-downloader -h
Yum Get RepoMD,  Version: 0.1.20220315.1450

Usage: ./filelist-mirror-downloader [options...]

  -attempts int
        Attempts for each file (default 5)
  -debug
        Turn on debug comments
  -list string
        Filelist to be fetched (one per line with: HASH SIZE PATH) (default "filelist.txt")
  -mirrors string
        Mirror / directory list of prefixes to use (default "mirrorlist.txt")
  -output string
        Path to put the repo files (default ".")
  -shuffle int
        Shuffle the mirror list ever N downloads (default 10)
  -threads int
        Concurrent downloads (default 1)
```
