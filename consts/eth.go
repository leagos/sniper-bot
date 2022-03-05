package consts

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	WBNBAddress           = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
	PancakeRouterAddress  = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
	PancakeFactoryAddress = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"

	WETHAddress             = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	UniSwapV2RouterAddress  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniSwapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"

	//test net
	WBNBAddressTest           = "0x1e33833a035069f42d68D1F53b341643De1C018D"
	PancakeRouterAddressTest  = "0x07d090e7FcBC6AFaA507A3441C7c5eE507C457e6"
	PancakeFactoryAddressTest = "0xd417A0A4b65D24f5eBD0898d9028D92E3592afCC"

	//eth test
	WETHAddressTest             = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	UniSwapV2RouterAddressTest  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniSwapV2FactoryAddressTest = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
)

const (
	BscRpcAddr     = "https://bsc-dataseed.binance.org/"
	BscRpcAddrTest = "https://data-seed-prebsc-2-s3.binance.org:8545/"
	EthRpcAddr     = "https://main-light.eth.linkpool.io/"
	EthRpcAddrTest = "https://main-light.eth.linkpool.io/"
)

const (
	ChainTypeEth     = "eth"
	ChainTypeEthTest = "ethTest"
	ChainTypeBsc     = "bsc"
	ChainTypeBscTest = "bscTest"
)

var (
	ZeroAddress = common.Address{}
	ChainIdMap  = map[string]int64{
		ChainTypeEth:     1,
		ChainTypeBsc:     56,
		ChainTypeEthTest: 3,
		ChainTypeBscTest: 97,
	}
	UniSwapWrapperTokenContractMap = map[string]common.Address{
		ChainTypeEth:     common.HexToAddress(WETHAddress),
		ChainTypeBsc:     common.HexToAddress(WBNBAddress),
		ChainTypeEthTest: common.HexToAddress(WETHAddressTest),
		ChainTypeBscTest: common.HexToAddress(WBNBAddressTest),
	}

	UniSwapFactoryContractMap = map[string]common.Address{
		ChainTypeBsc:     common.HexToAddress(PancakeFactoryAddress),
		ChainTypeBscTest: common.HexToAddress(PancakeFactoryAddressTest),
		ChainTypeEth:     common.HexToAddress(UniSwapV2FactoryAddress),
		ChainTypeEthTest: common.HexToAddress(UniSwapV2FactoryAddressTest),
	}

	UniSwapRouterContractMap = map[string]common.Address{
		ChainTypeEth:     common.HexToAddress(UniSwapV2RouterAddress),
		ChainTypeBsc:     common.HexToAddress(PancakeRouterAddress),
		ChainTypeEthTest: common.HexToAddress(UniSwapV2RouterAddressTest),
		ChainTypeBscTest: common.HexToAddress(PancakeRouterAddressTest),
	}
)
