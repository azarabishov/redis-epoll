package core

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func readLength(buf *bytes.Buffer) (int64, error) {
	s, err := readStringUntilSr(buf)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func readStringUntilSr(buf *bytes.Buffer) (string, error) {
	s, err := buf.ReadString('\r')
	if err != nil {
		return "", err
	}
	// increamenting to skip `\n`
	buf.ReadByte()
	return s[:len(s)-1], nil
}

func readSimpleString(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	return readStringUntilSr(buf)
}

func readError(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	return readStringUntilSr(buf)
}

func readInt64(c io.ReadWriter, buf *bytes.Buffer) (int64, error) {
	s, err := readStringUntilSr(buf)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func readBulkString(c io.ReadWriter, buf *bytes.Buffer) (string, error) {
	len, err := readLength(buf)
	if err != nil {
		return "", err
	}

	var bytesRem int64 = len + 2 // 2 for \r\n
	bytesRem = bytesRem - int64(buf.Len())
	for bytesRem > 0 {
		tbuf := make([]byte, bytesRem)
		n, err := c.Read(tbuf)
		if err != nil {
			return "", nil
		}
		buf.Write(tbuf[:n])
		bytesRem = bytesRem - int64(n)
	}

	bulkStr := make([]byte, len)
	_, err = buf.Read(bulkStr)
	if err != nil {
		return "", err
	}

	// moving buffer pointer by 2 for \r and \n
	buf.ReadByte()
	buf.ReadByte()

	// reading `len` bytes as string
	return string(bulkStr), nil
}

func readArray(c io.ReadWriter, buf *bytes.Buffer, rp *RESPParser) (interface{}, error) {
	count, err := readLength(buf)
	if err != nil {
		return nil, err
	}

	var elems []interface{} = make([]interface{}, count)
	for i := range elems {
		elem, err := rp.DecodeOne()
		if err != nil {
			return nil, err
		}
		elems[i] = elem
	}
	return elems, nil
}

func encodeString(v string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return encodeString(v)
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case []string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, b := range value.([]string) {
			buf.Write(encodeString(b))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(v), buf.Bytes()))
	case []*Obj:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, b := range value.([]*Obj) {
			buf.Write(Encode(b.Value, false))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(v), buf.Bytes()))
	case []interface{}:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, b := range value.([]string) {
			buf.Write(Encode(b, false))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(v), buf.Bytes()))
	case []int64:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, b := range value.([]int64) {
			buf.Write(Encode(b, false))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(v), buf.Bytes()))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	default:
		return RESP_NIL
	}
}
