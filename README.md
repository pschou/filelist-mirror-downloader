# Filelist Mirror Downloader

This shim takes a file list and a mirror list and downloads the files in the
file list while doing checksums and parallel downloads for speed.

When packages are downloaded (specifically DEB, RPM and APK), the time stamps
of the file are derived from the contents of the file.  This way if a mirror
has had all the files touched to another time stamp, the correct timestamp is
restored.

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
Yum Get RepoMD

Usage: ./filelist-mirror-downloader [options...]

  -after string
        Select packages after specified date
        Date formats supported: https://github.com/araddon/dateparse
  -attempts int
        Attempts for each file (default 40)
  -background
        Ignore all keyboard inputs, background mode (for non-interactive scripting)
  -before string
        Select packages before specified date
        Date formats supported: https://github.com/araddon/dateparse
  -client-cert string
        Satellite repo, CERT for using PKI auth
  -client-key string
        Satellite repo, KEY for using PKI auth
  -client-pass string
        Satellite repo, PASS for USER
  -client-user string
        Satellite repo, using basic USER auth
  -debug
        Turn on debug comments
  -dup string
        What to do with duplicates: omit, copy, symlink, hardlink (default "symlink")
  -keyring string
        Use keyring for verifying signed package files (example: keyring.gpg or keys/ directory)
  -list string
        Filelist to be fetched (one per line with: HASH SIZE PATH) (default "filelist.txt")
  -log string
        File in which to store a log of files downloaded
        Line prefixes (OnDisk, OnDiskSkip, Skipped, Fetched, Uncompressed, Failed), indicate action taken.
        Skip means that a file falls outside the required date bounds
  -mirrors string
        Mirror / directory list of prefixes to use (default "mirrorlist.txt")
  -output string
        Path to put the repo files (default ".")
  -shuffle int
        Shuffle the mirror list every N downloads (default 100)
  -threads int
        Concurrent downloads (default 4)
  -timeout duration
        Max connection time per file, in case a mirror slows significantly
        If one is downloading large ISO files, a longer time may be needed. (default 10m0s)
```
