package generator

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
)

func Generate(path string, linesCount int, maxWordsPerLine int) (err error) {
	var file *os.File

	fmt.Printf("Creating file %s\n", path)

	if file, err = os.Create(path); err != nil {
		return
	}

	defer file.Close()

	loremIpsumGenerator := lorem.NewGenerator()

	fmt.Printf("Generating %d lines\n", linesCount)
	startTime := time.Now()

	linesCountTenth := linesCount / 10

	for index := 0; index < linesCount; index += 1 {
		if index%linesCountTenth == 0 {
			completePercent := int((float64(index) / float64(linesCount)) * 100)

			if completePercent > 0 {
				fmt.Printf("    %d%% complete in %v\n", completePercent, time.Now().Sub(startTime))
			}
		}

		wordsCount := rand.IntN(maxWordsPerLine-1) + 1
		loremIpsumGenerator.WordsPerSentence = wordsCount

		line := loremIpsumGenerator.Generate(wordsCount)

		file.Write([]byte(line))

		if index < linesCount-1 {
			file.Write([]byte{'\n'})
		}
	}

	duration := time.Now().Sub(startTime)

	fmt.Printf("\nLines generated in %v\n", duration)

	return
}
