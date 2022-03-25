package grt

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

const maxLineLength = 5 * 2 << 19 // 5M

func (t *Transformer) ReadStream(r io.Reader) error {

	scanner := bufio.NewScanner(r)

	// as some secrets, e.g. certificates or keys are quite long, provide a custom buffer
	scannerBuffer := []byte{}
	scanner.Buffer(scannerBuffer, maxLineLength)

	emptyLineRE := regexp.MustCompile("^ *$")

	buf := bytes.Buffer{}
	objectIsEmpty := true
	for scanner.Scan() {
		line := scanner.Text()
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
