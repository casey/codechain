package hashchain

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/frankbraun/codechain/util/def"
	"github.com/mutecomm/mute/util/fuzzer"
)

func ignoreWindow(buf []byte) (int, int, error) {
	// determine position of last date field to exclude it from fuzzing
	c, err := Read(bytes.NewBuffer(buf))
	if err != nil {
		return 0, 0, err
	}
	var length int
	for i := 0; i < len(c.chain)-1; i++ {
		length += len(c.chain[i].String()) + 1
	}
	length += 65 // linkhash + white space
	ignoreStart := length * 8
	ignoreEnd := (length + 20 - 1) * 8 // 20 = length of date field
	return ignoreStart, ignoreEnd, nil
}

func TestFuzzer(t *testing.T) {
	if testing.Short() {
		t.Skip("skip fuzzer test in short mode")
	}
	buf, err := ioutil.ReadFile(filepath.Join("..", def.HashchainFile))
	if err != nil {
		t.Fatalf("ioutil.ReadFile() failed: %v", err)
	}

	ignoreStart, ignoreEnd, err := ignoreWindow(buf)
	if err != nil {
		t.Fatalf("ignoreWindow() failed: %v", err)
	}

	// Fuzz entire hash chain.
	divider := 100
	var (
		currentBit int
		failedBits int
		counter    int
		part       int
	)
	parts := len(buf) * 8 / divider
	fzzr := &fuzzer.SequentialFuzzer{
		Data: buf,
		// Let's not consider the last newline and the character before,
		// base64 can be modified somewhat without changing the byte sequence.
		End: (len(buf) - 2) * 8,
		TestFunc: func(buf []byte) error {
			defer func() {
				currentBit++
			}()
			if counter%divider == 0 {
				part++
				fmt.Printf("hashchain fuzzing part %d/%d\n", part, parts)
			}
			counter++
			if ignoreStart <= currentBit && currentBit <= ignoreEnd {
				return errors.New("ignore bit")
			}
			if _, err := Read(bytes.NewBuffer(buf)); err != nil {
				return err
			}
			failedBits++
			return nil
		},
	}
	if ok := fzzr.Fuzz(); !ok {
		t.Errorf("fuzzer failed (%d bits of %d)", failedBits, fzzr.End)
	}
}
