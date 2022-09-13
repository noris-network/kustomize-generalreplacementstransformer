package grt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

const maxLineLength = 5 * 2 << 19 // 5M

func (t *Transformer) ReadStream(r io.Reader) error {

	var dumpTo *os.File

	dumpDest := os.Getenv("GRT_DUMP_INPUT")
	if dumpDest != "" {
		file, err := os.Create(dumpDest)
		if err != nil {
			return err
		}
		defer file.Close()
		dumpTo = file
	}

	scanner := bufio.NewScanner(r)

	// as some secrets, e.g. certificates or keys are quite long, provide a custom buffer
	scannerBuffer := []byte{}
	scanner.Buffer(scannerBuffer, maxLineLength)

	emptyLineRE := regexp.MustCompile("^ *$")

	buf := bytes.Buffer{}
	objectIsEmpty := true
	for scanner.Scan() {
		line := scanner.Text()
		if dumpTo != nil {
			fmt.Fprintln(dumpTo, line)
		}
		if line == "---" {
			if buf.Len() >= 2 {
				if err := t.RegisterRaw(buf.Bytes()); err != nil {
					t.uus = nil
					return err
				}
			}
			buf.Reset()
			objectIsEmpty = true
		} else {
			if !emptyLineRE.MatchString(line) {
				objectIsEmpty = false
			}
			if !objectIsEmpty {
				buf.WriteString(line)
				buf.WriteString("\n")
			}
		}
	}

	if buf.Len() >= 2 {
		if err := t.RegisterRaw(buf.Bytes()); err != nil {
			t.uus = nil
			return err
		}
	}
	return nil
}
