package rpc

import (
	"encoding/hex"
	"gother/chapter8/internal/utils"
	"gother/chapter8/internal/wallet"
)

type JsonRpcService struct {
}

type walletInfoReq struct {
	address string
}
type walletInfoRes struct {
	address string
	pubKey  string
	refName string
}

func (cli *JsonRpcService) walletInfo(req *walletInfoReq, res *walletInfoRes) error {
	wlt, err := wallet.LoadWallet(req.address)
	utils.Handle(err)
	refList := wallet.LoadRefList()

	res.address = req.address
	res.pubKey = hex.EncodeToString(wlt.PublicKey)
	res.refName, _ = refList.FindRefName(req.address)

	return nil
}
