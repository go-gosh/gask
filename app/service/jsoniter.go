//go:build jsoniter
// +build jsoniter

package service

import (
	"time"
	"unsafe"

	"github.com/json-iterator/go"
)

type timeFormatDecoder struct{}

func (t *timeFormatDecoder) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.UnixNano() == 0
}

func (t *timeFormatDecoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	stream.WriteString(ts.In(time.Local).Format(time.RFC3339))
}

func (t *timeFormatDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	parse, _ := time.Parse(time.RFC3339, iter.ReadString())
	*((*time.Time)(ptr)) = parse.In(time.Local)
}

func init() {
	jsoniter.RegisterTypeDecoder("time.Time", &timeFormatDecoder{})
	jsoniter.RegisterTypeEncoder("time.Time", &timeFormatDecoder{})
}
