package blockchain

import "sync"

type Mempool struct {
	Transactions []*Transaction
	mu           sync.Mutex
}

func NewMempool() *Mempool {
	return &Mempool{
		Transactions: make([]*Transaction, 0),
	}
}

func (m *Mempool) AddTransaction(trans *Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Transactions = append(m.Transactions, trans)
}

func (m *Mempool) GetTransactions() []*Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Transactions
}

func (m *Mempool) ClearTransactions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Transactions = make([]*Transaction, 0)
}
