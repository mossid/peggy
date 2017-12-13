package etgate

import (
    "bytes"

    sdk "github.com/cosmos/cosmos-sdk"
    "github.com/cosmos/cosmos-sdk/state"

    wire "github.com/tendermint/go-wire"

    "github.com/ethereum/go-ethereum/common"
)

var (
    InfoKey       = []byte{0x01} // Info
    ValidatorsKey = []byte{0x02} // []sdk.Actor

    BufferPrefix  = []byte{0x03} // uint => []SignedHeader
    FinalPrefix   = []byte{0x04} // uint => Header
    BalancePrefix = []byte{0x05} // string => uint
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

func GetBufferKey(height uint64) []byte { 
    return append(BufferPrefix, marshal(height)...)
}

func GetFinalKey(height uint64) []byte { 
    return append(FinalPrefix, marshal(height)...)
}

func GetBalanceKey(chainid string) []byte {
    return append(BalancePrefix, chainid...)
}

// ------------i-------------------

func loadInfo(store state.SimpleDB) (res Info) {
    b := store.Get(InfoKey)
    if b == nil {
        panic("Info not found")
    }
    return unmarshal(b).(Info)
}

func saveInfo(store state.SimpleDB, info Info) {
    store.Set(InfoKey, marshal(info))
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

func saveBalance(store state.SimpleDB, chainid string, balance uint64) {
    store.Set(GetBalanceKey(chainid), marshal(balance))
}

func loadBalance(store state.SimpleDB, chainid string) (res uint64) {
    b := store.Get(GetBalanceKey(chainid))
    if b == nil {
        return 0
    }
    return unmarshal(b).(uint64)
}

func increaseBalance(store state.SimpleDB, chainid string, diff uint64) {
    saveBalance(store, chainid, loadBalance(store, chainid) + diff)
}

func decreaseBalance(store state.SimpleDB, chainid string, diff uint64) {
    saveBalance(store, chainid, loadBalance(store, chainid) - diff)
}
