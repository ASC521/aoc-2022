package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
--- Day 7: No Space Left On Device ---

You can hear birds chirping and raindrops hitting leaves as the expedition proceeds. Occasionally, you can even hear much louder sounds in the distance; how big do the animals get out here, anyway?

The device the Elves gave you has problems with more than just its communication system. You try to run a system update:

$ system-update --please --pretty-please-with-sugar-on-top
Error: No space left on device

Perhaps you can delete some files to make space for the update?

You browse around the filesystem to assess the situation and save the resulting terminal output (your puzzle input). For example:

$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k

The filesystem consists of a tree of files (plain data) and directories (which can contain other directories or files). The outermost directory is called /. You can navigate around the filesystem, moving into or out of directories and listing the contents of the directory you're currently in.

Within the terminal output, lines that begin with $ are commands you executed, very much like some modern computers:

    cd means change directory. This changes which directory is the current directory, but the specific result depends on the argument:
        cd x moves in one level: it looks in the current directory for the directory named x and makes it the current directory.
        cd .. moves out one level: it finds the directory that contains the current directory, then makes that directory the current directory.
        cd / switches the current directory to the outermost directory, /.
    ls means list. It prints out all of the files and directories immediately contained by the current directory:
        123 abc means that the current directory contains a file named abc with size 123.
        dir xyz means that the current directory contains a directory named xyz.

Given the commands and output in the example above, you can determine that the filesystem looks visually like this:

- / (dir)
  - a (dir)
    - e (dir)
      - i (file, size=584)
    - f (file, size=29116)
    - g (file, size=2557)
    - h.lst (file, size=62596)
  - b.txt (file, size=14848514)
  - c.dat (file, size=8504156)
  - d (dir)
    - j (file, size=4060174)
    - d.log (file, size=8033020)
    - d.ext (file, size=5626152)
    - k (file, size=7214296)

Here, there are four directories: / (the outermost directory), a and d (which are in /), and e (which is in a). These directories also contain files of various sizes.

Since the disk is full, your first step should probably be to find directories that are good candidates for deletion. To do this, you need to determine the total size of each directory. The total size of a directory is the sum of the sizes of the files it contains, directly or indirectly. (Directories themselves do not count as having any intrinsic size.)

The total sizes of the directories above can be found as follows:

    The total size of directory e is 584 because it contains a single file i of size 584 and no other directories.
    The directory a has total size 94853 because it contains files f (size 29116), g (size 2557), and h.lst (size 62596), plus file i indirectly (a contains e which contains i).
    Directory d has total size 24933642.
    As the outermost directory, / contains every file. Its total size is 48381165, the sum of the size of every file.

To begin, find all of the directories with a total size of at most 100000, then calculate the sum of their total sizes. In the example above, these directories are a and e; the sum of their total sizes is 95437 (94853 + 584). (As in this example, this process can count files more than once!)

Find all of the directories with a total size of at most 100000. What is the sum of the total sizes of those directories?

*/

var (
	CDCmdRX = regexp.MustCompile("[$]{1} cd (?P<arg>.*)")
	LSCmdRX = regexp.MustCompile("[$]{1} ls")
	DirRX   = regexp.MustCompile("dir (?P<name>[a-zA-Z]*)")
	FileRX  = regexp.MustCompile("(?P<size>[0-9]*) (?P<name>[a-zA-Z]*.[a-zA-Z]*)")
)

type File struct {
	name string
	size int
}

func NewFile(name string, size int) *File {
	return &File{
		name: name,
		size: size,
	}
}

type Directory struct {
	parent         *Directory
	name           string
	files          map[string]*File
	subDirectories map[string]*Directory
	size           int
}

func NewDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		parent:         parent,
		name:           name,
		files:          map[string]*File{},
		subDirectories: map[string]*Directory{},
		size:           -1,
	}
}

func (d *Directory) AddSubDirectory(name string) {
	d.subDirectories[name] = NewDirectory(name, d)
}

func (d *Directory) AddFile(name string, size int) {
	d.files[name] = NewFile(name, size)
}

func (d *Directory) Show(depth int) {
	pad := strings.Repeat("  ", depth)
	fmt.Printf("%v - %v (dir, size=%v)\n", pad, d.name, d.Size())

	for _, s := range d.subDirectories {
		s.Show(depth + 1)
	}

	for _, f := range d.files {
		pad = strings.Repeat("  ", depth+1)
		fileInfo := "file, size=" + strconv.Itoa(f.size)
		fmt.Printf("%v - %v (%v)\n", pad, f.name, fileInfo)
	}

}

func (d *Directory) Size() int {

	if d.size > -1 {
		return d.size
	}

	size := 0
	for _, f := range d.files {
		size += f.size
	}

	for _, s := range d.subDirectories {
		size += s.Size()
	}

	d.size = size

	return size

}

func (d *Directory) FilterMaxSize(maxSize int) []*Directory {

	dirs := []*Directory{}
	if d.Size() <= maxSize {
		dirs = append(dirs, d)
	}

	for _, s := range d.subDirectories {
		dirs = append(dirs, s.FilterMaxSize(maxSize)...)
	}

	return dirs

}

func (d *Directory) FilterMinSize(minSize int) []*Directory {

	dirs := []*Directory{}
	if d.Size() >= minSize {
		dirs = append(dirs, d)
	}

	for _, s := range d.subDirectories {
		dirs = append(dirs, s.FilterMinSize(minSize)...)
	}

	return dirs
}

type FileSystem struct {
	root *Directory
}

func NewFileSytem() *FileSystem {
	return &FileSystem{root: &Directory{
		parent:         nil,
		name:           "/",
		files:          map[string]*File{},
		subDirectories: map[string]*Directory{},
	}}
}

func buildFileSystem(cmds io.Reader) *FileSystem {
	fileSystem := NewFileSytem()
	currentLocation := fileSystem.root
	r := bufio.NewReader(cmds)
	for {
		line, readErr := r.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			panic(readErr)
		}

		switch {
		case CDCmdRX.Match(line):

			matches := CDCmdRX.FindStringSubmatch(string(line))
			arg := matches[1]

			if arg == ".." {
				currentLocation = currentLocation.parent
			} else if arg != "/" {
				currentLocation = currentLocation.subDirectories[arg]
			}

		case LSCmdRX.Match(line):

		case DirRX.Match(line):
			matches := DirRX.FindStringSubmatch(string(line))
			currentLocation.AddSubDirectory(string(matches[1]))
		case FileRX.Match(line):
			matches := FileRX.FindStringSubmatch(string(line))
			name := matches[2]
			size, err := strconv.Atoi(matches[1])
			if err != nil {
				panic(err)
			}
			currentLocation.AddFile(name, size)
		}

		if readErr == io.EOF {
			break
		}
	}

	return fileSystem
}

func main() {

	file, err := os.Open("/Users/alex.curto/code/aoc-2022/day7/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileSystem := buildFileSystem(file)
	fileSystem.root.Show(0)
	totalSize := 0
	for _, dir := range fileSystem.root.FilterMaxSize(100000) {
		totalSize += dir.Size()
	}
	fmt.Println(totalSize)

	totalAvailableSpace := 70000000
	freeSpace := totalAvailableSpace - fileSystem.root.Size()
	updateSize := 30000000
	spaceNeeded := updateSize - freeSpace
	possibleDirs := fileSystem.root.FilterMinSize(spaceNeeded)

	var smallest *Directory
	for _, d := range possibleDirs {
		if smallest == nil {
			smallest = d
			continue
		}

		if smallest.Size() > d.Size() {
			smallest = d
		}
	}

	fmt.Printf("Name: %v Size: %v\n", smallest.name, smallest.Size())

}
