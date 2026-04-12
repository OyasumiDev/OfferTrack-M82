package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "offertrack",
	Short: "OfferTrack M82 - Plataforma de busqueda de empleo asistida por IA",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(adaptCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(configCmd)
}
