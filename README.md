# torrentcleaner
torrentcleaner is a simple tool that syncs a torrent content directory with the
files listed in the ".torrent" file. 

It lists all the files that are on disk but not in the torrent and lets you
easily delete them. 

## Building torrentcleaner
```bash
cd torrentcleaner
go build
```

## Running torrentcleaner
```
torrentcleaner <torrent-file> <content-directory>



