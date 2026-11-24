package main

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
	MapIndex map[int]Metadata // Map key to metadata on disk
}

func (i *Index) Add(k int, fileName string, offset int64) {
	m := NewMetadata(fileName, offset)
	i.MapIndex[k] = m
}
