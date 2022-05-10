package rpc

import (
	"encoding/hex"
	"gother/chapter8/internal/utils"
	"gother/chapter8/internal/wallet"
)

type JsonRpcService struct {
}

type WalletInfoReq struct {
	Address string
}
type WalletInfoRes struct {
	Address string
	PubKey  string
	RefName string
}

func (s *JsonRpcService) WalletInfo(req *WalletInfoReq, res *WalletInfoRes) error {
	wlt, err := wallet.LoadWallet(req.Address)
	utils.Handle(err)
	refList := wallet.LoadRefList()

	res.Address = req.Address
	res.PubKey = hex.EncodeToString(wlt.PublicKey)
	res.RefName, _ = refList.FindRefName(req.Address)

	return nil
}
