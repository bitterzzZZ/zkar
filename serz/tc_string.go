package serz

import (
	"fmt"
	"github.com/phith0n/zkar/commons"
)

type TCString struct {
	Utf     *TCUtf
	Handler uint32
}

func (so *TCString) ToBytes() []byte {
	var bs []byte
	var length = len(so.Utf.Data)
	if length <= 0xFFFF {
		bs = append(bs, JAVA_TC_STRING)
	} else {
		bs = append(bs, JAVA_TC_LONGSTRING)
	}

	return append(bs, so.Utf.ToBytes()...)
}

func (so *TCString) ToString() string {
	var b = commons.NewPrinter()
	var length = len(so.Utf.Data)
	if length <= 0xFFFF {
		b.Printf("TC_STRING - %s", commons.Hexify(JAVA_TC_STRING))
	} else {
		b.Printf("TC_LONGSTRING - %s", commons.Hexify(JAVA_TC_LONGSTRING))
	}
	b.IncreaseIndent()
	b.Printf("@Handler - %v", so.Handler)
	b.Print(so.Utf.ToString())
	return b.String()
}

func readTCString(stream *ObjectStream) (*TCString, error) {
	var s = new(TCString)
	var err error
	flag, _ := stream.ReadN(1)
	if flag[0] == JAVA_TC_STRING {
		s.Utf, err = readUTF(stream)
	} else if flag[0] == JAVA_TC_LONGSTRING {
		s.Utf, err = readLongUTF(stream)
	} else {
		return nil, fmt.Errorf("readTCString flag error on index %v", stream.CurrentIndex())
	}

	if err != nil {
		return nil, err
	}

	stream.AddReference(s)
	return s, nil
}
