package order_brc20_service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"ordbook-aggregation/controller/request"
	"ordbook-aggregation/controller/respond"
	"ordbook-aggregation/major"
	"ordbook-aggregation/model"
	"ordbook-aggregation/service/create_key"
	"ordbook-aggregation/service/inscription_service"
	"ordbook-aggregation/service/mongo_service"
	"ordbook-aggregation/service/oklink_service"
	"ordbook-aggregation/service/unisat_service"
	"ordbook-aggregation/tool"
	"strconv"
)

func ColdDownUtxo(req *request.ColdDownUtxo) (string, error){
	var (
		netParams *chaincfg.Params = GetNetParams(req.Net)
		err error
		fromPriKeyHex, fromSegwitAddress string = "", ""
		txRaw string = ""
		//latestUtxo *model.OrderUtxoModel
		utxoList []*model.OrderUtxoModel = make([]*model.OrderUtxoModel, 0)
		startIndex int64 = GetSaveStartIndex(req.Net, req.UtxoType)
	)

	fromPriKeyHex, fromSegwitAddress, err = create_key.CreateSegwitKey(netParams)
	if err != nil {
		return "", err
	}

	//latestUtxo, _ = mongo_service.GetLatestStartIndexUtxo(req.Net, req.UtxoType)
	//if latestUtxo != nil {
	//	startIndex = latestUtxo.SortIndex
	//}

	inputs := make([]*TxInputUtxo, 0)
	inputs = append(inputs, &TxInputUtxo{
		TxId:     req.TxId,
		TxIndex:  req.Index,
		PkScript: req.PkScript,
		Amount:   req.Amount,
		PriHex:   req.PriKeyHex,
	})
	addr, err := btcutil.DecodeAddress(fromSegwitAddress, netParams)
	if err != nil {
		return "", err
	}
	//addrHash, err := btcutil.NewAddressWitnessPubKeyHash(addr.ScriptAddress(), netParams)
	//if err != nil {
	//	fmt.Printf("NewAddressWitnessPubKeyHash err: %s\n", err.Error())
	//	return "", err
	//}
	pkScriptByte, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}
	pkScript := hex.EncodeToString(pkScriptByte)
	//count := req.Amount/req.PerAmount
	outputs := make([]*TxOutput, 0)
	for i := int64(0); i < req.Count; i++ {
		outputs = append(outputs, &TxOutput{
			Address: fromSegwitAddress,
			Amount:  int64(req.PerAmount),
		})

		utxoList = append(utxoList, &model.OrderUtxoModel{
			//UtxoId:     "",
			Net:           req.Net,
			UtxoType:      req.UtxoType,
			Amount:        req.PerAmount,
			Address:       fromSegwitAddress,
			PrivateKeyHex: fromPriKeyHex,
			TxId:          "",
			Index:         i,
			PkScript:      pkScript,
			UsedState:     model.UsedNo,
			//UseTx:      "",
			SortIndex: startIndex + i+1,
			Timestamp: tool.MakeTimestamp(),
		})
	}


	if req.ChangeAddress == "" {
		req.ChangeAddress = req.Address
	}
	tx, err := BuildCommonTx(netParams, inputs, outputs, req.ChangeAddress, req.FeeRate)
	if err != nil {
		fmt.Printf("BuildCommonTx err:%s\n", err.Error())
		return "", err
	}
	txRaw, err = ToRaw(tx)
	if err != nil {
		return "", err
	}
	for _, u := range utxoList {
		u.TxId = tx.TxHash().String()
		u.UtxoId = fmt.Sprintf("%s_%d", u.TxId, u.Index)

		_, err := mongo_service.SetOrderUtxoModel(u)
		if err != nil {
			major.Println(fmt.Sprintf("SetOrderUtxoModel for cold down err:%s", err.Error()))
			return "", err
		}
	}

	txId := ""
	//if req.Net == "testnet" {
	//	txResp, err := mempool_space_service.BroadcastTx(req.Net, txRaw)
	//	if err != nil {
	//		return "", err
	//	}
	//	txId = txResp
	//}else {
	//	txResp, err := oklink_service.BroadcastTx(txRaw)
	//	if err != nil {
	//		return "", err
	//	}
	//	txId = txResp.TxId
	//}

	txResp, err := unisat_service.BroadcastTx(req.Net, txRaw)
	if err != nil {
		return "", err
	}
	txId = txResp.Result

	//txResp, err := node.BroadcastTx(req.Net, txRaw)
	//if err != nil {
	//	return "", err
	//}
	//txId = txResp

	return txId, nil
}

func saveNewDummyFromBid(net string, out Output, priKeyHex string, index int64, txId string) error {
	startIndex := GetSaveStartIndex(net, model.UtxoTypeDummy)
	//startIndex := int64(0)
	//latestUtxo, _ := mongo_service.GetLatestStartIndexUtxo(net, model.UtxoTypeDummy)
	//if latestUtxo != nil {
	//	startIndex = latestUtxo.SortIndex
	//}
	netParams := GetNetParams(net)
	addr, err := btcutil.DecodeAddress(out.Address, netParams)
	if err != nil {
		return err
	}
	addrHash, err := btcutil.NewAddressPubKeyHash(addr.ScriptAddress(), netParams)
	if err != nil {
		return err
	}
	pkScriptByte, err := txscript.PayToAddrScript(addrHash)
	if err != nil {
		return err
	}
	pkScript := hex.EncodeToString(pkScriptByte)

	newDummy := &model.OrderUtxoModel{
		UtxoId:        fmt.Sprintf("%s_%d", txId, index),
		Net:           net,
		UtxoType:      model.UtxoTypeDummy,
		Amount:        out.Amount,
		Address:       out.Address,
		PrivateKeyHex: priKeyHex,
		TxId:          txId,
		Index:         index,
		PkScript:      pkScript,
		UsedState:     model.UsedNo,
		SortIndex:     startIndex + 1,
		Timestamp:     tool.MakeTimestamp(),
	}

	_, err = mongo_service.SetOrderUtxoModel(newDummy)
	if err != nil {
		major.Println(fmt.Sprintf("SetOrderUtxoModel from bid err:%s", err.Error()))
		return nil
	}
	return nil
}

func CollectionUtxo(req *request.CollectionUtxo) (string, error){
	var (
		netParams *chaincfg.Params = GetNetParams(req.Net)
		err error
		txRaw string = ""
		totalIn uint64 = 0
		totalAmount int64 = 0
	)
	inputs := make([]*TxInputUtxo, 0)
	for _, v := range req.UtxoList {
		inputs = append(inputs, &TxInputUtxo{
			TxId:     v.TxId,
			TxIndex:  v.Index,
			PkScript: v.PkScript,
			Amount:   v.Amount,
			PriHex:   req.PriKeyHex,
		})
		totalIn = totalIn + v.Amount
	}

	totalSize := int64(len(inputs)) * SpendSize + 1 * OutSize + OtherSize


	totalAmount = int64(totalIn)- totalSize*req.FeeRate-546

	fmt.Printf("totalSize:%d, totalIn:%d, totalSize*req.FeeRate:%d, totalAmount:%d\n", totalSize, totalIn, totalSize*req.FeeRate, totalAmount)

	outputs := make([]*TxOutput, 0)
	outputs = append(outputs, &TxOutput{
		Address: req.Address,
		Amount:  totalAmount,
	})

	tx, err := BuildCommonTx(netParams, inputs, outputs, req.Address, req.FeeRate)
	if err != nil {
		fmt.Printf("BuildCommonTx err:%s\n", err.Error())
		return "", err
	}
	txRaw, err = ToRaw(tx)
	if err != nil {
		return "", err
	}
	txId := ""
	//if req.Net == "testnet" {
	//	txResp, err := mempool_space_service.BroadcastTx(req.Net, txRaw)
	//	if err != nil {
	//		return "", err
	//	}
	//	txId = txResp
	//}else {
	//	txResp, err := oklink_service.BroadcastTx(txRaw)
	//	if err != nil {
	//		return "", err
	//	}
	//	txId = txResp.TxId
	//}

	txResp, err := unisat_service.BroadcastTx(req.Net, txRaw)
	if err != nil {
		return "", err
	}
	txId = txResp.Result

	//txResp, err := node.BroadcastTx(req.Net, txRaw)
	//if err != nil {
	//	return "", err
	//}
	//txId = txResp

	return txId, nil
}


//cold down the brc20 transfer
func ColdDownBrc20Transfer(req *request.ColdDownBrcTransfer) (*respond.Brc20TransferCommitResp, error){
	var (
		netParams *chaincfg.Params = GetNetParams(req.Net)
		_, platformAddressSendBrc20 string = GetPlatformKeyAndAddressSendBrc20(req.Net)
		transferContent string = fmt.Sprintf(`{"p":"brc-20", "op":"transfer", "tick":"%s", "amt":"%d"}`, req.Tick, req.InscribeTransferAmount)
		commitTxHash, revealTxHash, inscriptionId string = "", "", ""
		err error
		brc20BalanceResult *oklink_service.OklinkBrc20BalanceDetails
		availableBalance int64 = 0
		//transferContentMap map[string]interface{} = map[string]interface{}{
		//	"p":"brc-20",
		//	"op":"transfer",
		//	"tick":req.Tick,
		//	"amt":fmt.Sprintf("%d", req.InscribeTransferAmount),
		//}
	)
	//transferContentMapStr, _ := tool.ObjectToJson(transferContentMap)

	fmt.Println(transferContent)
	//fmt.Println(transferContentMapStr)
	brc20BalanceResult, err = oklink_service.GetAddressBrc20BalanceResult(platformAddressSendBrc20, req.Tick, 1, 50)
	if err != nil  {
		return nil, err
	}
	availableBalance, _ = strconv.ParseInt(brc20BalanceResult.AvailableBalance, 10, 64)
	if availableBalance < req.InscribeTransferAmount {
		return nil, errors.New("AvailableBalance not enough. ")
	}
	commitTxHash, revealTxHash, inscriptionId, err = inscription_service.InscribeOneData(netParams, req.PriKeyHex, platformAddressSendBrc20, transferContent, req.FeeRate, req.ChangeAddress)
	if err != nil {
		return nil, err
	}
	return &respond.Brc20TransferCommitResp{
		CommitTxHash:  commitTxHash,
		RevealTxHash:  revealTxHash,
		InscriptionId: inscriptionId,
	}, nil
}

//func ColdDownBrc20TransferBatch(req *request.ColdDownBrcTransferBatch) (*respond.Brc20TransferCommitResp, error){
//
//
//
//
//}