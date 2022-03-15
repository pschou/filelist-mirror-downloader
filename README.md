# Filelist Mirror Downloader

This shim takes a file list and a mirror list and downloads the files in the file list
while doing checksums and parallel downloads for speed.

# Example usage:
```bash
./filelist-mirror-downloader -output output
```


# Usage help:
```bash
$ ./filelist-mirror-downloader -h
Yum Get RepoMD,  Version: 0.1.20220315.1310

Usage: ./filelist-mirror-downloader [options...]

  -debug
        Turn on debug comments
  -list string
        Filelist to be fetched (one per line with: HASH SIZE PATH) (default "filelist.txt")
  -mirrors string
        Mirror / directory list of prefixes to use (default "mirrorlist.txt")
  -output string
        Path to put the repo files (default ".")
  -threads int
        Concurrent downloads (default 1)

```




