package page




type Page struct {
	bytes []byte
}


func NewPage(size uint) *Page {
	page := Page {
		bytes : make([]byte, size),
	}

	return &page
}

func (p *Page) SetInt(val int, offset uint) error {

	
	return nil;
}

func (p *Page) GetInt(offset uint) (int, error) {
	return 0, nil;
}