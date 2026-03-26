package core

type BlockChain struct {
	headers   []*Header
	store     Storage
	validator Validator
}

func NewBlockChain(genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		headers: []*Header{},
		store:   NewMemorystore(),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *BlockChain) SetValidator(v Validator) {
	bc.validator = v
}
func (bc *BlockChain) AddBlock(b *Block) error {
	if err := bc.validator.ValidatorBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *BlockChain) HasBlock(heigth uint32) bool {
	return heigth <= bc.Heigth()
}

func (bc *BlockChain) Heigth() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)
	return bc.store.Put(b)
}
func (bc *BlockChain) addgenesisBlock(b *Block) {

}
