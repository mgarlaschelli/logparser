package stream

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mgarlaschelli/logparser/filter"
	"github.com/mgarlaschelli/logparser/offset"
)

func ParseFiles(filePath string, filters []filter.LineFilter, offsets map[string]*offset.Offset, processAllFiles bool, ignoreOld int) error {
	fmt.Println("<matching-lines>")
	defer fmt.Println("</matching-lines>")

	files, err := filepath.Glob(filePath)
	if err != nil {
		return err
	}

	if !processAllFiles {

		lastUpdatedFile := getLastUpdatedFile(files)

		info, err := os.Stat(lastUpdatedFile)
		if err != nil {
			return err
		}

		if info.ModTime().Before(time.Now().AddDate(0, 0, -ignoreOld)) {
			return nil
		}

		processFile(lastUpdatedFile, filters, offsets, ignoreOld)

		return nil
	}

	for _, file := range files {
		processFile(file, filters, offsets, ignoreOld)
	}

	return nil
}

func getLastUpdatedFile(files []string) string {

	var lastUpdateFile string
	var lastUpdateTime time.Time

	for _, file := range files {

		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if lastUpdateTime.IsZero() {
			lastUpdateTime = info.ModTime()
			lastUpdateFile = file
		} else if info.ModTime().After(lastUpdateTime) {
			lastUpdateTime = info.ModTime()
			lastUpdateFile = file
		}
	}

	return lastUpdateFile
}

func processFile(file string, filters []filter.LineFilter, offsets map[string]*offset.Offset, ignoreOld int) {

	info, err := os.Stat(file)
	if err != nil {
		return
	}

	if info.ModTime().Before(time.Now().AddDate(0, 0, -ignoreOld)) {
		return
	}

	fp, err := os.Open(file)
	if err != nil {
		return
	}
	defer fp.Close()

	offsetValue, offsetFound := offsets[file]

	if offsetFound {
		fp.Seek(offsetValue.Offset, 0)
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {

		line := scanner.Text()

		for _, filter := range filters {

			match, err := filter.Match(line)

			if err != nil {
				continue
			}

			if match {
				fmt.Println(line)
				break
			}
		}
	}

	newOffset, err := fp.Seek(0, io.SeekCurrent)

	if err != nil {
		return
	}

	if offsetFound {
		offsetValue.Offset = newOffset
	} else {
		offsets[file] = &offset.Offset{file, newOffset}
	}
}
