package runner

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	return &ethRunner{
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		uniAbi:       uniAbi,
		ethClientMap: make(map[string]*ethclient.Client),
	}
}

func (e *ethRunner) SniperDxsale(chain string) {
	dxsaleContractAddress := common.HexToAddress(viper.GetString("targetContract"))
	value, _ := big.NewFloat(viper.GetFloat64("buyingBnbAmount") * params.Ether).Int(nil)
	estimateTransferGasData, err := e.uniAbi.Pack("transfer", dxsaleContractAddress, value)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	gas, err := e.getClient(chain).EstimateGas(ctx, ethereum.CallMsg{
		To:   &dxsaleContractAddress,
		Data: estimateTransferGasData,
	})

	if err != nil {
		if err.Error() != "execution reverted" {
			log.Fatal(err)
		}
		interval := viper.GetInt64("sniperInterval")
		log.Printf("contract not active, retry in %d ms", interval)
		time.AfterFunc(time.Duration(interval)*time.Millisecond, func() {
			e.SniperDxsale(chain)
		})
		return
	}

	if viper.GetUint64("gasLimit") < gas {
		log.Println("config gas limit less than estimate gas ", gas, "auto set to estimate gasLimit")
		viper.Set("gasLimit", gas)
	}

	log.Println("EstimateGas", gas, "ready to transfer")

	e.transfer(ctx, chain, dxsaleContractAddress, value)
}

func (e *ethRunner) transfer(ctx context.Context, chain string, toAddress common.Address, value *big.Int) {

	nonce, err := e.getClient(chain).PendingNonceAt(ctx, e.fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	println(nonce)

	tx, err := types.SignNewTx(e.privateKey, types.LatestSignerForChainID(big.NewInt(consts.ChainIdMap[chain])), &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: big.NewInt(viper.GetInt64("gasPrice") * params.GWei),
		Gas:      viper.GetUint64("gasLimit"),
		To:       &toAddress,
		Value:    value,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = e.getClient(chain).SendTransaction(ctx, tx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
}

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

func (e *ethRunner) SniperUniCake(chain string, quickMode bool) {
	wrapperTokenAddress := consts.UniSwapWrapperTokenContractMap[chain]
	targetTokenAddress := common.HexToAddress(viper.GetString("targetContract"))
	log.Printf("token address %s", targetTokenAddress)
	factory, err := uniswap.NewUniswapV2(consts.UniSwapFactoryContractMap[chain], e.getClient(chain))
	if err != nil {
		log.Fatal(err)
	}
	router, err := uniswap.NewUniswapV2(consts.UniSwapRouterContractMap[chain], e.getClient(chain))
	if err != nil {
		log.Fatal(err)
	}
	ethToken, err := uniswap.NewUniswapV2(wrapperTokenAddress, e.getClient(chain))
	if err != nil {
		log.Fatal(err)
	}

	interval := time.Duration(viper.GetInt64("sniperInterval"))
	path := []common.Address{wrapperTokenAddress, targetTokenAddress}
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
	log.Println(quickMode)
	if quickMode {
		auth.NoSend = true
		tx, err := router.SwapExactETHForTokens(auth, amountOutMin, path, e.fromAddress, big.NewInt(time.Now().Add(2*time.Minute).Unix()))
		if err != nil {
			log.Fatal(err)
		}
		to := consts.UniSwapRouterContractMap[chain]
	estimat:
		_, err = e.getClient(chain).EstimateGas(context.Background(), ethereum.CallMsg{
			From:  e.fromAddress,
			To:    &to,
			Data:  tx.Data(),
			Value: amountIn,
		})
		if err != nil {
			log.Printf("estimate GasLimit failed, retry in %d s", interval)
			time.Sleep(interval * time.Millisecond)
			goto estimat
		}
		err = e.getClient(chain).SendTransaction(context.Background(), tx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Transaction has been sent, tx hash: %s", tx.Hash().Hex())
		return
	}

getPair:
	pairAddress, err := factory.GetPair(nil, wrapperTokenAddress, targetTokenAddress)
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
		balance, err := ethToken.BalanceOf(nil, pairAddress)
		if err != nil {
			log.Fatal(err)
		}
		if balance.Cmp(minPoolLiquidityAdded) >= 1 {
			log.Printf("pool liquidity %s, start to buy ------", utils.WeiToEtherFloatByDecimals(18, balance).String())
			break
		}
		log.Println("pool liquidity too low:", utils.WeiToEtherFloatByDecimals(18, balance).String(), "BNB")
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
