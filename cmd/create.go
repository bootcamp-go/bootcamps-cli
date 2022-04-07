package cmd

import (
	"context"
	"fmt"

	"github.com/ezedh/bootcamps/internal/config"
	"github.com/ezedh/bootcamps/internal/invitation"
	"github.com/ezedh/bootcamps/internal/repo"
	"github.com/ezedh/bootcamps/internal/secrets"
	"github.com/ezedh/bootcamps/internal/template"
	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/ezedh/bootcamps/pkg/confirm"
	"github.com/ezedh/bootcamps/pkg/ghclient"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Crear repositorios semilla para el bootcamp",
	Long: `Crear repositorios semilla para el bootcamp indicando la wave y la
	cantidad de repos..`,
	Run: func(cmd *cobra.Command, args []string) {
		uuid := uuid.New().String()

		c, err := config.GetConfiguration()
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al obtener la configuración: %s", err))
			return
		}

		ghclient := ghclient.NewGhClient(context.Background(), c.Token)
		repoM := repo.NewRepoManager(c.Token, c.Username)
		templateM := template.NewTemplateManager(c.Username, c.Company, uuid, repoM)
		defer templateM.Clean()

		inviter := invitation.NewInviter(repoM, c.Company, uuid)

		inv, err := inviter.GetCreationConfig()
		if err != nil {
			color.Print("red", err.Error())
			return
		}

		if !confirm.Ask("¿Desea crear " + fmt.Sprintf("%d", inv.Amount) + " repositorios para la wave " + inviter.Wave() + " de " + inviter.Company() + "?") {
			return
		}

		for i := 1; i <= inv.Amount; i++ {
			// index to string
			index := fmt.Sprintf("%d", i)
			repoName := fmt.Sprintf(config.GoRepoNameFormat, c.Company, inviter.Wave(), index)
			secretsM := secrets.NewSecretsManager(ghclient, c.Username, repoName)
			color.Print("cyan", fmt.Sprintf("Crear repositorio %s", repoName))

			repoM.SetName(repoName)

			err := repoM.CreateRepo(repoName)
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al crear el repositorio: %s", err.Error()))
				return
			}

			color.Print("cyan", "Invitar usuarios")
			// i to string to get the correct index
			strI := fmt.Sprintf("%d", i)
			err = repoM.InviteUsers(inv.Groups[strI])
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al invitar usuarios: %s", err.Error()))
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

			repoM.SetName(repoName)

			err = repoM.PushChanges("add template")
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al subir los cambios: %s", err.Error()))
				return
			}

			templateM.RemoveRepoFolder()

			color.Print("green", "Repositorio creado y configurado con éxito")

			color.Print("cyan", "Crear secrets")
			err = secretsM.SetSecret(context.Background(), "API_URL", c.ApiUrl)
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al crear el secret API_URL: %s", err.Error()))
				return
			}
			err = secretsM.SetSecret(context.Background(), "API_KEY", c.ApiKey)
			if err != nil {
				color.Print("red", fmt.Sprintf("Error al crear el secret API_KEY: %s", err.Error()))
				return
			}

			fmt.Printf("\n\n")
		}

		color.Print("green", "Todos los repositorios fueron creados y configurados con éxito")
	},
}

func init() {
	// bootcamps create
	rootCmd.AddCommand(createCmd)
}
