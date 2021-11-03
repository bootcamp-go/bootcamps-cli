/*
Copyright © 2021 Ezekiel Grosfeld grosfeldezekiel@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ezedh/bootcamps/internal/repo"
	"github.com/ezedh/bootcamps/internal/template"
	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Crear repositorios semilla para el bootcamp",
	Long: `Crear repositorios semilla para el bootcamp indicando la wave y la
	cantidad de repos..`,
	Run: func(cmd *cobra.Command, args []string) {
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
			color.Print("red", "Cancelado")
			return
		}

		fmt.Printf("\n\n")

		token := viper.GetString("token")
		if token == "" {
			color.Print("red", "No se pudo obtener el token")
			return
		}

		username := viper.GetString("username")
		if username == "" {
			color.Print("red", "No se pudo obtener el username")
			return
		}

		repoM := repo.NewRepoManager(token, username)

		for i := 1; i <= ammount; i++ {
			repoName := fmt.Sprintf("meli_bootcamp_w%s-%d", wave, i)

			color.Print("cyan", fmt.Sprintf("Crear repositorio %s", repoName))

			repoM.SetName(repoName)

			err := repoM.CreateRepo(repoName)
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al crear el repositorio: %s", err.Error()))
				return
			}

			templateM := template.NewTemplateManager(repoName, username)

			color.Print("cyan", "Colocar el template")

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

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
