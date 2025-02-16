package serz

import (
	"encoding/binary"
	"fmt"
	"github.com/phith0n/zkar/commons"
)

type TCReference struct {
	Handler         uint32
	Flag            byte
	Object          *TCObject
	Class           *TCClass
	NormalClassDesc *TCClassDesc
	ProxyClassDesc  *TCProxyClassDesc
	String          *TCString
	Array           *TCArray
	Enum            *TCEnum
}

func (r *TCReference) ToBytes() []byte {
	result := []byte{JAVA_TC_REFERENCE}
	bs := commons.NumberToBytes(r.Handler)
	return append(result, bs...)
}

func (r *TCReference) ToString() string {
	var b = commons.NewPrinter()
	b.Printf("TC_REFERENCE - %s", commons.Hexify(JAVA_TC_REFERENCE))
	b.IncreaseIndent()
	b.Printf("@Handler - %v - %s", r.Handler, commons.Hexify(r.Handler))
	return b.String()
}

func readTCReference(stream *ObjectStream) (*TCReference, error) {
	// read JAVA_TC_REFERENCE flag
	_, _ = stream.ReadN(1)

	bs, err := stream.ReadN(4)
	if err != nil {
		return nil, fmt.Errorf("read JAVA_TC_REFERENCE failed on index %v", stream.CurrentIndex())
	}

	handler := binary.BigEndian.Uint32(bs)
	reference := &TCReference{
		Handler: handler,
	}

	obj := stream.GetReference(handler)
	if obj != nil {
		switch obj := obj.(type) {
		case *TCObject:
			reference.Flag = JAVA_TC_OBJECT
			reference.Object = obj
		case *TCClass:
			reference.Flag = JAVA_TC_CLASS
			reference.Class = obj
		case *TCClassDesc:
			reference.Flag = JAVA_TC_CLASSDESC
			reference.NormalClassDesc = obj
		case *TCProxyClassDesc:
			reference.Flag = JAVA_TC_PROXYCLASSDESC
			reference.ProxyClassDesc = obj
		case *TCString:
			reference.Flag = JAVA_TC_STRING
			reference.String = obj
		case *TCArray:
			reference.Flag = JAVA_TC_ARRAY
			reference.Array = obj
		case *TCEnum:
			reference.Flag = JAVA_TC_ENUM
			reference.Enum = obj
		default:
			goto Failed
		}

		return reference, nil
	}

Failed:
	return nil, fmt.Errorf("object reference %v is not found on index %v", handler, stream.CurrentIndex())
}
