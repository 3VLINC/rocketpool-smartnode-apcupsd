package account

import (
    "github.com/urfave/cli"

    cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)


// Register subcommands
func RegisterSubcommands(command *cli.Command, name string, aliases []string) {
    command.Subcommands = append(command.Subcommands, cli.Command{
        Name:      name,
        Aliases:   aliases,
        Usage:     "Manage the node account",
        Subcommands: []cli.Command{

            cli.Command{
                Name:      "status",
                Aliases:   []string{"s"},
                Usage:     "Get the node password and account status",
                UsageText: "rocketpool api account status",
                Action: func(c *cli.Context) error {

                    // Validate args
                    if err := cliutils.CheckAPIArgCount(c, 0); err != nil { return err }

                    // Run
                    return getStatus(c)

                },
            },

            cli.Command{
                Name:      "init-password",
                Aliases:   []string{"p"},
                Usage:     "Initialize the node password",
                UsageText: "rocketpool api account init-password password",
                Action: func(c *cli.Context) error {

                    // Validate args
                    if err := cliutils.CheckAPIArgCount(c, 1); err != nil { return err }
                    password := c.Args().Get(0)

                    // Run
                    return initPassword(c, password)

                },
            },

            cli.Command{
                Name:      "init-account",
                Aliases:   []string{"a"},
                Usage:     "Initialize the node account",
                UsageText: "rocketpool api account init-account",
                Action: func(c *cli.Context) error {

                    // Validate args
                    if err := cliutils.CheckAPIArgCount(c, 0); err != nil { return err }

                    // Run
                    return initAccount(c)

                },
            },

            cli.Command{
                Name:      "export",
                Aliases:   []string{"e"},
                Usage:     "Export the node account in JSON format",
                UsageText: "rocketpool api account export",
                Action: func(c *cli.Context) error {

                    // Validate args
                    if err := cliutils.CheckAPIArgCount(c, 0); err != nil { return err }

                    // Run
                    return exportAccount(c)

                },
            },

        },
    })
}
