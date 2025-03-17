package page_test

import "testing"
import p "Clotho/file_manager/page"


func TestGetInt(t *testing.T) {
	t.Run("reading Int with Offset 0", func(t *testing.T) {
		page := p.NewPage(42)


		_ = page.SetInt(1337, 0);

		expected := 1337
		val, err := page.GetInt(0)

		if err != nil {
			t.Errorf("got an error but not expecting it: %q", err.Error())
		}

		if val != expected {
			t.Errorf("got %d but expected %d", val, expected)
		}
	})
}