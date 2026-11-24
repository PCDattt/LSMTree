package main

import (
	"bufio"
	"io"
	"os"
	"errors"
)

const compactedFileName = "compactedFile.txt"
const currentFileName = "currentFile.txt"
const newCompactedFileName = "newCompactedFile.txt"

type Storage struct {
	MemTable BinaryTree
	Index Index 
	currentFile *os.File
	compactedFile *os.File
}

func NewStorage() (*Storage, error) {
	err := os.Remove(currentFileName)
	if err != nil {
		return nil, err
	}
	err = os.Remove(compactedFileName)
	if err != nil {
		return nil, err
	}
	currentFile, err := os.OpenFile(currentFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	compactedFile, err := os.OpenFile(compactedFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	mapIndex := make(map[int]Metadata, 0)
	s := &Storage{
		MemTable: BinaryTree{},
		Index: Index{
			MapIndex: mapIndex,
		},
		currentFile: currentFile,
		compactedFile: compactedFile,
	}
	return s, nil
}

func (s *Storage) Add(k int, v string) {
	n := NewNode(k, v)
	s.MemTable.AddNode(n)
}

func (s *Storage) Find(k int) (string,error) {
	n := s.MemTable.FindNode(s.MemTable.Head, k)
	if n != nil {
		return n.Value, nil
	}
	metadata := s.Index.MapIndex[k]
	if metadata.FileName == "" {
		return "", errors.New("Key doesn't exist")
	}
	if metadata.FileName == currentFileName {
		_, err := s.currentFile.Seek(metadata.Offset, io.SeekStart)
		if err != nil {
			return "", err
		}
		reader := bufio.NewReader(s.currentFile)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return line, nil
	} else if metadata.FileName == compactedFileName {
		_, err := s.compactedFile.Seek(metadata.Offset, io.SeekStart)
		if err != nil {
			return "", err
		}
		reader := bufio.NewReader(s.compactedFile)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return line, nil
	}
	return "", nil
}

func (s *Storage) Flush() error {
	head := s.MemTable.Head
	listNode := make([]Node, 0)
	s.MemTable.InorderTraversal(head, &listNode)
	_, err := s.currentFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	for _, n := range listNode {
		offset, err := s.currentFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		str := NodeMarshal(n) + "\n"
		_, err = s.currentFile.WriteString(str)
		s.Index.Add(n.Key, currentFileName, offset)
		if err != nil {
			return err
		}
	}
	s.MemTable.Clear()
	return nil
}

func (s *Storage) Compact() error {
	newCompactFile, err := os.OpenFile(newCompactedFileName,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = s.compactedFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = s.currentFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	oldCompactedFileScanner:= bufio.NewScanner(s.compactedFile)
	currentFileScanner := bufio.NewScanner(s.currentFile)
	oldCompactedLine, currentLine := "", ""
	nodeCompacted, nodeCurrent := Node{}, Node{}

	// Merge Sort on 2 files - incomplete
	for oldCompactedFileScanner.Scan() && currentFileScanner.Scan() {
		if oldCompactedLine == "" {
			oldCompactedLine = oldCompactedFileScanner.Text()
			nodeCompacted, err = NodeUnmarshal(oldCompactedLine)
		}
		if currentLine == "" {
			currentLine = currentFileScanner.Text()
			nodeCurrent, err = NodeUnmarshal(currentLine)
		}
		if nodeCompacted.Key < nodeCurrent.Key {
			offset, err := newCompactFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			newCompactFile.WriteString(NodeMarshal(nodeCompacted) + "\n")
			s.Index.Add(nodeCompacted.Key, compactedFileName, offset)
			oldCompactedLine = ""
		} else {
			offset, err := newCompactFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			newCompactFile.WriteString(NodeMarshal(nodeCurrent)+ "\n")
			s.Index.Add(nodeCurrent.Key, compactedFileName, offset)
			currentLine = ""
		}
	}
	for oldCompactedFileScanner.Scan() {
		if oldCompactedLine == "" {
			oldCompactedLine = oldCompactedFileScanner.Text()
			nodeCompacted, err = NodeUnmarshal(oldCompactedLine)
		}
		offset, err := newCompactFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		newCompactFile.WriteString(NodeMarshal(nodeCompacted)+ "\n") 
		s.Index.Add(nodeCompacted.Key, compactedFileName, offset)
		oldCompactedLine = ""
	}
	for currentFileScanner.Scan() {
		if currentLine == "" {
			currentLine = currentFileScanner.Text()
			nodeCurrent, err = NodeUnmarshal(currentLine)
		}
		offset, err := newCompactFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
		newCompactFile.WriteString(NodeMarshal(nodeCurrent) + "\n")
			s.Index.Add(nodeCurrent.Key, compactedFileName, offset)
			currentLine = ""
	}
	// end of Merge sort

	newCompactFile.Close()
	s.compactedFile.Close()
	err = os.Remove(compactedFileName)
	if err != nil {
		return err
	}
	err = os.Rename(newCompactedFileName, compactedFileName)
	if err != nil {
		return err
	}
	s.compactedFile, err = os.OpenFile(compactedFileName,os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = s.currentFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = s.currentFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
