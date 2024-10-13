package reader

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"
)

const (
	INDEX_LINES_OFFSET_SIZE = 1000
	READ_BUFFER_SIZE        = 1024 * 8
)

type Index struct {
	Path  string
	Files map[string]*FileIndex
}

type FileIndex struct {
	Path         string
	LastModified time.Time
	LinesTotal   int
	LineOffsets  [][]int
}

func openIndex(indexPath string) (index *Index, err error) {
	var file *os.File

	index = &Index{
		Path:  indexPath,
		Files: map[string]*FileIndex{},
	}

	if file, err = os.Open(indexPath); err != nil {
		if !os.IsNotExist(err) {
			return
		}

		err = nil
		return
	}

	defer file.Close()

	err = index.read(file)
	return
}

func (index *Index) read(indexFile *os.File) (err error) {
	var data []byte

	if data, err = ioutil.ReadAll(indexFile); err != nil {
		return
	}

	dataLength := len(data)

	if dataLength < 1 {
		return
	}

	offset := 0

	for offset < dataLength && data[offset] == byte(28) {
		fileIndex := &FileIndex{}

		offset += 1

		filePathLength := int(data[offset])

		for itemOffset := 1; itemOffset <= filePathLength; itemOffset += 1 {
			fileIndex.Path += string(data[offset+itemOffset])
		}

		offset += filePathLength + 1

		lastModified := ""
		lastModifiedLength := int(data[offset])

		for itemOffset := 1; itemOffset <= lastModifiedLength; itemOffset += 1 {
			lastModified += string(data[offset+itemOffset])
		}

		if fileIndex.LastModified, err = time.Parse(time.RFC822, lastModified); err != nil {
			err = fmt.Errorf("Reading index failed: \"lastModified\" is not a time (%s)", err)
			return
		}

		offset += lastModifiedLength + 1

		linesTotal := ""
		linesTotalLength := int(data[offset])

		for itemOffset := 1; itemOffset <= linesTotalLength; itemOffset += 1 {
			linesTotal += string(data[offset+itemOffset])
		}

		if fileIndex.LinesTotal, err = strconv.Atoi(linesTotal); err != nil {
			err = fmt.Errorf("Reading index failed: \"linesTotal\" is not a number (%s)", err)
			return
		}

		offset += linesTotalLength + 1

		hasMore := true
		for hasMore {
			var (
				line     int
				position int
			)

			lineRaw := ""
			for {
				lineRaw += string(data[offset])
				offset += 1

				if data[offset] == byte(29) {
					offset += 1
					break
				}
			}

			positionRaw := ""
			for {
				positionRaw += string(data[offset])
				offset += 1

				if data[offset] == byte(30) {
					offset += 1
					break
				}
			}

			if line, err = strconv.Atoi(lineRaw); err != nil {
				err = fmt.Errorf("Reading index failed: \"line\" is not a number (%s)", err)
			}

			if position, err = strconv.Atoi(positionRaw); err != nil {
				err = fmt.Errorf("Reading index failed: \"position\" is not a number (%s)", err)
			}

			fileIndex.LineOffsets = append(fileIndex.LineOffsets, []int{line, position})

			if offset >= dataLength-1 || data[offset] == byte(28) {
				hasMore = false
			}
		}

		index.Files[fileIndex.Path] = fileIndex
	}

	return
}

func (index *Index) write() (err error) {
	var file *os.File

	if file, err = os.Open(index.Path); err != nil {
		if !os.IsNotExist(err) {
			return
		}

		err = nil

		if file, err = os.Create(index.Path); err != nil {
			return
		}
	}

	defer file.Close()

	for filePath, fileIndex := range index.Files {
		file.Write([]byte{28, byte(len(filePath))})
		file.Write([]byte(filePath))

		lastModified := fileIndex.LastModified.Format(time.RFC822)

		file.Write([]byte{byte(len(lastModified))})
		file.Write([]byte(lastModified))

		linesTotal := fmt.Sprintf("%d", fileIndex.LinesTotal)

		file.Write([]byte{byte(len(linesTotal))})
		file.Write([]byte(linesTotal))

		for _, offset := range fileIndex.LineOffsets {
			line := fmt.Sprintf("%d", offset[0])
			position := fmt.Sprintf("%d", offset[1])

			file.Write([]byte(line))
			file.Write([]byte{29})
			file.Write([]byte(position))
			file.Write([]byte{30})
		}
	}

	return
}

func (index *Index) isIndexed(path string, lastModified time.Time) bool {
	fileIndex, hasIndex := index.Files[path]

	if !hasIndex {
		return false
	}

	// TODO: Fix index invalidation
	// return fileIndex.LastModified == lastModified

	_ = fileIndex
	return true
}

func (index *Index) indexFile(path string, file *os.File, fileInfo os.FileInfo) (err error) {
	fileIndex := &FileIndex{
		Path:         path,
		LastModified: fileInfo.ModTime(),
		LinesTotal:   1,
		LineOffsets:  [][]int{{1, 0}},
	}

	fileSize := fileInfo.Size()

	position := 0
	for int64(position) < fileSize {
		buffer := make([]byte, int(math.Min(READ_BUFFER_SIZE, float64(fileSize-int64(position)))))

		_, err = file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return
			}

			err = nil
		}

		for _, char := range buffer {
			if char == '\n' {
				fileIndex.LinesTotal += 1

				if fileIndex.LinesTotal%INDEX_LINES_OFFSET_SIZE == 0 {
					fileIndex.LineOffsets = append(fileIndex.LineOffsets, []int{fileIndex.LinesTotal, position + 1})
				}
			}

			position += 1
		}
	}

	index.Files[path] = fileIndex
	return
}
