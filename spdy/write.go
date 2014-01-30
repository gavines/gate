package spdy

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
)

func (f *SynStreamFrame) write(w io.Writer, buf *bytes.Buffer, zw *zlib.Writer) {
	headFirst := uint32(0x80020001)
	binary.Write(w, binary.BigEndian, headFirst)

	writeHeader(f.Header, zw)
	zheader := buf.Bytes()

	flagsLength := uint32(f.Flags<<24) + uint32(len(zheader)) + 10
	binary.Write(w, binary.BigEndian, flagsLength)

	binary.Write(w, binary.BigEndian, f.StreamId)
	binary.Write(w, binary.BigEndian, f.AssociatedId)

	priority := f.Priority << 14
	binary.Write(w, binary.BigEndian, priority)

	w.Write(zheader)

	log.Debug("Send frame: %v", f)
}

func writeHeader(header map[string]string, zw *zlib.Writer) {
	binary.Write(zw, binary.BigEndian, uint16(len(header)))

	for k, v := range header {
		binary.Write(zw, binary.BigEndian, uint16(len(k)))
		io.WriteString(zw, k)
		binary.Write(zw, binary.BigEndian, uint16(len(v)))
		io.WriteString(zw, v)
	}
	zw.Flush()
}

func (f *DataFrame) Write(w io.Writer) {
}
