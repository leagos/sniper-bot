package consts

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	WBNBAddress           = "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"
	USDTAddressBSC        = "0x55d398326f99059fF775485246999027B3197955"
	BUSDAddress           = "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"
	PancakeRouterAddress  = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
	PancakeFactoryAddress = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"

	WETHAddress             = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	UniSwapV2RouterAddress  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniSwapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"

	//test net
	WBNBAddressTest           = "0x0cBE3D856B6A03B4dEeCbdfBB35e834Cad1aA50D"
	USDTAddressBSCTest        = "0x12cf459B5d90fA6D697e245e1d0CFAe2068fb766"
	BUSDAddressTest           = "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"
	PancakeRouterAddressTest  = "0x220b01F9192A435F3Fc33fD5F0b205BCec3d3fA1"
	PancakeFactoryAddressTest = "0x13700dAe4b84E546D7e493BA1d1FA15f190dB10E"

	//eth test
	WETHAddressTest             = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	UniSwapV2RouterAddressTest  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniSwapV2FactoryAddressTest = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
)

const (
	BscRpcAddr     = "https://bsc-dataseed.binance.org/"
	BscRpcAddrTest = "https://data-seed-prebsc-1-s1.binance.org:8545"
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
	USDTContractMap = map[string]common.Address{
		ChainTypeBsc:     common.HexToAddress(USDTAddressBSC),
		ChainTypeBscTest: common.HexToAddress(USDTAddressBSCTest),
	}
	BUSDContractMap = map[string]common.Address{
		ChainTypeBsc:     common.HexToAddress(BUSDAddress),
		ChainTypeBscTest: common.HexToAddress(BUSDAddressTest),
	}
)
