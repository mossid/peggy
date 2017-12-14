package etgate

import (

    sdk "github.com/cosmos/cosmos-sdk" // dev branch
    "github.com/cosmos/cosmos-sdk/modules/coin"
//    "github.com/cosmos/cosmos-sdk/errors"
    
    //"github.com/tendermint/iavl" // dev branch

    "github.com/ethereum/go-ethereum/common"
)

var (
    ByteUpdate    = byte(0xe0)
    ByteFinalize  = byte(0xe1)
//    ByteValChange = byte(0xe2)
    ByteDeposit   = byte(0xe3)
    ByteWithdraw  = byte(0xe4)
    ByteTransfer  = byte(0xe5)

    TypeUpdate    = NameETGate + "/update"
    TypeFinalize  = NameETGate + "/finalize"
//    TypeValChange = NameETGate + "/valchange"
    TypeDeposit   = NameETGate + "/deposit"
    TypeWithdraw  = NameETGate + "/withdraw"
    TypeTransfer  = NameETGate + "/transfer"
)

type UpdateTx struct {
    Header []byte
}

func (tx UpdateTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx UpdateTx) ValidateBasic() error {
    return nil
}

type FinalizeTx struct {
    Header []byte  
    Witness Header
}

func (tx FinalizeTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx FinalizeTx) ValidateBasic() error {
    return nil
}

type DepositTx struct {
    Proof LogProof
}

func (tx DepositTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx DepositTx) ValidateBasic() error {
    return nil
    /*
    // remove this part for efficiency?
    log, err := tx.Proof.Log()
    if err != nil {
        return err
    }

    deposit := new(Deposit)
    return depositabi.Unpack(deposit, "Deposit", log)
    */
}

type WithdrawTx struct {
    To [20]byte
    Value []byte
    Token common.Address
}

func (tx WithdrawTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx WithdrawTx) ValidateBasic() error {
    return nil
}

type TransferTx struct {
    DestChain string
    To [20]byte
    Value coin.Coins
}

func (tx TransferTx) Wrap() sdk.Tx {
    return sdk.Tx{tx}
}

func (tx TransferTx) ValidateBasic() error {
    if !tx.Value.IsValid() {
        return errInvalidCoins
    }
    return nil
}
