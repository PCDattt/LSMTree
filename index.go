package main

import (
	"bufio"
	"io"
	"os"
)

type Metadata struct {
	FileName string
	Offset int64
}

func NewMetadata(fileName string, offset int64) Metadata {
	return Metadata{
		FileName: fileName,
		Offset: offset,
	}
}

type Index struct {
	MapIndex map[int]Metadata // Map key to node metadata on disk
}

func NewIndex() Index {
	return Index{
		MapIndex: make(map[int]Metadata, 0),
	}
}

func (i *Index) Add(k int, fileName string, offset int64) {
	m := NewMetadata(fileName, offset)
	i.MapIndex[k] = m
}

func (i *Index) ConstructFromFile(f *os.File) error {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	var offset int64
	for scanner.Scan() {
		line := scanner.Text()
		n, err := NodeUnmarshal(line)
		if err != nil {
			return err
		}
		metadata := NewMetadata(f.Name(), offset)
		i.MapIndex[n.Key] = metadata
		offset += int64(len(line) + 1) // Plus 1 because of the end of line \n
	}
	return scanner.Err()
}
