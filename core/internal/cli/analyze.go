package cli

import "github.com/spf13/cobra"

var analyzeCmd = &cobra.Command{
	Use:   "analyze <job-id>",
	Short: "Analizar una vacante contra tu perfil",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementar analisis
		return nil
	},
}
