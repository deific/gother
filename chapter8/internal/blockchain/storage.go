package blockchain

type Storage interface {
	Save(chain *Blockchain)
	Load() *Blockchain
	Update(chain *Blockchain)
}
