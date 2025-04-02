package block

import "fmt"

type BlockID struct {
	Filename string
	Blknum   int
}

func (b *BlockID) String() string {
	return fmt.Sprintf("[file: %s, block: %d]", b.Filename, b.Blknum)
}
