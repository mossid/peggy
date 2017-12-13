package etgate

import (
    "strings"

    abci "github.com/tendermint/abci/types"

    sdk "github.com/cosmos/cosmos-sdk"
    "github.com/cosmos/cosmos-sdk/state"
    "github.com/cosmos/cosmos-sdk/stack"
    "github.com/cosmos/cosmos-sdk/modules/auth"
    "github.com/cosmos/cosmos-sdk/modules/ibc"

    "github.com/tendermint/tmlibs/log"

    eth "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/rlp"

    "./abi"
    "../../contracts"
    "../etend"
)

const (
    NameETGate = "etgate"

    Delay = 2 // for now. 30+ for production?
    Pending = 6
)

var (
    depositabi abi.ABI
)

func init() {
    deposit, err := abi.JSON(strings.NewReader(contracts.ETGateABI))
    if err != nil {
        panic(err)
    }
    depositabi = deposit
}

type Handler struct { 
}

func (Handler) AssertDispatcher() {}

var _ stack.Dispatchable = Handler{}

func (Handler) Name() string {
    return NameETGate
}

func (h Handler) CheckTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx, next sdk.Checker) (res sdk.CheckResult, err error) {
    return
}

func (h Handler) DeliverTx(ctx sdk.Context, store state.SimpleDB, tx sdk.Tx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    err = tx.ValidateBasic()
    if err != nil {
        return res, err
    }

    switch t := tx.Unwrap().(type) {
    case UpdateTx:
        return h.updateTx(ctx, store, t)
    case FinalizeTx:
        return h.finalizeTx(ctx, store, t)
//    case ValChaingeTx:
//        return h.valChangeTx(ctx, store, t)
    case DepositTx:
        return h.depositTx(ctx, store, t, next)
    case WithdrawTx:
        return h.withdrawTx(ctx, store, t)
    }

    return next.DeliverTx(ctx, store, tx)
}

func (h Handler) updateTx(ctx sdk.Context, store state.SimpleDB, tx UpdateTx) (res sdk.DeliverResult, err error) {
    header, err := decodeHeader(tx.Header)
    if err != nil {
        return res, err
    }

    info := loadInfo(store)
    if header.Number < info.LastFinalized || header.Number > info.LastFinalized + Delay + Pending {
        return res, errHeaderOutOfVisibleRange
    }

    sender, err := getTxSender(ctx)
    if err != nil {
        return res, err
    }

    if !isValidator(store, sender) {
        return res, errNotValidator
    }

    signedHeaders := loadSignedHeaders(store, header.Number)
    if signedHeaders == nil {
        saveSignedHeaders(store, header.Number, []SignedHeader{ singleSignedHeader(header, sender) })
        return res, nil
    }

    index := -1
    for i, s := range signedHeaders {
        if s.Header.Equals(header) {
            index = i
            break
        }
    }

    if index == -1 {        
        saveSignedHeaders(store, header.Number, append(signedHeaders, singleSignedHeader(header, sender)))
        return res, nil
    }

    signed := signedHeaders[index]
    for _, s := range signed.Signers {
        if s.Equals(sender) {
            return res, ErrAlreadySignedHeader(sender)
        }
    }   
    signed.Signers = append(signed.Signers, sender)
    signedHeaders = append(signedHeaders, signed)
    saveSignedHeaders(store, header.Number, signedHeaders)
    return res, nil
}

func (h Handler) finalizeTx(ctx sdk.Context, store state.SimpleDB, tx FinalizeTx) (res sdk.DeliverResult, err error) {
    info := loadInfo(store)

    if info.LastFinalized != tx.Witness.Number - Delay - 1 {
        return res, errNoncontinuousFinalization
    }
    
    validators := loadValidators(store)

    signers := loadSigners(store, tx.Witness.Number, tx.Witness.Hash)

    sum := 0

    for _, s := range signers {
        flag := false
        for _, v := range validators {
            if v.Equals(s) {
                flag = true
                break
            }
        }
        if flag {
            sum++
        }
    }

    if sum * 3 < len(validators) * 2 {
        return res, errNotEnoughSigns
    }

    ancestor, err := getAncestor(store, tx.Witness)
    if err != nil {
        return res, err
    }

    parent, exists := loadFinalized(store, info.LastFinalized)
    if !exists {
        return res, ErrHeaderNotFound(info.LastFinalized) // must not be happend
    }

    if ancestor.ParentHash != parent.Hash {
        return res, errConflictingChain
    }

    info.LastFinalized += 1;
    saveInfo(store, info)

    saveFinalized(store, ancestor.Number, ancestor)

    removeSignedHeaders(store, ancestor.Number)

    return res, nil
}

func (h Handler) depositTx(ctx sdk.Context, store state.SimpleDB, tx DepositTx, next sdk.Deliver) (res sdk.DeliverResult, err error) {
    log, err := tx.Proof.Log()
    if err != nil {
        return res, ErrInvalidLogProof(log, err)
    }

    header, exists := loadFinalized(store, tx.Proof.Number)
    if !exists {
        return res, errLogHeaderNotFound
    }

    if !tx.Proof.IsValid(header.ReceiptHash) {
        return res, errInvalidLogProof
    }

    deposit := new(DepositLog)
    if err = depositabi.Unpack(deposit, "Deposit", log); err != nil {
        return res, errLogUnpackingError
    }

    increaseBalance(store, string(deposit.Chain), deposit.Value.Uint64())

    outTx := etend.DepositTx{
        To:    deposit.To,
        Value: deposit.Value.Bytes(),
        Token: deposit.Token,
    }

    packet := ibc.CreatePacketTx {
        DestChain:   string(deposit.Chain),
        Permissions: []sdk.Actor{},
        Tx:          outTx.Wrap(),
    }

    ibcCtx := ctx.WithPermissions(ibc.AllowIBC(NameETGate))
    _, err = next.DeliverTx(ibcCtx, store, packet.Wrap())
    
    return res, err
}

func (h Handler) InitState(l log.Logger, store state.SimpleDB, module, key, value string, cb sdk.InitStater) (log string, err error) {
    return 
}

func (h Handler) InitValidate(log log.Logger, store state.SimpleDB, vals []*abci.Validator, next sdk.InitValidater) {
    
}

func decodeHeader(headerb []byte) (Header, error) {
    var header eth.Header
    if err := rlp.DecodeBytes(headerb, &header); err != nil {
        return Header{}, err
    }

    return Header {
        ParentHash:  header.ParentHash,
        Hash:        header.Hash(),
        ReceiptHash: header.ReceiptHash,
        Number:      header.Number.Uint64(),
    }, nil
}

func getTxSender(ctx sdk.Context) (sender sdk.Actor, err error) {
    senders := ctx.GetPermissions("", auth.NameSigs)
    if len(senders) != 1 {
        return sender, errMissingSignature
    }
    return senders[0], nil
}

func isValidator(store state.SimpleDB, sender sdk.Actor) bool {
    validators := loadValidators(store)
    flag := false
    for _, v := range validators {
        if v.Equals(sender) {
            flag = true
            break
        }
    }

    return flag
}

func singleSignedHeader(header Header, signer sdk.Actor) (SignedHeader) {
    return SignedHeader{
        Header: header,
        Signers: []sdk.Actor{ signer },
    }
}

func getAncestor(store state.SimpleDB, header Header) (res Header, err error) {
    for i := 0; i < Delay; i++ {
        signedHeaders := loadSignedHeaders(store, header.Number)
        flag := false
        for _, s := range signedHeaders {
            if s.Header.Hash == header.ParentHash {
                flag = true
                header = s.Header
                break
            }
        }
        if !flag {
            return res, errAncestorNotFound
        }
    }
    return header, nil
}
