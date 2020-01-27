# Torrent tools

This project contains three small utilities to help you manage torrents.

## torrentcleaner
`torrentcleaner` is a simple tool that syncs a torrent content directory with the
files listed in the ".torrent" file. 

It lists all the files that are on disk but not in the torrent and lets you
easily delete them. 

## torrentlister
`torrentlister` lists all files that are contained in the ".torrent" file. 

## torrentdumper
`torrentdumper` dumps all the information in the ".torrent" file in a human-readable form.

## Building the tools
```bash
cd torrenttools
go build ./cmd/torrentcleaner 
go build ./cmd/torrentlister 
go build ./cmd/torrentdumper
```

## Running torrentcleaner
```
torrentcleaner <torrent-file> <content-directory>
```

## Running torrentlister
```
torrentlister <torrent-file>
```

## Running torrentdumper
```
torrentdumper <torrent-file>
```
