package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

type MyWriter struct {
	writtenChunkSizes []int
	writtenData       bytes.Buffer
}

func (mw *MyWriter) Write(data []byte) (n int, err error) {
	mw.writtenChunkSizes = append(mw.writtenChunkSizes, len(data))
	n, err = mw.writtenData.Write(data)
	log.Printf("mw.Write() called with %v bytes, returning %v %v", len(data), n, err)
	return
}

func (mw *MyWriter) GetTotalOutput() []byte {
	return mw.writtenData.Bytes()

}

func TestChunkedCopy(t *testing.T) {

	CHUNK_SIZE := int64(337)

	bldr := strings.Builder{}

	for i := 0; i < 1000; i++ {
		bldr.WriteString(fmt.Sprintf("x%07d", i)) // can ignore return value
	}

	input_bytes := []byte(bldr.String())

	mw := &MyWriter{}

	bytes_written, err := chunkedCopy(mw, bytes.NewBuffer(input_bytes), CHUNK_SIZE) // use a strange chunk size just to prove it works

	if err != io.EOF {
		t.Errorf("chunkedCopy: expected EOF but got: %v", err)
	}

	if int(bytes_written) != len(input_bytes) {
		t.Errorf("chunkedCopy: did not copy full input: expected %v, got %v", len(input_bytes), bytes_written)
	}

	output_bytes := mw.GetTotalOutput()

	if !bytes.Equal(output_bytes, input_bytes) {
		t.Errorf("chunkedCopy: output did not match input")
	}

	// This is what we really want to check. All of the calls to Write (except the last) must be of size CHUNK_SIZE
	for i := 0; i < len(mw.writtenChunkSizes)-2; i++ {
		if mw.writtenChunkSizes[i] != int(CHUNK_SIZE) {
			t.Errorf("chunkedCopy: writtenChunkSizes[%d] = %v, expected %v", i, mw.writtenChunkSizes[i], CHUNK_SIZE)
		}
	}

}
