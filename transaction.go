package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
)

// 1. 定义交易结构
// 2. 提供创建交易方法
// 3. 创建挖矿交易
// 4. 根据交易调整程序

const Reward = 12.5

type Transaction struct {
	TxID      []byte      // 交易ID
	TxInputs  []*TxInput  // 交易输入数组
	TxOutputs []*TxOutput // 交易输出数组
	Timestamp uint64      // 交易产生时间戳
}

type TxInput struct {
	TxID  []byte // 引用的交易ID
	Index int    // 引用的output的索引值
	//Sig   string // 解锁脚本，我们用地址来模拟

	// 真正的数字签名，由r，s拼成的[]byte
	Signature []byte

	// 约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分(参考r,s传递)
	// 注意：是公钥，不是hash，也不是地址
	PubKey []byte
}

type TxOutput struct {
	Amount float64 // 转账金额
	//PubKeyHash string  // 锁定脚本，我们用地址模拟

	// 收款方的公钥hash，注意：是hash而不是公钥，也不是地址
	PubKeyHash []byte
}

// 给TxOutput提供一个创建方法，否则无法调用Lock
func NewTxOutput(amount float64, address string) *TxOutput {
	output := &TxOutput{
		Amount: amount,
	}
	output.Lock(address)
	return output
}

// 由于现在存储的字段是地址的公钥hash，所以无法直接创建TxOutput，
//  为了能够得到公钥hash，我们需要处理一下，写一个Lock函数
func (o *TxOutput) Lock(address string) {
	// 真正的锁定动作！！！
	o.PubKeyHash = GetPubKeyFromAddress(address)
}

// 添加交易的Hash ID（设置Tx的TxID）
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TxID = hash[:]
}

// 实现一个函数，判断当前的交易是否为挖矿交易
func (tx *Transaction) IsCoinBase() bool {
	// 1. 交易的input只有一个
	// 2. 交易id为空
	// 3. 交易的index为-1
	if len(tx.TxInputs) == 1 {
		input := tx.TxInputs[0]
		if bytes.Equal(input.TxID, []byte{}) && input.Index == -1 {
			return true
		}
	}
	return false
}

// 创建挖矿奖励的交易
func NewCoinBaseTx(addr string, data string) *Transaction {
	// 1. 校验地址
	if !IsValidAddress(addr) {
		log.Printf("address %s is invalid\n", addr)
		return nil
	}

	// 挖矿交易的特点:
	// 1. 只有一个input
	// 2. 无需引用交易id
	// 3. 无需引用output 的 index

	// 矿工由于挖矿时无需指定签名，所以PubKey这个字段可以由矿工自由填写数据，一般是填写矿池的名字
	input := &TxInput{
		TxID:      []byte{},
		Index:     -1,
		Signature: nil,
		PubKey:    []byte(data),
	}
	//output := &TxOutput{
	//	Amount:     Reward,
	//	PubKeyHash: addr,
	//}

	// 新的创建方法
	output := NewTxOutput(Reward, addr)

	// 对于挖矿交易来说，只有一个input和一个output
	tx := &Transaction{
		TxInputs:  []*TxInput{input},
		TxOutputs: []*TxOutput{output},
		Timestamp: uint64(time.Now().Unix()),
	}
	tx.SetHash()
	return tx
}

// 创建普通的转账交易
//  1. 找到最合理的UTXO集合 map[string][]int64
//  2. 将这些UTXO逐一转成input
//  3. 创建outputs
//  4. 如果有零钱要找零
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 1. 校验地址
	if !IsValidAddress(from) {
		log.Printf("address %s is invalid\n", from)
		return nil
	} else if !IsValidAddress(to) {
		log.Printf("address %s is invalid\n", to)
		return nil
	}

	// 1. 创建交易之后要进行数字签名->所以需要私钥->打开钱包"NewWallets()"
	// 2. 找到自己的钱包，根据地址返回自己的wallet
	// 3. 得到对应的公钥、私钥
	ws := NewWallets()
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		log.Printf("not found address %s, Tx create fail!\n", from)
		return nil
	}
	pubKey := wallet.PubKey
	privateKey := wallet.Private

	// 传递公钥的hash，而不是传递地址
	pubKeyHash := HashPubKey(pubKey)

	utxos, totalAmount := bc.FindNeedUTXOs(pubKeyHash, amount)
	if totalAmount < amount {
		log.Printf("Insufficient balance, your: %f, need: %f\n", totalAmount, amount)
		return nil
	}

	var inputs = make([]*TxInput, 0, 4)
	var outputs = make([]*TxOutput, 0, 4)

	// 创建交易输入，并将这些UTXO添加到inputs中
	for txID, indexArray := range utxos {
		for _, i := range indexArray {
			input := &TxInput{
				TxID:      []byte(txID),
				Index:     i,
				Signature: nil,
				PubKey:    pubKey,
			}
			inputs = append(inputs, input)
		}
	}

	// 创建交易输出
	//output := &TxOutput{
	//	Amount:     amount,
	//	PubKeyHash: to,
	//}
	output := NewTxOutput(amount, to)
	outputs = append(outputs, output)

	// 找零
	if totalAmount > amount {
		output = NewTxOutput(totalAmount-amount, from)
		outputs = append(outputs, output)
	}

	tx := &Transaction{
		TxInputs:  inputs,
		TxOutputs: outputs,
		Timestamp: uint64(time.Now().Unix()),
	}
	tx.SetHash()

	// 签名，交易创建的最后进行签名
	bc.SignTransaction(tx, privateKey)

	return tx
}

// 签名的具体实现，参数为：私钥、inputs里面所有引用的交易结构map[string]Transaction
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTxs map[string]*Transaction) {

	if tx.IsCoinBase() {
		return
	}

	// 1. 创建一个当前交易的txCopy，使用函数：TrimmedCopy()，要把Signature和PubKey字段设为nil
	txCopy := tx.TrimmedCopy()
	// 2. 循环遍历txCopy的inputs，得到input所引用的output公钥hash
	for i, input := range txCopy.TxInputs {
		prevTx, ok := prevTxs[string(input.TxID)]
		if !ok {
			log.Panic("invalid tx")
		}
		txCopy.TxInputs[i].PubKey = prevTx.TxOutputs[input.Index].PubKeyHash

		// 3. 生成要签名的数据，要签名的数据一定是hash值
		// 3.1 我们对每一个input都要签名一次，签名的数据是由当前input引用的output的hash+当前的outputs（都承载在当前这个txCopy里面）
		// 3.2 要对这个拼好的txCopy进行hash处理，SetHash得到TxID，这个TxID就是我们要签名的最终数据
		txCopy.SetHash()
		txCopy.TxInputs[i].PubKey = nil
		signDataHash := txCopy.TxID

		// 4. 执行签名动作得到r，s字节流
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic(err)
		}

		// 5. 放到我们所签名的input的Signature中
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TxInputs[i].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TxInput
	var outputs []*TxOutput
	for _, input := range tx.TxInputs {
		inputs = append(inputs, &TxInput{
			TxID:      input.TxID,
			Index:     input.Index,
			Signature: nil,
			PubKey:    nil,
		})
	}
	for _, output := range tx.TxOutputs {
		outputs = append(outputs, &TxOutput{
			Amount:     output.Amount,
			PubKeyHash: output.PubKeyHash,
		})
	}
	return &Transaction{
		TxID:      tx.TxID,
		TxInputs:  inputs,
		TxOutputs: outputs,
		Timestamp: tx.Timestamp,
	}
}

// 分析校验
//  所需要的数据：公钥、数据（txCopy，生成hash）、签名
//  我们要对每一个签名过的input进行校验
func (tx *Transaction) Verify(prevTxs map[string]*Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	// 1. 得到签名的数据
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TxInputs {
		prevTx, ok := prevTxs[string(input.TxID)]
		if !ok {
			log.Panic("invalid tx")
		}
		txCopy.TxInputs[i].PubKey = prevTx.TxOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		dataHash := txCopy.TxID

		// 2. 得到Signature，反推回r，s
		signature := input.Signature // 拆 r, s
		// 3. 拆解PubKey，得到原生的公钥X，Y
		pubKey := input.PubKey // 拆 X，Y

		// 定义两个辅助的big.Int
		r := new(big.Int)
		s := new(big.Int)
		// 拆分我们signature，前半部分给r，后半部分给s
		idx := len(signature) / 2
		r.SetBytes(signature[:idx])
		s.SetBytes(signature[idx:])

		// 定义两个辅助的big.Int
		x := new(big.Int)
		y := new(big.Int)
		// 拆分我们pubKey，前半部分给X，后半部分给Y
		idx = len(pubKey) / 2
		x.SetBytes(pubKey[:idx])
		y.SetBytes(pubKey[idx:])

		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), x, y}

		// 4. Verify
		if !ecdsa.Verify(&pubKeyOrigin, dataHash, r, s) {
			return false
		}
	}

	return true
}

func (tx *Transaction) String() string {
	var lines = make([]string, 0, 16)
	lines = append(lines, fmt.Sprintf("--- Transaction %x", tx.TxID))

	for i, input := range tx.TxInputs {
		lines = append(lines, fmt.Sprintf("    Input: %d", i))
		lines = append(lines, fmt.Sprintf("      TxID: %x", input.TxID))
		lines = append(lines, fmt.Sprintf("      Index: %d", input.Index))
		lines = append(lines, fmt.Sprintf("      Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("      PubKey: %x", input.PubKey))
	}

	for i, output := range tx.TxOutputs {
		lines = append(lines, fmt.Sprintf("    Output: %d", i))
		lines = append(lines, fmt.Sprintf("      Amount: %f", output.Amount))
		lines = append(lines, fmt.Sprintf("      PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
