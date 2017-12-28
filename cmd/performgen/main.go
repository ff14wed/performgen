package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ff14wed/performgen"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	output, err := mainWithError(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Print(output)
}

func mainWithError(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString(byte(0))
	if err != nil && err != io.EOF {
		return "", err
	}
	segments, err := performgen.Generate(input)
	if err != nil {
		return "", err
	}
	output := bytes.NewBufferString("data,duration(ms)\n")
	writer := bufio.NewWriter(output)
	for _, segment := range segments {
		_, _ = writer.WriteString(segment.Block.String())
		_ = writer.WriteByte(',')
		ms := int64(segment.Length / time.Millisecond)
		_, _ = writer.WriteString(strconv.FormatInt(ms, 10))
		_ = writer.WriteByte('\n')
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}
	return output.String(), nil
}
