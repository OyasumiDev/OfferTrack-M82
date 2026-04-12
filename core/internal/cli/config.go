package cli

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gestionar configuracion del usuario y proveedor IA",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementar configuracion
		return nil
	},
}
