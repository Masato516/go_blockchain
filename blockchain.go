package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
	}
}

func (b *Block) Print() {
	fmt.Printf("nonce: %d\n", b.nonce)
	fmt.Printf("previous_hash: %s\n", b.previousHash)
	fmt.Printf("timestamp: %d\n", b.timestamp)
	fmt.Printf("transactions: %v\n", b.transactions)
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Timestamp    int64    `json:"timestamp"`
		Transactions []string `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

// １つ目のブロックを作成
func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previous_hash [32]byte) *Block {
	b := NewBlock(nonce, previous_hash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

type Transaction struct {
	sendBlockchainAddress      string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sendBlockchainAddress: %s\n", t.sendBlockchainAddress)
	fmt.Printf(" recipientBlockchainAddress: %s\n", t.sendBlockchainAddress)
	fmt.Printf(" value: %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"send_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.sendBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	blockChain := NewBlockchain()
	blockChain.Print()

	blockChain.CreateBlock(1, blockChain.LastBlock().Hash())
	blockChain.Print()

	blockChain.CreateBlock(2, blockChain.LastBlock().Hash())
	blockChain.Print()

	t := NewTransaction("Aさん", "Bさん", 3.219)
	t.Print()
	jt, _ := t.MarshalJSON()
	fmt.Println(string(jt))
}
