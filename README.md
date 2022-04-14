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
# Test the downloads, note that the output is in the filelist format for missing or invalid files
./filelist-mirror-downloader -list alpine_filelist.txt -output alpine_test -test
# The return value will be 0 for success or 1 for any failures
echo $?
```

# Usage help:
```bash
$ ./filelist-mirror-downloader -h
Yum Get RepoMD,  Version: 0.1.20220413.2147

Usage: ./filelist-mirror-downloader [options...]

  -after string
        Select packages after specified date
        Date formats supported: https://github.com/araddon/dateparse
  -attempts int
        Attempts for each file (default 40)
  -before string
        Select packages before specified date
        Date formats supported: https://github.com/araddon/dateparse
  -debug
        Turn on debug comments
  -dup string
        What to do with duplicates: omit, copy, symlink, hardlink (default "symlink")
  -list string
        Filelist to be fetched (one per line with: HASH SIZE PATH) (default "filelist.txt")
  -log string
        File in which to store a log of files downloaded
        Line prefixes (OnDisk, OnDiskSkip, Skipped, Fetched), indicate action taken.
        Skip means that a file falls outside the required date bounds
  -mirrors string
        Mirror / directory list of prefixes to use (default "mirrorlist.txt")
  -output string
        Path to put the repo files (default ".")
  -shuffle int
        Shuffle the mirror list every N downloads (default 100)
  -test
        Just validate downloaded files
  -threads int
        Concurrent downloads (default 4)
  -timeout duration
        Max connection time, in case a mirror slows significantly (default 10m0s)
```
