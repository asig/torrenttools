/*
 * torrentlister: List all files in a torrent.
 *
 * Copyright (C) 2020 Andreas Signer <asigner@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/asig/torrenttools/pkg/torrent"
)

func usage() {
	fmt.Println("Usage: torrentlister <torrentfile>")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	torrentFile, err := torrent.Load(os.Args[1]);
	if err != nil {
		log.Fatal(err)
	}
	r := torrentFile.Name()
	var files []string
	for _, f := range torrentFile.Files() {
		files = append(files, filepath.Join(r, f))
	}
	sort.Strings(files)
	for _, f := range files {
		fmt.Printf("%s\n", f)
	}
}
