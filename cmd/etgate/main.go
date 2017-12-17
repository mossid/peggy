package main

import (
    "os"

    "github.com/spf13/cobra"

    "github.com/tendermint/tmlibs/cli"

    scmd "github.com/cosmos/cosmos-sdk/server/commands"
//    "github.com/tendermint/basecoin/types"


    "../../modules/etgate"
    "./commands"
)

func main() {
    var RootCmd = &cobra.Command {
        Use: "etgate",
        Short: "ethereum log relaying plugin for basecoin",
    }

    RootCmd.AddCommand(
        scmd.InitCmd,
        scmd.StartCmd,
        scmd.RelayCmd,
        GateCmd,
        scmd.UnsafeResetAllCmd,
        scmd.VersionCmd,
    )
    
    basecmd.RegisterStartPlugin("ETGATE", func() types.Plugin { return etgate.New() })

    cmd := cli.PrepareMainCmd(RootCmd, "ETGATE", os.ExpandEnv("$HOME/.etgate/server"))
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}

