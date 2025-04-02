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

		got, err := page.SetInt32(val, offset)

		asserNoError(t, err)

		expected, err := p.EncodingSize(val)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("SetInt failed: got '%d' but expected '%d' written bytes", got, val)
		}
	})

	t.Run("setting Int32 at offset 13", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = 13

		got, err := page.SetInt32(val, offset)

		asserNoError(t, err)

		expected, err := p.EncodingSize(val)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("SetInt failed: got '%d' but expected '%d' written bytes", got, val)
		}
	})

	t.Run("setting Int32 at offset more than page size", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = 40

		_, err := page.SetInt32(val, offset)

		asserError(t, err)
	})

	t.Run("setting Int32 at negative offset", func(t *testing.T) {
		page := p.NewPage(42)

		var val int32 = 1337
		var offset p.Offset = -40

		_, err := page.SetInt32(val, offset)

		asserError(t, err)
	})
}

func TestGetInt(t *testing.T) {
	t.Run("reading Int32 at Offset 0", func(t *testing.T) {
		page := p.NewPage(42)

		var expected int32 = 1337
		var offset p.Offset = 0

		_, err := page.SetInt32(expected, offset)

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

		_, err := page.SetInt32(expected, offset)

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

func TestEncodingSize(t *testing.T) {
	t.Run("int32 Encoding size", func(t *testing.T) {

		var val int32 = 42

		expected := binary.Size(val)

		got, err := p.EncodingSize(val)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("size : got '%d' but expected '%d'", got, expected)
		}
	})

	t.Run("string Encoding size", func(t *testing.T) {

		var val string = "blabla"

		expected := binary.Size([]byte(val)) + binary.Size(int32(1))

		got, err := p.EncodingSize(val)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("size: got '%d' but expected '%d'", got, expected)
		}
	})

	t.Run("unsuported type Encoding size", func(t *testing.T) {

		var val map[int]int

		expected := binary.Size(val)

		got, err := p.EncodingSize(val)

		asserError(t, err)

		if got != expected {
			t.Errorf("size: got '%d' but expected '%d'", got, expected)
		}
	})

	t.Run("unsuported platform specific (e.g int, uint) type Encoding size", func(t *testing.T) {

		var val int

		expected := binary.Size(val)

		got, err := p.EncodingSize(val)

		asserError(t, err)

		if got != expected {
			t.Errorf("size: got '%d' but expected '%d'", got, expected)
		}
	})

}

func TestSetGetBlob(t *testing.T) {
	t.Run("setting & getting a blob at offset 0", func(t *testing.T) {
		page := p.NewPage(42)

		blob := []byte{1, 2, 3, 4, 5, 6}
		offset := p.Offset(0)

		got, err := page.SetBlob(blob, offset)

		asserNoError(t, err)

		expected, err := p.EncodingSize(blob)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("Setblob failed: got '%v' but expected '%v' writen bytes", got, expected)
		}
	})

	t.Run("setting & getting a blob at offset 13", func(t *testing.T) {
		page := p.NewPage(42)

		blob := []byte{1, 2, 3, 4, 5, 6}
		offset := p.Offset(13)

		got, err := page.SetBlob(blob, offset)

		asserNoError(t, err)

		expected, err := p.EncodingSize(blob)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("Setblob failed: got '%v' but expected '%v' writen bytes", got, expected)
		}
	})

	t.Run("setting & getting a blob at offset more than page size", func(t *testing.T) {
		page := p.NewPage(42)

		blob := []byte{1, 2, 3, 4, 5, 6}
		offset := p.Offset(42)

		_, err := page.SetBlob(blob, offset)

		asserError(t, err)

		_, err = page.GetBlob(offset)

		asserError(t, err)
	})

	t.Run("setting & getting a blob at a negative offset", func(t *testing.T) {
		page := p.NewPage(42)

		blob := []byte{1, 2, 3, 4, 5, 6}
		offset := p.Offset(-42)

		_, err := page.SetBlob(blob, offset)

		asserError(t, err)

		_, err = page.GetBlob(offset)

		asserError(t, err)
	})
}

func TestSetGetString(t *testing.T) {
	t.Run("setting & getting a string at offset 0", func(t *testing.T) {
		page := p.NewPage(42)

		str := "hello there"
		offset := p.Offset(0)

		got, err := page.SetString(str, offset)

		asserNoError(t, err)

		expected, err := p.EncodingSize(str)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("SetString failed: got '%d' but expected '%d' written bytes", got, expected)
		}
	})

	t.Run("setting & getting a string at offset 13", func(t *testing.T) {
		page := p.NewPage(42)

		str := "hello there"
		offset := p.Offset(13)

		got, err := page.SetString(str, offset)

		asserNoError(t, err)

		expected, err := p.EncodingSize(str)

		asserNoError(t, err)

		if got != expected {
			t.Errorf("SetString failed: got '%d' but expected '%d' written bytes", got, expected)
		}
	})
}
