package main

// 定义一个Wallets结构，它保存所有的wallet以及它的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

// 创建方法
func NewWallets() *Wallets {
	wallet := NewWallet()
	address := wallet.NewAddress()

	var wallets = &Wallets{
		WalletsMap: map[string]*Wallet{
			address: wallet,
		},
	}
	return wallets
}

// 保存方法，把新建的wallet添加进去

// 读取文件方法，把所有的wallet读出来
