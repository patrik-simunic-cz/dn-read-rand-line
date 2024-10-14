package reader

import (
	"fmt"
	"io"
	"os"
	"time"
)

const MAX_BYTES_PER_LINE = 1000

func ReadLine(
	indexPath string,
	path string,
	outputLine int,
	enableVerbose bool,
) (err error) {
	var (
		fileInfo os.FileInfo
		file     *os.File
		index    *Index
		output   string
	)

	if file, err = os.Open(path); err != nil {
		return
	}

	if fileInfo, err = file.Stat(); err != nil {
		return
	}

	defer file.Close()

	if index, err = openIndex(indexPath); err != nil {
		return
	}

	if !index.isIndexed(path, fileInfo.ModTime()) {
		if enableVerbose {
			fmt.Printf("Index not found. Indexing file...\n")
		}

		indexingStart := time.Now()

		if index.indexFile(path, file, fileInfo); err != nil {
			return
		}

		if err = index.write(); err != nil {
			return
		}

		if enableVerbose {
			fmt.Printf("File has been indexed in %v\n", time.Now().Sub(indexingStart))
		}
	} else if enableVerbose {
		fmt.Printf("Index loaded\n")
	}

	fileIndex, hasFileIndex := index.Files[path]
	if !hasFileIndex {
		panic("File index not found")
	}

	if outputLine > fileIndex.LinesTotal {
		return fmt.Errorf("Cannot read line %d, the source file has only %d\n", outputLine, fileIndex.LinesTotal)
	}

	if enableVerbose {
		fmt.Printf("Looking up line %d out of %d\n", outputLine, fileIndex.LinesTotal)
	}

	lookupStart := time.Now()

	lineOffsetIndex := 0
	for lineOffsetIndex < len(fileIndex.LineOffsets) {
		if fileIndex.LineOffsets[lineOffsetIndex][0] < outputLine {
			if lineOffsetIndex+1 >= len(fileIndex.LineOffsets) || fileIndex.LineOffsets[lineOffsetIndex+1][0] > outputLine {
				break
			}

			lineOffsetIndex += 1
			continue
		}

		break
	}

	startBufferPosition := fileIndex.LineOffsets[lineOffsetIndex][1]

	buffer := make([]byte, INDEX_LINES_OFFSET_SIZE*(MAX_BYTES_PER_LINE+1))
	bufferLine := fileIndex.LineOffsets[lineOffsetIndex][0]

	if _, err = file.Seek(int64(startBufferPosition), 0); err != nil {
		return
	}

	if _, err = file.Read(buffer); err != nil {
		if err != io.EOF {
			return
		}

		err = nil
	}

	for _, char := range buffer {
		if bufferLine >= outputLine {
			if char == '\n' {
				break
			}

			output += string(char)
		}

		if char == '\n' {
			bufferLine += 1
		}
	}

	if enableVerbose {
		fmt.Printf("Lookup finished in %v\n", time.Now().Sub(lookupStart))
	}

	if enableVerbose {
		fmt.Printf("Looked up line:\n---\n%s\n---\n", output)
	} else {
		fmt.Printf("%s\n", output)
	}

	return
}
