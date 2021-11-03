package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ezedh/bootcamps/internal/config"
	"github.com/ezedh/bootcamps/internal/repo"
	"github.com/ezedh/bootcamps/internal/template"
	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Crear repositorios semilla para el bootcamp",
	Long: `Crear repositorios semilla para el bootcamp indicando la wave y la
	cantidad de repos..`,
	Run: func(cmd *cobra.Command, args []string) {
		wave, ammount, err := getCreationConfig()
		if err != nil {
			color.Print("red", err.Error())
			return
		}

		fmt.Printf("\n\n")

		c, err := config.GetConfiguration()
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al obtener la configuración: %s", err))
			return
		}

		repoM := repo.NewRepoManager(c.Token, c.Username)
		templateM := template.NewTemplateManager(c.Username, c.Company)

		for i := 1; i <= ammount; i++ {
			repoName := fmt.Sprintf("%s_bootcamp_w%s-%d", c.Company, wave, i)

			color.Print("cyan", fmt.Sprintf("Crear repositorio %s", repoName))

			repoM.SetName(repoName)

			err := repoM.CreateRepo(repoName)
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al crear el repositorio: %s", err.Error()))
				return
			}

			color.Print("cyan", "Colocar el template")

			templateM.SetName(repoName)

			err = templateM.PlaceTemplateInRepo()
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al crear el template: %s", err.Error()))
				return
			}

			err = templateM.ReplaceImportsInRepo()
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al reemplazar los imports: %s", err.Error()))
				return
			}

			err = repoM.PushChanges("add template")
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al subir los cambios: %s", err.Error()))
				return
			}

			templateM.RemoveRepoFolder()

			color.Print("green", "Repositorio creado y configurado con éxito")

			fmt.Printf("\n\n")
		}

		_ = os.RemoveAll("./template")

		color.Print("green", "Todos los repositorios fueron creados y configurados con éxito")
	},
}

func getCreationConfig() (string, int, error) {
	var wave string
	var ammount int

	fmt.Printf("Wave N°: ")
	fmt.Scan(&wave)

	fmt.Printf("Cantidad de grupos: ")
	fmt.Scan(&ammount)

	fmt.Printf("\n\n")

	fmt.Printf("Está a punto de crear %d repositorios para la wave %s, está de acuerdo? (y/N): ", ammount, wave)

	var answer string
	fmt.Scan(&answer)

	if strings.ToLower(answer) != "y" {
		return "", 0, errors.New("Cancelado")
	}

	return wave, ammount, nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}
