package cli

import "github.com/spf13/cobra"

var adaptCmd = &cobra.Command{
	Use:   "adapt-cv <job-id>",
	Short: "Generar CV adaptado para una vacante",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementar adaptacion de CV
		return nil
	},
}

func init() {
	adaptCmd.Flags().String("output", "./exports", "Directorio de salida")
}
