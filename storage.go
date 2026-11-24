package main

import (
	"bufio"
	"io"
	"os"
	"errors"
)

const CONST_COMPACTED_FILE_NAME = "compactedFile.txt"
const CONST_CURRENT_FILE_NAME = "currentFile.txt"
const CONST_NEW_COMPACTED_FILE_NAME = "newCompactedFile.txt"

type Storage struct {
	MemTable BinaryTree
	Index Index 
	currentFile *os.File
	compactedFile *os.File
}

func NewStorage(currentFileName, compactedFileName string) (*Storage, error) {
	currentFile, err := os.OpenFile(currentFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	compactedFile, err := os.OpenFile(compactedFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	i := NewIndex()
	// Construct current file later so can keep recent duplicate key
	i.ConstructFromFile(compactedFile)
	i.ConstructFromFile(currentFile)
	s := &Storage{
		MemTable: NewBinaryTree(),
		Index: i,
		currentFile: currentFile,
		compactedFile: compactedFile,
	}
	return s, nil
}

func (s *Storage) Close() {
	// Flush before close to persist data to disk
	s.Flush()
	s.compactedFile.Close()
	s.currentFile.Close()
}

func (s *Storage) Add(k int, v string) {
	n := NewNode(k, v)
	s.MemTable.AddNode(n)
}

func (s *Storage) Find(k int) (string,error) {
	n := s.MemTable.FindNode(s.MemTable.Head, k)
	if n != nil {
		res := NodeMarshal(*n)
		return res, nil
	}
	metadata := s.Index.MapIndex[k]
	if metadata.FileName == "" {
		return "", errors.New("Key doesn't exist")
	}
	if metadata.FileName == s.currentFile.Name() {
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
	} else if metadata.FileName == s.compactedFile.Name() {
		_, err := s.compactedFile.Seek(metadata.Offset, io.SeekStart)
		if err != nil {
			return "", err
		}
		reader := bufio.NewReader(s.compactedFile)
		res , err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return res, nil
	}
	return "", errors.New("Key doesn't exist")
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
		s.Index.Add(n.Key, s.currentFile.Name(), offset)
		if err != nil {
			return err
		}
	}
	s.MemTable.Clear()
	return nil
}

func (s *Storage) Compact() error {
	newCompactFile, err := os.OpenFile(CONST_NEW_COMPACTED_FILE_NAME,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	compactedFileScanner:= bufio.NewScanner(s.compactedFile)
	currentFileScanner := bufio.NewScanner(s.currentFile)
	compactedLine, currentLine := "", ""
	nodeCompacted, nodeCurrent := Node{}, Node{}
	canScanCompactedFile, canScanCurrentFile := true , true

	// Merge sort 2 files
	for {
		if compactedLine == "" {
			canScanCompactedFile = compactedFileScanner.Scan()
			if !canScanCompactedFile {
				break
			}
			compactedLine = compactedFileScanner.Text()
			nodeCompacted, err = NodeUnmarshal(compactedLine)
		}
		if currentLine == "" {
			canScanCurrentFile = currentFileScanner.Scan()
			if !canScanCurrentFile {
				break
			}
			currentLine = currentFileScanner.Text()
			nodeCurrent, err = NodeUnmarshal(currentLine)
		}
		if nodeCompacted.Key < nodeCurrent.Key {
			offset, err := newCompactFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			newCompactFile.WriteString(NodeMarshal(nodeCompacted) + "\n")
			s.Index.Add(nodeCompacted.Key, s.compactedFile.Name(), offset)
			compactedLine = ""
		} else {
			offset, err := newCompactFile.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}
			newCompactFile.WriteString(NodeMarshal(nodeCurrent)+ "\n")
			s.Index.Add(nodeCurrent.Key, s.compactedFile.Name(), offset)
			currentLine = ""
			// If equal then pick only node from current file
			// and skip node in compacted file
			if nodeCompacted.Key == nodeCurrent.Key {
				compactedLine = ""
			}
		}
	}
	for {
		if compactedLine == "" {
			canScanCompactedFile = compactedFileScanner.Scan()
			if !canScanCompactedFile {
				break
			}
			compactedLine = compactedFileScanner.Text()
			nodeCompacted, err = NodeUnmarshal(compactedLine)
		}
		offset, err := newCompactFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		newCompactFile.WriteString(NodeMarshal(nodeCompacted)+ "\n") 
		s.Index.Add(nodeCompacted.Key, s.compactedFile.Name(), offset)
		compactedLine = ""
	}
	for {
		if currentLine == "" {
			canScanCurrentFile = currentFileScanner.Scan()
			if !canScanCurrentFile {
				break
			}
			currentLine = currentFileScanner.Text()
			nodeCurrent, err = NodeUnmarshal(currentLine)
		}
		offset, err := newCompactFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		newCompactFile.WriteString(NodeMarshal(nodeCurrent) + "\n")
		s.Index.Add(nodeCurrent.Key, s.compactedFile.Name(), offset)
		currentLine = ""
	}
	// end of Merge sort

	newCompactFile.Close()
	s.compactedFile.Close()
	err = os.Remove(s.compactedFile.Name())
	if err != nil {
		return err
	}
	err = os.Rename(CONST_NEW_COMPACTED_FILE_NAME, s.compactedFile.Name())
	if err != nil {
		return err
	}
	s.compactedFile, err = os.OpenFile(s.compactedFile.Name(),os.O_RDWR|os.O_CREATE, 0644)
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
