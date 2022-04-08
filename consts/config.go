package consts

var DefaultConfig = `
gasLimit: 500000

#gas价格
gasPrice: 10
# eip1559 gas tip, bsc not support
gasTip: 10

#购买数量(单位BNB)
buyingBnbAmount: 0.5

# 池子大小，低于此值的池子不抢购
minPoolLiquidityAdded: 2

# 滑点，0为自动，仅普通模式下有效
slippage: 0

#钱包私钥
privateKey: 

# 监听间隔（毫秒） 1秒 = 1000毫秒
sniperInterval: 1000

# 开始抢购时间 格式 Y-m-d H:i:s 例如 2022-03-29 22:00:00
startTime: 2022-03-29 22:00:00

# 自动卖出比例(百分比，0为不自动卖出)
sellPercent: 0

# 自动卖出


#代币合约地址，替换成要购买的合约地址
targetContract: 0x31e7ddebc4b4c1a9ba91a761390445f887354b25


#运行命令：
# sniper-bot.exe cake -m f/t/n -p USDT/BUSD

#模式说明 -m mode
# f 抢跑模式 速度最快，烧gas，不支持滑点和最低池子要求
# t 定时模式，适合指定开始时间的合约，不支持滑点和最低池子要求
# n 普通模式（默认模式） ，监测到流动性才去购买，支持滑点和最低池子要求，其他模式不支持，速度较慢

#池子类型 -p USDT/BUSD
#默认bnb池，不需要传-p参数

#命令举例
#抢跑模式：设置开始时间startTime，到时间后就按照sniperInterval的频率开始发送交易，直到监测到购买成功才停止
sniper-bot.exe cake -m f  

`
