package torrent

import (
	"io/ioutil"
	"log"
	"strconv"
)

type Entity interface{}

type Dict map[string]Entity

type List []Entity

type File struct {
	content []byte
	pos     int

	rootDict Dict
}

func (tf *File) read() byte {
	res := tf.content[tf.pos]
	tf.pos++
	return res
}

func (tf *File) unread() {
	tf.pos--
}

func (tf *File) readInt() int {
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

func (tf *File) readString() string {
	var data []byte
	var c byte
	for c = tf.read(); c != ':'; c = tf.read() {
		data = append(data, c)
	}
	l, _ := strconv.Atoi(string(data))
	data = make([]byte, l)
	for i := 0; i < l; i++ {
		data[i] = tf.read()
	}
	return string(data)
}

func (tf *File) readList() List {
	if tf.read() != 'l' {
		log.Fatal("list needs to start with l")
	}
	var content List
	for {
		content = append(content, tf.readEntity())
		if tf.read() == 'e' {
			break
		}
		tf.unread()
	}
	return content
}

func (tf *File) readDict() Dict {
	if tf.read() != 'd' {
		log.Fatal("dict needs to start with d")
	}
	content := make(Dict)
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

func (tf *File) readEntity() Entity {
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

func (tf *File) Files() []string {
	var res []string

	info := tf.rootDict["info"].(Dict)
	files := info["files"].(List)
	for idx, _ := range files {
		f := files[idx].(Dict)
		parts := f["path"].(List)
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

func (tf *File) Name() string {
	info := tf.rootDict["info"].(Dict)
	return info["name"].(string)
}

func (tf *File) Root() Dict {
	return tf.rootDict
}

func Load(path string) (*File, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tf := File{content: content}
	tf.rootDict = tf.readEntity().(Dict)
	return &tf, nil
}

