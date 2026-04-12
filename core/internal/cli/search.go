package cli

import "github.com/spf13/cobra"

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Buscar vacantes en portales de empleo",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implementar busqueda
		return nil
	},
}

func init() {
	searchCmd.Flags().String("role", "", "Puesto o area profesional")
	searchCmd.Flags().Int("salary-min", 0, "Salario minimo mensual")
	searchCmd.Flags().String("modality", "", "Modalidad: remote | hybrid | onsite")
	searchCmd.Flags().String("location", "", "Ciudad o municipio base")
}
