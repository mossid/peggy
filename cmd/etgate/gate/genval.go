package gate

import (
    "fmt"
    "encoding/hex"
    "strings"
    
    "github.com/spf13/cobra"

    "github.com/ethereum/go-ethereum/common"

    "github.com/tendermint/go-crypto"

)

var GateGenValidatorCmd = &cobra.Command {
    Use: "genval",
    Short: "Generate secp256k1 validator",
    RunE: gateGenValidatorCmd,
}

func gateGenValidatorCmd(cmd *cobra.Command, args []string) error {
    privKey := crypto.GenPrivKeySecp256k1()
    pubKey := privKey.Wrap().PubKey().Unwrap()
    var addr common.Address
    copy(addr[:], pubKey.Address())
    fmt.Printf("Priv:\t%v\nPub:\t%v\nAddr:\t%v\n", strings.ToUpper(hex.EncodeToString(privKey[:])), pubKey.KeyString(), strings.ToUpper(addr.Hex()[2:]))

    return nil
}
