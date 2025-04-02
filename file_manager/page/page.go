package page

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type Page struct {
	bytes []byte
}

// Create and initialize a new Page
func NewPage(size uint) *Page {
	page := Page{
		bytes: make([]byte, size),
	}

	return &page
}

func EncodingSize(v any) (int, error) {

	switch t := v.(type) {
	case int32:
		return binary.Size(t), nil
	// case uint32:
	// 	return binary.Size(v), nil
	case string:
		buf := []byte(t)
		return EncodingSize(buf)
	case []byte:
		return binary.Size(int32(1)) + binary.Size(t), nil
	case int:
		return -1, errors.New("specify the size of the int (e.g: int8, in16...)")
	case uint:
		return -1, errors.New("specify the size of the int (e.g: uint8, uin16...)")
	default:
		return -1, errors.Errorf("Unsupported type: %t", t)
	}
}

// Encode the binary format of val (NativeEndian byte Order) into the
// page starting at offset, using Encode function from the binary std pkg.
//
// Return an error in case offset or offset + valsize is outside of page
// bounds
func (p *Page) SetInt32(val int32, offset Offset) (int, error) {

	var valSize int = binary.Size(val)

	if offset+valSize > cap(p.bytes) || offset < 0 {
		return 0, errors.New("attempting to write outside the page bounds")
	}

	buf := make([]byte, valSize)

	// encode val into buf
	count, err := binary.Encode(buf, binary.NativeEndian, val)

	if err != nil {
		return 0, errors.Wrapf(err, "failed to Encode val '%d', consumed bytes: '%d'", val, count)
	}

	start := offset
	end := offset + valSize

	n := copy(p.bytes[start:end], buf)

	return n, nil
}

// Decode the binary encoding (NativeEndian byte order) of type Int32
// from the page starting at offset, using Decode function from the binary std pkg
//
// Return the Decoded value and nil, or 0 and a descriptive error.
func (p *Page) GetInt32(offset Offset) (int32, error) {

	var val int32

	valsize := binary.Size(val)

	if offset > cap(p.bytes) || offset < 0 {
		return 0, errors.New("reading exceedes page bounds")
	}

	count, err := binary.Decode(p.bytes[offset:], binary.NativeEndian, &val)

	if err != nil || count < valsize {
		return 0, errors.Wrapf(err, "failed to Decode Int32 val '%d', consumed bytes: '%d'", val, count)
	}

	return val, nil
}

func (p *Page) SetBlob(blob []byte, offset Offset) (int, error) {
	len := int32(len(blob))

	size, err := EncodingSize(blob)

	if err != nil {
		return 0, err
	}

	if offset+size > cap(p.bytes) || offset < 0 {
		return 0, errors.New("attempting to write outside the page bounds")
	}

	n, err := p.SetInt32(len, offset)

	if err != nil {
		return 0, errors.Wrapf(err, "SetBlob failed")
	}

	prefixSize, _ := EncodingSize(len)

	start := offset + prefixSize
	end := offset + size

	n += copy(p.bytes[start:end], blob)

	return n, nil
}

func (p *Page) GetBlob(offset Offset) ([]byte, error) {
	len, _ := p.GetInt32(offset)
	prefixSize, _ := EncodingSize(len)

	if int(offset)+prefixSize+int(len) > cap(p.bytes) || offset < 0 {
		return nil, errors.New("attempting to write outside the page bounds")
	}

	buff := make([]byte, len)

	start := int(offset) + prefixSize
	end := int(offset) + prefixSize + int(len)

	copy(buff, p.bytes[start:end])

	return buff, nil
}

func (p *Page) SetString(str string, offset Offset) (int, error) {
	return p.SetBlob([]byte(str), offset)
}

func (p *Page) GetString(offset Offset) (string, error) {
	buff, err := p.GetBlob(offset)

	if err != nil {
		return "", errors.Wrap(err, "failed to get string")
	}

	str := string(buff)

	return str, nil
}

func (p *Page) Bytes() []byte {
	return p.bytes
}

func (p *Page) String() string {
	return "TODO: implement this function for the Page"
}
