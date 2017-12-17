package etgate

import (
    "math/big"
    "bytes"

    sdk "github.com/cosmos/cosmos-sdk"
    "github.com/cosmos/cosmos-sdk/state"

    wire "github.com/tendermint/go-wire"

    "github.com/ethereum/go-ethereum/common"
)

var (
    LastHeaderKey  = []byte{0x01} // uint64
    ValidatorsKey  = []byte{0x02} // []sdk.Actor

    InfoPrefix     = []byte{0x03} // string => Info
    BufferPrefix   = []byte{0x04} // uint => []SignedHeader
    FinalPrefix    = []byte{0x05} // uint => Header
    BalancePrefix  = []byte{0x06} // ChainTokenPair => uint
    WithdrawPrefix = []byte{0x07} // uint => []byte
)

// this repo is based on master branch of go-wrie\
// Marshal() is added on develop branch
// just copied it for now

func wireMarshal(o interface{}) ([]byte, error) {
    w, n, err := new(bytes.Buffer), new(int), new(error)
    wire.WriteBinary(o, w, n, err)
    if *err != nil {
        return nil, *err
    }   
    return w.Bytes(), nil
}

func wireUnmarshal(d []byte, ptr interface{}) error {
    r, n, err := bytes.NewBuffer(d), new(int), new(error)
    wire.ReadBinaryPtr(ptr, r, len(d), n, err)
    return *err
}

func marshal(o interface{}) []byte {
    res, err := wireMarshal(o)
    if err != nil {
        panic(err)
    }
    return res
}

func unmarshal(d []byte) interface{} {
    var ptr interface{}
    if err := wireUnmarshal(d, ptr); err != nil {
        panic(err)
    }
    return ptr
}

func GetInfoKey(chainid string) []byte {
    return append(InfoPrefix, []byte(chainid)...)
}

func GetBufferKey(height uint64) []byte { 
    return append(BufferPrefix, marshal(height)...)
}

func GetFinalKey(height uint64) []byte { 
    return append(FinalPrefix, marshal(height)...)
}

func GetBalanceKey(chainid string, token common.Address) []byte {
    return append(BalancePrefix, marshal(ChainTokenPair{chainid, token})...)
}

func GetWithdrawKey(seq uint64) []byte {
    return append(WithdrawPrefix, marshal(seq)...)
}

// ------------i-------------------

func loadLastHeader(store state.SimpleDB) (res uint64) {
    b := store.Get(LastHeaderKey)
    if b == nil {
        return 0 // change it to (pseudo) Genesis header number
    }
    return unmarshal(b).(uint64)
}

func loadInfo(store state.SimpleDB, chainid string) (res Info) {
    b := store.Get(GetInfoKey(chainid))
    if b == nil {
        panic("Info not found")
    }
    return unmarshal(b).(Info)
}

func saveInfo(store state.SimpleDB, chainid string, info Info) {
    store.Set(GetInfoKey(chainid), marshal(info))
}

func loadSignedHeaders(store state.SimpleDB, height uint64) (res []SignedHeader) {
    b := store.Get(GetBufferKey(height))
    if b == nil {
        return nil
    }
    return unmarshal(b).([]SignedHeader)
}

func saveSignedHeaders(store state.SimpleDB, height uint64, signedHeaders []SignedHeader) {
    store.Set(GetBufferKey(height), marshal(signedHeaders))
}

func removeSignedHeaders(store state.SimpleDB, height uint64) {
    store.Remove(GetBufferKey(height))
}

func loadValidators(store state.SimpleDB) (res []sdk.Actor) {
    b := store.Get(ValidatorsKey)
    if b == nil {
        return nil
    }
    return unmarshal(b).([]sdk.Actor)
}

func saveFinalized(store state.SimpleDB, height uint64, header Header) {
    store.Set(GetFinalKey(height), marshal(header))
}

func loadFinalized(store state.SimpleDB, height uint64) (res Header, exists bool) {
    b := store.Get(GetFinalKey(height))
    if b == nil {
        return res, false
    }
    return unmarshal(b).(Header), true
}

func loadSigners(store state.SimpleDB, height uint64, hash common.Hash) []sdk.Actor {
    signedHeaders := loadSignedHeaders(store, height)
    for _, s := range signedHeaders {
        if s.Header.Hash == hash {
            return s.Signers
        }
    }
    return nil
}

func saveBalance(store state.SimpleDB, chainid string, token common.Address, balance *big.Int) {
    store.Set(GetBalanceKey(chainid, token), marshal(balance))
}

func loadBalance(store state.SimpleDB, chainid string, token common.Address) (res *big.Int) {
    b := store.Get(GetBalanceKey(chainid, token))
    if b == nil {
        return new(big.Int).SetUint64(0)
    }
    return new(big.Int).SetBytes(unmarshal(b).([]byte))
}

func saveWithdraw(store state.SimpleDB, to common.Address, value *big.Int, token common.Address, seq uint64) {
    data := [][]byte {
        to.Bytes(),
        value.Bytes(),
        token.Bytes(),
    }
    
    var encoded []byte
    for _, d := range data {
        encoded = append(encoded, d...)
    }

    store.Set(GetWithdrawKey(seq), encoded)
}

func increaseBalance(store state.SimpleDB, chainid string, token common.Address, diff *big.Int) {
    saveBalance(store, chainid, token, new(big.Int).Add(loadBalance(store, chainid, token), diff))
}

func decreaseBalance(store state.SimpleDB, chainid string, token common.Address, diff *big.Int) {
    saveBalance(store, chainid, token, new(big.Int).Sub(loadBalance(store, chainid, token), diff))
}
