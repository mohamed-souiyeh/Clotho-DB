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

// Encode the binary format of val (NativeEndian byte Order) into the 
// page starting at offset, using Encode function from the binary std pkg.
//
// Return an error in case offset or offset + valsize is outside of page
// bounds 
func (p *Page) SetInt32(val int32, offset Offset) error {

	var valSize int = binary.Size(val)

	if int(offset)+valSize > cap(p.bytes) || offset < 0 {
		return errors.New("attempting to write outside the page bounds")
	}

	buf := make([]byte, valSize)

	// encode val into buf
	count, err := binary.Encode(buf, binary.NativeEndian, val)

	if err != nil {
		return errors.Wrapf(err, "failed to Encode val '%d', consumed bytes: '%d'", val, count)
	}

	copy(p.bytes[offset:offset+Offset(valSize)], buf)

	return nil
}

// Decode the binary encoding (NativeEndian byte order) of type Int32 
// from the page starting at offset, using Decode function from the binary std pkg
// 
// Return the Decoded value and nil, or 0 and a descriptive error.
func (p *Page) GetInt32(offset Offset) (int32, error) {

	var val int32

	valsize := binary.Size(val)

	if int(offset) > cap(p.bytes) || offset < 0{
		return 0, errors.New("reading exceedes page bounds")
	}

	count, err := binary.Decode(p.bytes[offset:], binary.NativeEndian, &val)

	if err != nil || count < valsize {
		return 0, errors.Wrapf(err, "failed to Decode Int32 val '%d', consumed bytes: '%d'", val, count)
	}

	return val, nil
}

func (p *Page) Bytes() []byte {
	return p.bytes
}

func (p *Page) String() string {
	return "TODO: implement this function for the Page"
}