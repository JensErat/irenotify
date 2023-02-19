# irenotify

[inotify](https://en.wikipedia.org/wiki/Inotify) is a Linux kernel subsystem to efficiently monitor changes to the filesystem. It allows applications to monitor changes to files, but requires special support for network file system like NFS, SMB and CephFS that are often missing.

This tool, `irenotify`, monitors a directory to observe file changes and touches changed files to enforce inotify events to be sent out. It recognizes changed files by periodically and recursively listing files in a directory and comparing modify timestamps.

## Limitations

- `irenotify` needs to poll for changes, suffering a performance penalty for large number of files and putting load on a network filesystem
- changes to a file without modification of the last modified timestamp will not be observed
- changes to a file's last modified timestamp without actual changes to the file will trigger a touch event anyway
- must not be used with multiple instances on the same directory -- the last modified timestamp will be synced and other instances will change the timestamp again, resulting in a loop (could be resolved by remembering a checksum, but not implemented yet)

## Usage

```angular2html
Usage of irenotify:
  -delay duration
        delay between subsequent filesystem scans (default 1s)
  -path string
        path to monitor, watching working directory if not set
  -v int
        log level (default 2)
```