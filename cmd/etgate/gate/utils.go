package gate

import (
    "fmt"
    "os"
    "errors"
    "path/filepath"
    "context"
 
    "github.com/spf13/viper"

    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    ecrypto "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/common"

    mintclient "github.com/tendermint/tendermint/rpc/client"

    "github.com/tendermint/go-crypto"

    basecmd "github.com/cosmos/cosmos-sdk/server/commands"
    querycmd "github.com/cosmos/cosmos-sdk/client/commands/query"

    "../../../contracts"
    "../../../modules/etgate"
)

type gateway struct {
    ethclient *ethclient.Client
    mintclient *mintclient.HTTP
    ethauth *bind.TransactOpts
    mintkey *basecmd.Key
}

func newGateway() (*gateway, error) {
    var datadir string
    if viper.GetBool(FlagTestnet) {
        datadir = filepath.Join(viper.GetString(FlagDatadir), "testnet")
    } else {
        datadir = viper.GetString(FlagDatadir)
    }

    clientpath := filepath.Join(datadir, viper.GetString(FlagIpcpath))

    ethclient, err := ethclient.Dial(clientpath)
    if err != nil {
        return nil, err
    }

    mintkey, err := basecmd.LoadKey(filepath.Join(os.Getenv("HOME"), ".etgate", "server", "key.json"))
    if err != nil {
        return nil, err
    }

    priv, err := getSecp256k1Priv(mintkey.PrivKey)
    if err != nil {
        return nil, err
    }

    ecdsa, err := ecrypto.ToECDSA(priv[:])
    if err != nil {
        return nil, err
    }

    fmt.Printf("Using Ethereum address %+v\n", ecrypto.PubkeyToAddress(ecdsa.PublicKey).Hex())

    ethauth := bind.NewKeyedTransactor(ecdsa)

    return &gateway{
        ethclient: ethclient,
        mintclient: mintclient.NewHTTP(viper.GetString(FlagNodeaddr), "/websocket"),
        ethauth: ethauth,
        mintkey: mintkey,
    }, nil
}

func (g *gateway) start() {
    contract, err := contracts.NewETGate(common.HexToAddress(viper.GetString(FlagAddress)), g.ethclient)
    if err != nil {
        panic(err)
    }

    go g.depositRelay()
    go g.withdrawRelay(contract)
}

func (g *gateway) depositRelay() { // eth->mint
    querycmd.Get(etgate.Get)
}

func (g *gateway) withdrawRelay(c *contracts.ETGate) { // mint->eth also with header when necessary

}

func (g *gateway) headerRelay() { // eth->mint onlyvalidators
    heads := make(chan *types.Header)
    headsub, err := g.ethclient.SubscribeNewHead(context.Background(), heads)
    if err != nil {
        panic("Failed to subscribe to new headers")
    }

    defer headsub.Unsubscribe()
    
    for {
        select {
        case head := <-heads:
            
        }
    }
}

func (g *gateway) changeRelay() { // eth->mint

}
 
func getSecp256k1Priv(priv crypto.PrivKey) (crypto.PrivKeySecp256k1, error) {
    switch inner := priv.Unwrap().(type) {
    case crypto.PrivKeySecp256k1:
        return inner, nil
    default:
        return crypto.PrivKeySecp256k1{}, errors.New("PrivKey is not secp256k1")
    }
}

func getSecp256k1Pub(pub crypto.PubKey) (crypto.PubKeySecp256k1, error) {
    switch inner := pub.Unwrap().(type) {
    case crypto.PubKeySecp256k1:
        return inner, nil
    default:
        return crypto.PubKeySecp256k1{}, errors.New("PubKey is not secp256k1")
    }
}


