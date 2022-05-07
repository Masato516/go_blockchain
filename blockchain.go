package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

func (b *Block) Print() {
	fmt.Printf("nonce: %d\n", b.nonce)
	fmt.Printf("previous_hash: %s\n", b.previousHash)
	fmt.Printf("timestamp: %d\n", b.timestamp)
	fmt.Printf("transactions: %v\n", b.transactions)
}

// ブロックをハッシュ化
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string // マイニング実行者のアドレス
}

// １つ目のブロックを作成
func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

// 新しいブロックを作成、Blockchainへの追加、transactionPoolを空にする
func (bc *Blockchain) CreateBlock(nonce int, previous_hash [32]byte) *Block {
	b := NewBlock(nonce, previous_hash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// transactionをblockchainのtransactionPoolに追加
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	transaction := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, transaction)
}

// BlockchainのTransactionPoolをコピー
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.sendBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previous_hash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	candidateBlock := Block{nonce, previous_hash, 0, transactions}
	candidateHashStr := fmt.Sprintf("%x", candidateBlock.Hash())
	fmt.Println(candidateHashStr)
	return candidateHashStr[:difficulty] == zeros
}

// 正解のnonceを返す
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().previousHash
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	// 自分のアドレスをTransactionPoolに入れてから、ProofOfWork(nonceを求める)を行う
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, bc.LastBlock().Hash())
	log.Println("action=mining, status=success")
	return true
}

// 引数のアドレスの仮想通貨の総額を返す
func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			} else if blockchainAddress == t.sendBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
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
	blockChain := NewBlockchain("My_Blockchain_Minor_Address")
	blockChain.Print()

	blockChain.AddTransaction("Gaethje", "Poirier", 2.01)
	blockChain.AddTransaction("Khabib", "Mcgregor", 10.187)
	blockChain.Mining()
	blockChain.Print()

	blockChain.AddTransaction("Chandler", "Oliveira", 2.01)
	blockChain.AddTransaction("Khabib", "Mcgregor", 10.187)
	blockChain.Mining()
	blockChain.Print()

	fmt.Println(blockChain)

	fmt.Printf("My_Blockchain_Minor_Address: %.1f\n", blockChain.CalculateTotalAmount("My_Blockchain_Minor_Address"))
}
