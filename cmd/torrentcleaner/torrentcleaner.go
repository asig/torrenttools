/*
 * torrentcleaner: Delete all files not listed in a torrent.
 *
 * Copyright (C) 2017 Andreas Signer <asigner@gmail.com>
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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/asig/torrenttools/pkg/torrent"
)

const (
	ansiFgRed   = "31"
	ansiFgGreen = "32"
)

func usage() {
	fmt.Println("Usage: torrentcleaner <torrentfile> <contentdir>")
	os.Exit(1)
}

func readFiles(dir string, m map[string]bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		p := filepath.Clean(filepath.Join(dir, f.Name()))
		if f.IsDir() {
			readFiles(p, m)
		} else {
			m[p] = true
		}
	}
}

func askDelete() bool {
	res := false
	for {
		fmt.Print("Delete? (Y/n) ")
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			if input == "y" || input == "Y" || input == "" {
				res = true
				break
			} else if input == "n" || input == "N" {
				res = false
				break
			} else {
				fmt.Println("I'm sorry, but I don't understand your choice.")
			}
		}
	}
	return res
}

func listFiles(files []string) {
	for _, s := range files {
		fmt.Printf("  %s\n", s)
	}
}

func colored(s string, col string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", col, s)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	torrentFile, err := torrent.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	contentDir := os.Args[2]
	for contentDir[len(contentDir)-1] == '/' {
		contentDir = contentDir[0 : len(contentDir)-1]
	}

	filesInTorrent := make(map[string]bool)
	r := torrentFile.Name()
	for _, f := range torrentFile.Files() {
		filesInTorrent[filepath.Join(contentDir, r, f)] = true
	}

	filesInDirectory := make(map[string]bool)
	readFiles(contentDir, filesInDirectory)

	var onlyInTorrent []string
	var onlyInDirectory []string

	for key, _ := range filesInTorrent {
		if _, ok := filesInDirectory[key]; !ok {
			onlyInTorrent = append(onlyInTorrent, key)
		}
	}
	sort.Strings(onlyInTorrent)

	for key, _ := range filesInDirectory {
		if _, ok := filesInTorrent[key]; !ok {
			onlyInDirectory = append(onlyInDirectory, key)
		}
	}
	sort.Strings(onlyInDirectory)

	if len(onlyInTorrent) > 0 {
		fmt.Println("The following files are only in the torrent, but not in the directory:")
		listFiles(onlyInTorrent)
		fmt.Printf("----------------------------------------------\n")
	}
	if len(onlyInDirectory) > 0 {
		fmt.Println("The following files are only in the directory, but not in the torrent:")
		listFiles(onlyInDirectory)
		if askDelete() {
			for _, s := range onlyInDirectory {
				var res string
				if os.Remove(s) != nil {
					res = colored("FAILED", ansiFgRed)
				} else {
					res = colored("ok", ansiFgGreen)
				}
				fmt.Printf("Deleting %s: %s\n", s, res)
			}
		}
	}
}
