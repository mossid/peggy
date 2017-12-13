package etend

import (

    sdk "github.com/cosmos/cosmos-sdk" // dev branch

    "github.com/ethereum/go-ethereum/common"
)

const (
    NameETEnd = "etend"
    ByteDeposit = byte(0xe8)
    ByteWithdraw = byte(0xe9)
    ByteTransfer = byte(0xea)
)

var (
    TypeDeposit = NameETEnd + "/deposit"
    TypeWithdraw = NameETEnd + "/withdraw"
    TypeTransfer = NameETEnd + "/transfer"
)

type DepositTx struct {
    To [20]byte
    Value []byte //big.Int.Bytes    
    Token common.Address

}

func (tx DepositTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx DepositTx) ValidateBasic() error {
    return nil
}

type WithdrawTx struct {
    To [20]byte
    Value int64
    Token string
    OriginChain string
    Sequence uint64
}


