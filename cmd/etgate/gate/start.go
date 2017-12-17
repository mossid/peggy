package gate

import (
    "github.com/spf13/cobra"
)

var GateStartCmd = &cobra.Command {
    Use: "start",
    Short: "Start etgate relayer",
    RunE: gateStartCmd,
}


func gateStartCmd(cmd *cobra.Command, args []string) error {
    gateway, err := newGateway()
    if err != nil {
        return err
    }

    gateway.start()
}
