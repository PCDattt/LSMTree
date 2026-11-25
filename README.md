# LSM Tree

A minimal, educational implementation of a Log-Structured Merge Tree (LSM Tree) in Go.

The goal of this project is to understand the fundamental concepts behind LSM Trees, including in-memory data structures, disk-based storage, indexing, flushing, and compaction.

## Implementation

- **MemTable:** Implemented using a Binary Search Tree (BST) to store key-value pairs in memory.
- **Index:** Uses a Go map to map keys to their corresponding file names and disk offsets.
- **Storage:** Persists data to disk using files.
- **Compaction:** Merges sorted data from two files using a merge-sort approach.
- **Unit Testing:** 75%

## Features

- **Write:** Add key-value pairs to the MemTable. Time complexity: O(h), where h is the height of the BST.
- **Read:** Search for a key in the MemTable first. If the key is not found, look it up in the in-memory index and use the corresponding file offset to read the value from disk.
  - MemTable lookup: O(h).
  - Index lookup: O(1) average-case complexity, excluding disk I/O.
- **Flush:** Write the MemTable's contents to disk in sorted order and update the in-memory index. Time complexity: O(n).
- **Compact:** Merge the sorted contents of the current and compacted files, keeping the latest value for duplicate keys. Time complexity: O(n + m), where n and m are the number of entries in each file.
- **Index Construction:** Rebuild the in-memory index by scanning the existing data files. Time complexity: O(n), where n is the total number of entries.

## Future Improvements

- Replace text-based storage (`.txt`) with binary files format (`.bin`) for more efficient storage and serialization.
- Implement a **Sparse Index** to reduce memory consumption by avoiding the need to keep an index entry for every key.
- Support efficient range queries using the Sparse Index.
- Implement crash-safe writes and recovery to prevent data loss after unexpected shutdowns.
- Ensure thread-safe concurrent reads and writes.
- Replace the unbalanced BST with a self-balancing tree to guarantee O(log n) insertion and lookup complexity.
