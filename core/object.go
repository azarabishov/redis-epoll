package core

type Obj struct {
	TypeEncoding uint8
	LastAccessedAt uint32
	Value          interface{}
}

var OBJ_TYPE_STRING uint8 = 0 << 4

var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_EMBSTR uint8 = 8

var OBJ_TYPE_BYTELIST uint8 = 1 << 4
var OBJ_ENCODING_QINT uint8 = 0
var OBJ_ENCODING_QREF uint8 = 1

var OBJ_ENCODING_STACKINT uint8 = 2
var OBJ_ENCODING_STACKREF uint8 = 3

var OBJ_TYPE_BITSET uint8 = 1 << 5 // 00100000
var OBJ_ENCODING_BF uint8 = 2      // 00000010

func ExtractTypeEncoding(obj *Obj) (uint8, uint8) {
	return obj.TypeEncoding & 0b11110000, obj.TypeEncoding & 0b00001111
}
