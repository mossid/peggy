package gate

import (
    "github.com/spf13/cobra"   



//    basecmd "github.com/cosmos/cosmos-sdk/server/commands"
)

var GateCmd = &cobra.Command {
    Use: "gate",
    Short: "Relay Ethereum logs to Tendermint",
}

var (
    depositABI abi.ABI
)

const (
    FlagTestnet  = "testnet"
    FlagDatadir  = "datadir"
    FlagGenesis  = "genesis"
    FlagIpcpath  = "ipcpath"
    FlagNodeaddr = "nodeaddr"
    FlagAddress  = "address"
)

func init() {
    GateInitCmd.Flags().String(FlagGenesis, "", "Path to genesis file")

    GateStartCmd.Flags().Bool(FlagTestnet, false, "Ropsten network: pre-configured test network")

    GateCmd.AddCommand(GateStartCmd)
    GateCmd.AddCommand(GateInitCmd)
    GateCmd.AddCommand(GateGenValidatorCmd)
}



