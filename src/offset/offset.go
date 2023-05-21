package offset

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Offset struct {
	FileName string
	Offset   int64
}

func ReadOffsetFile(offsetFile string) (map[string]*Offset, error) {

	var offsetMap map[string]*Offset = make(map[string]*Offset)

	if _, err := os.Stat(offsetFile); err != nil { // file does not exists ... return empty map
		return offsetMap, nil
	}

	f, err := os.Open(offsetFile)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		offsetLine := fileScanner.Text()

		// <file-name> <offset>
		split := strings.Split(offsetLine, " ")

		if len(split) != 2 {

			err = errors.New("invalid offset line " + offsetLine + ": too many items")

			return nil, err
		}

		fileName := split[0]

		offsetValue, err := strconv.ParseInt(split[1], 10, 64)

		if err != nil {
			err = errors.New("invalid offset line " + offsetLine + ": error in offset conversion")

			return nil, err
		}

		offset := Offset{fileName, offsetValue}

		offsetMap[fileName] = &offset
	}

	return offsetMap, nil
}

func WriteOffsetFile(offsetFile string, offsetMap map[string]*Offset) error {

	f, err := os.OpenFile(offsetFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)

	if err != nil {
		return err
	}

	defer f.Close()

	for key, element := range offsetMap {

		f.WriteString(fmt.Sprintf("%s %v\n", key, element.Offset))
	}

	return nil
}

func ClearOffsetMap(offsetMap map[string]*Offset) map[string]*Offset {

	newOffsetMap := make(map[string]*Offset)

	for fileName, _ := range offsetMap {

		if _, err := os.Stat(fileName); err == nil { // file exists ... put it to the map

			newOffsetMap[fileName] = offsetMap[fileName]
		}
	}

	return newOffsetMap
}
