package page_test

import (
	p "Clotho/file_manager/page"
	"encoding/binary"
	"testing"
)

func asserNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("did not expect an error but got %q", err.Error())
	}
}

func asserError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestSetInt(t *testing.T) {
	t.Run("setting Int32 at offset 0", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = 0

		err := page.SetInt32(val, offset)

		asserNoError(t, err)

		var got int32

		binary.Decode(page.Bytes()[offset:], binary.NativeEndian, &got)

		if got != val {
			t.Errorf("got '%d' but expected '%d'", got, val)
		}
	})

	t.Run("setting Int32 at offset 13", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = 13

		err := page.SetInt32(val, offset)

		asserNoError(t, err)

		valsize := binary.Size(val)
		var got int32

		binary.Decode(page.Bytes()[offset:offset+p.Offset(valsize)], binary.NativeEndian, &got)

		if got != val {
			t.Errorf("got '%d' but expected '%d'", got, val)
		}
	})

	t.Run("setting Int32 at offset more than page size", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = 40

		err := page.SetInt32(val, offset)

		asserError(t, err)
	})

	t.Run("setting Int32 at negative offset", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = -40

		err := page.SetInt32(val, offset)

		asserError(t, err)
	})
}

func TestGetInt(t *testing.T) {
	t.Run("reading Int32 at Offset 0", func(t *testing.T) {
		page := p.NewPage(42)

		var expected int32 = 1337
		var offset p.Offset = 0

		err := page.SetInt32(expected, offset)

		asserNoError(t, err)

		val, err := page.GetInt32(offset)

		asserNoError(t, err)

		if val != expected {
			t.Errorf("got %d but expected %d", val, expected)
		}
	})

	t.Run("reading Int32 at Offset 13", func(t *testing.T) {
		page := p.NewPage(42)

		var expected int32 = 1337
		var offset p.Offset = 13

		err := page.SetInt32(expected, offset)

		asserNoError(t, err)

		val, err := page.GetInt32(offset)

		asserNoError(t, err)

		if val != expected {
			t.Errorf("got %d but expected %d", val, expected)
		}
	})

	t.Run("reading Int32 at offset more than page size", func(t *testing.T) {
		page := p.NewPage(42)

		_, err := page.GetInt32(42)

		asserError(t, err)
	})

	t.Run("reading Int32 at negative offset", func(t *testing.T) {
		page := p.NewPage(42)

		_, err := page.GetInt32(-42)

		asserError(t, err)
	})
}


func TestSetString(t *testing.T) {
	t.Run("setting a string at offset 0", func(t *testing.T) {
		
	})
}
