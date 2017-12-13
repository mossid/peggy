package etgate

import (
    "math/big"

    sdk "github.com/cosmos/cosmos-sdk"

    common "github.com/ethereum/go-ethereum/common"

)

type Info struct {
    LastFinalized uint64
}

type Header struct {
    ParentHash common.Hash
    Hash common.Hash
    Number uint64
    ReceiptHash common.Hash
}

func (h1 Header) Equals(h2 Header) bool {
    return h1.ParentHash  == h2.ParentHash &&
           h1.Hash        == h2.Hash       &&
           h1.Number      == h2.Number     &&
           h1.ReceiptHash == h2.ReceiptHash
}

type SignedHeader struct {
    Header Header
    Signers []sdk.Actor
}

type DepositLog struct {
    To [20]byte
    Value *big.Int
    Token common.Address
    Chain []byte
    Sequence uint64
}
