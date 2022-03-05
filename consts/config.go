package consts

var DefaultConfig = `
gasLimit: 500000

#gas价格
gasPrice: 10
# eip1559 gas tip, bsc not support
gasTip: 10

#购买数量(单位BNB)
buyingBnbAmount: 0.5

# 池子最低BNB数量，否则不会购买
minPoolLiquidityAdded: 2

# 滑点，0为自动
slippage: 12

#钱包私钥
privateKey: 

# 监听间隔（毫秒） 1秒 = 1000毫秒
sniperInterval: 1000

#代币合约地址，替换成要购买的合约地址
targetContract: 0x31e7ddebc4b4c1a9ba91a761390445f887354b25
`
