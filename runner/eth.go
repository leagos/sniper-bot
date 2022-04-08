package runner

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/leagos/sniper-bot/consts"
	"github.com/leagos/sniper-bot/contract/uniswap"
	"github.com/leagos/sniper-bot/utils"
	"github.com/spf13/viper"
)

type ethRunner struct {
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	uniAbi       abi.ABI
	ethClientMap map[string]*ethclient.Client
}

var FromAddr common.Address

func NewEthRunner() *ethRunner {
	uniAbi, err := abi.JSON(strings.NewReader(uniswap.UniswapV2ABI))
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(viper.GetString("privateKey"))
	if err != nil {
		log.Fatal("privateKey error:", err)
	}

	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	FromAddr = fromAddress

	return &ethRunner{
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		uniAbi:       uniAbi,
		ethClientMap: make(map[string]*ethclient.Client),
	}
}

//func (e *ethRunner) SniperDxsale(chain string) {
//	dxsaleContractAddress := common.HexToAddress(viper.GetString("targetContract"))
//	value, _ := big.NewFloat(viper.GetFloat64("buyingBnbAmount") * params.Ether).Int(nil)
//	estimateTransferGasData, err := e.uniAbi.Pack("transfer", dxsaleContractAddress, value)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx := context.Background()
//
//	gas, err := e.getClient(chain).EstimateGas(ctx, ethereum.CallMsg{
//		To:   &dxsaleContractAddress,
//		Data: estimateTransferGasData,
//	})
//
//	if err != nil {
//		if err.Error() != "execution reverted" {
//			log.Fatal(err)
//		}
//		interval := viper.GetInt64("sniperInterval")
//		log.Printf("contract not active, retry in %d ms", interval)
//		time.AfterFunc(time.Duration(interval)*time.Millisecond, func() {
//			e.SniperDxsale(chain)
//		})
//		return
//	}
//
//	if viper.GetUint64("gasLimit") < gas {
//		log.Println("config gas limit less than estimate gas ", gas, "auto set to estimate gasLimit")
//		viper.Set("gasLimit", gas)
//	}
//
//	log.Println("EstimateGas", gas, "ready to transfer")
//
//	e.transfer(ctx, chain, dxsaleContractAddress, value)
//}

//func (e *ethRunner) transfer(ctx context.Context, chain string, toAddress common.Address, value *big.Int) {
//
//	nonce, err := e.getClient(chain).PendingNonceAt(ctx, e.fromAddress)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	println(nonce)
//
//	tx, err := types.SignNewTx(e.privateKey, types.LatestSignerForChainID(big.NewInt(consts.ChainIdMap[chain])), &types.LegacyTx{
//		Nonce:    nonce,
//		GasPrice: big.NewInt(viper.GetInt64("gasPrice") * params.GWei),
//		Gas:      viper.GetUint64("gasLimit"),
//		To:       &toAddress,
//		Value:    value,
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = e.getClient(chain).SendTransaction(ctx, tx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
//}

func (e *ethRunner) getClient(chain string) *ethclient.Client {
	client, ok := e.ethClientMap[chain]
	if ok {
		return client
	}
	var rpcAddress string
	switch chain {
	case consts.ChainTypeBsc:
		//rpcAddress = viper.GetString("bscWssUrl")
		rpcAddress = consts.BscRpcAddr
	case consts.ChainTypeBscTest:
		//rpcAddress = viper.GetString("ethWssUrl")
		rpcAddress = consts.BscRpcAddrTest
	case consts.ChainTypeEth:
		//rpcAddress = viper.GetString("ethWssUrl")
		rpcAddress = consts.EthRpcAddr
	case consts.ChainTypeEthTest:
		//rpcAddress = viper.GetString("ethWssUrl")
		rpcAddress = consts.EthRpcAddrTest
	default:
		panic("not support chain")
	}
	log.Println("chain Type:", chain)

	//wssClient, err := rpc.DialWebsocket(context.Background(), rpcAddress, "")
	//if err != nil {
	//	panic("Dial Websocket failed")
	//}
	//client = ethclient.NewClient(wssClient)
	//if err != nil {
	//	log.Fatal(err)
	//}
	client, err := ethclient.Dial(rpcAddress)
	if err != nil {
		log.Fatal(err)
	}
	e.ethClientMap[chain] = client
	return client
}

func (e *ethRunner) SniperUniCake(chain string, mode string, poolType string) {
	wrapperTokenAddress := consts.UniSwapWrapperTokenContractMap[chain]
	USDTAddress := consts.USDTContractMap[chain]
	BUSDAddress := consts.BUSDContractMap[chain]
	targetTokenAddress := common.HexToAddress(viper.GetString("targetContract"))
	log.Printf("Your address %s", e.fromAddress)
	log.Printf("token address %s", targetTokenAddress)
	factory, err := uniswap.NewUniswapV2(consts.UniSwapFactoryContractMap[chain], e.getClient(chain))
	if err != nil {
		log.Fatal(err)
	}
	router, err := uniswap.NewUniswapV2(consts.UniSwapRouterContractMap[chain], e.getClient(chain))
	if err != nil {
		log.Fatal(err)
	}

	var poolToken *uniswap.UniswapV2
	var poolTokenAddress common.Address

	interval := time.Duration(viper.GetInt64("sniperInterval"))
	var path []common.Address
	switch strings.ToUpper(poolType) {
	case "USDT":
		path = []common.Address{wrapperTokenAddress, USDTAddress, targetTokenAddress}
		poolToken, err = uniswap.NewUniswapV2(USDTAddress, e.getClient(chain))
		if err != nil {
			log.Fatal(err)
		}
		poolTokenAddress = USDTAddress
	case "BUSD":
		path = []common.Address{wrapperTokenAddress, BUSDAddress, targetTokenAddress}
		poolToken, err = uniswap.NewUniswapV2(BUSDAddress, e.getClient(chain))
		if err != nil {
			log.Fatal(err)
		}
		poolTokenAddress = BUSDAddress
	default:
		path = []common.Address{wrapperTokenAddress, targetTokenAddress}
		poolToken, err = uniswap.NewUniswapV2(wrapperTokenAddress, e.getClient(chain))
		if err != nil {
			log.Fatal(err)
		}
		poolTokenAddress = wrapperTokenAddress
	}
	amountIn, _ := big.NewFloat(viper.GetFloat64("buyingBnbAmount") * params.Ether).Int(nil)
	amountOutMin := big.NewInt(0)

	//### auth ####
	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, big.NewInt(consts.ChainIdMap[chain]))
	if err != nil {
		log.Fatal(err)
	}
	auth.Value = amountIn
	auth.GasLimit = viper.GetUint64("gasLimit")
	auth.GasPrice = big.NewInt(viper.GetInt64("gasPrice") * params.GWei)

	timeFormat := "2006-01-02 15:04:05"
	startTime := viper.GetString("startTime")
	startAt, err := time.ParseInLocation(timeFormat, startTime, time.Local)
	if err != nil {
		log.Fatalf("开始时间格式错误:%s", err)
	}
	startTimestamp := startAt.Unix()
	//抢跑模式
	if mode == "f" {
		targetToken, err := uniswap.NewUniswapV2(targetTokenAddress, e.getClient(chain))
		if err != nil {
			log.Fatal(err)
		}
		tokenBalance, err := targetToken.BalanceOf(nil, FromAddr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Init token balance %s", utils.WeiToEtherFloatByDecimals(9, tokenBalance).String())
		for {
			//到了指定时间才开始发起交易
			block, _ := e.getClient(chain).BlockByNumber(context.Background(), nil)
			log.Printf("当前区块时间：%s,%d", time.Unix(int64(block.Time()), 0).In(time.Local).Format(timeFormat), block.Time())
			if block.Time()+6 < uint64(startTimestamp) {
				time.Sleep(interval * time.Millisecond)
				log.Printf("未到开始时间:%s,%d", startTime, startTimestamp)
				continue
			}
			nowBalance, err := targetToken.BalanceOf(nil, FromAddr)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Now token balance %s", utils.WeiToEtherFloatByDecimals(9, nowBalance).String())
			buyAmount := new(big.Int).Sub(nowBalance, tokenBalance)
			if buyAmount.Cmp(big.NewInt(0)) >= 1 {
				log.Printf("Buy amount %s, over ------", utils.WeiToEtherFloatByDecimals(9, buyAmount).String())
				return
			}
			tx, err := router.SwapExactETHForTokens(auth, amountOutMin, path, e.fromAddress, big.NewInt(time.Now().Add(2*time.Minute).Unix()))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
			time.Sleep(interval * time.Millisecond)
		}
	}
	if mode == "t" {
		for {
			block, _ := e.getClient(chain).BlockByNumber(context.Background(), nil)
			log.Printf("当前区块时间：%s,%d", time.Unix(int64(block.Time()), 0).In(time.Local).Format(timeFormat), block.Time())
			if block.Time()+6 < uint64(startTimestamp) {
				time.Sleep(interval * time.Millisecond)
				log.Printf("未到开始时间:%s,%d", startTime, startTimestamp)
				continue
			}
			block0, _ := e.getClient(chain).BlockByNumber(context.Background(), nil)
			log.Printf("当前区块高度： %d,区块时间：%s", block0.Number(), time.Unix(int64(block0.Time()), 0).In(time.Local).Format(timeFormat))
			tx, err := router.SwapExactETHForTokens(auth, amountOutMin, path, e.fromAddress, big.NewInt(time.Now().Add(2*time.Minute).Unix()))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("交易已发送, hash: %s", tx.Hash().Hex())
			time.Sleep(2500 * time.Millisecond)

			block1, _ := e.getClient(chain).BlockByNumber(context.Background(), nil)
			log.Printf("当前区块高度： %d,区块时间：%s", block1.Number(), time.Unix(int64(block1.Time()), 0).In(time.Local).Format(timeFormat))
			tx1, err := router.SwapExactETHForTokens(auth, amountOutMin, path, e.fromAddress, big.NewInt(time.Now().Add(2*time.Minute).Unix()))
			log.Printf("交易已发送,hash: %s", tx1.Hash().Hex())
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

getPair:
	pairAddress, err := factory.GetPair(nil, poolTokenAddress, targetTokenAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("pairAddress:", pairAddress)

	if pairAddress == consts.ZeroAddress {
		log.Printf("pair not create, retry in %d s", interval)
		time.Sleep(interval * time.Millisecond)
		goto getPair
	}
	minPoolLiquidityAdded, _ := big.NewFloat(viper.GetFloat64("minPoolLiquidityAdded") * params.Ether).Int(nil)

	for {
		balance, err := poolToken.BalanceOf(nil, pairAddress)
		if err != nil {
			log.Fatal(err)
		}
		if balance.Cmp(minPoolLiquidityAdded) >= 1 {
			log.Printf("pool liquidity %s, start to buy ------", utils.WeiToEtherFloatByDecimals(18, balance).String())
			break
		}
		log.Println("pool liquidity too low:", utils.WeiToEtherFloatByDecimals(18, balance).String(), poolType)
		time.Sleep(interval * time.Millisecond)
	}
	if slippage := viper.GetInt64("slippage"); slippage != 0 && slippage < 100 {
		amounts, err := router.GetAmountsOut(nil, amountIn, path)
		if err != nil {
			log.Fatal(err)
		}
		amountOutMin = new(big.Int).Div(new(big.Int).Mul(amounts[1], big.NewInt(100-slippage)), big.NewInt(100))
	}

	tx, err := router.SwapExactETHForTokens(auth, amountOutMin, path, e.fromAddress, big.NewInt(time.Now().Add(2*time.Minute).Unix()))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
}
