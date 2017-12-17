package gate

import (
    "math/big"
    "io/ioutil"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var GateInitCmd = &cobra.Command {
    Use: "init",
    Short: "Register Ethereum contract",
    RunE: gateInitCmd,
}

func gateInitCmd(cmd *cobra.Command, args []string) error {
    g, err := newGateway()
    if err != nil {
        return err
    }

    g.ethauth.GasLimit = big.NewInt(4700000)

    genesisBytes, err := ioutil.ReadFile(viper.GetString(FlagGenesis))
}
