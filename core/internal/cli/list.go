package cli

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar vacantes recopiladas",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementar listado
		return nil
	},
}
