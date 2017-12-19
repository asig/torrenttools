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
	"sort"
	"strconv"
)

const (
	ansiFgRed   = "31"
	ansiFgGreen = "32"
)

type TorrentFile struct {
	content []byte
	pos     int

	rootDict map[string]interface{}
}

var (
	torrentFile TorrentFile
	contentDir  string
)

func (tf *TorrentFile) read() byte {
	res := tf.content[tf.pos]
	tf.pos++
	return res
}

func (tf *TorrentFile) unread() {
	tf.pos--
}

func (tf *TorrentFile) readInt() int {
	if tf.read() != 'i' {
		log.Fatal("int needs to start with i")
	}
	var data []byte
	var c byte
	for c = tf.read(); c != 'e'; c = tf.read() {
		data = append(data, c)
	}
	res, _ := strconv.Atoi(string(data))
	return res
}

func (tf *TorrentFile) readString() string {
	var data []byte
	var c byte
	for c = tf.read(); c != ':'; c = tf.read() {
		data = append(data, c)
	}
	len, _ := strconv.Atoi(string(data))

	data = make([]byte, len)
	for i := 0; i < len; i++ {
		data[i] = tf.read()
	}

	return string(data)
}

func (tf *TorrentFile) readList() []interface{} {
	if tf.read() != 'l' {
		log.Fatal("list needs to start with l")
	}
	var content []interface{}
	for {
		content = append(content, tf.readEntity())
		if tf.read() == 'e' {
			break
		}
		tf.unread()
	}
	return content
}

func (tf *TorrentFile) readDict() map[string]interface{} {
	if tf.read() != 'd' {
		log.Fatal("dict needs to start with d")
	}
	content := make(map[string]interface{})
	for {
		key := tf.readString()
		val := tf.readEntity()
		content[key] = val
		if tf.read() == 'e' {
			break
		}
		tf.unread()
	}
	return content
}

func (tf *TorrentFile) readEntity() interface{} {
	c := tf.read()
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tf.unread()
		return tf.readString()
	case 'i':
		tf.unread()
		return tf.readInt()
	case 'l':
		tf.unread()
		return tf.readList()
	case 'd':
		tf.unread()
		return tf.readDict()
	default:
		log.Fatalf("Invalid torrent file. Unexecpected character %c", c)
	}
	return ""
}

func (tf *TorrentFile) Load() {
	tf.rootDict = tf.readEntity().(map[string]interface{})
}

func (tf *TorrentFile) Files() []string {
	var res []string

	info := tf.rootDict["info"].(map[string]interface{})
	files := info["files"].([]interface{})
	for idx, _ := range files {
		f := files[idx].(map[string]interface{})
		parts := f["path"].([]interface{})
		path := ""
		for _, p := range parts {
			if len(path) > 0 {
				path = path + "/"
			}
			path = path + p.(string)
		}
		res = append(res, path)
	}
	return res
}

func usage() {
	fmt.Println("Usage: torrentcleaner <torrentfile> <contentdir>")
	os.Exit(1)
}

func init() {
	if len(os.Args) != 3 {
		usage()
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	torrentFile = TorrentFile{content: content}
	torrentFile.Load()

	contentDir = os.Args[2]
	for contentDir[len(contentDir)-1] == '/' {
		contentDir = contentDir[0 : len(contentDir)-1]
	}
}

func readFiles(dir string, m map[string]bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			readFiles(dir+"/"+f.Name(), m)
		} else {
			m[dir+"/"+f.Name()] = true
		}
	}
}

func askDelete() bool {
	delete := false
	for {
		fmt.Print("Delete? (Y/n) ")
		var input string
		fmt.Scanln(&input)
		if input == "y" || input == "Y" || input == "" {
			delete = true
			break
		} else if input == "n" || input == "N" {
			delete = false
			break
		} else {
			fmt.Println("I'm sorry, but I don't understand your choice.")
		}
	}
	return delete
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
	filesInTorrent := make(map[string]bool)
	for _, f := range torrentFile.Files() {
		filesInTorrent[contentDir+"/"+f] = true
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
