package cmd

import (
	"fmt"
	"os/exec"

	"github.com/ezedh/bootcamps/internal/config"
	"github.com/ezedh/bootcamps/internal/repo"
	"github.com/ezedh/bootcamps/internal/test"
	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetConfiguration()
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al obtener la configuración: %s", err))
			return
		}

		repoName := fmt.Sprintf(config.GoRepoNameFormat, c.Company, wave, group)
		repoM := repo.NewRepoManager(c.Token, owner)

		sprintTag := fmt.Sprintf("sprint_%s.0.0", sprint)
		repoFolder := fmt.Sprintf("/tmp/%s", repoName)

		err = exec.Command("rm", "-rf", repoFolder).Run()
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al borrar el repositorio: %s", err))
			return
		}

		repoM.SetName(repoName)

		err = repoM.CloneFromBranch(repoFolder, sprintTag)
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al obtener el sprint: %s", err))
			return
		}

		tester := test.NewTester(uuid.New().String(), c.Company, repoFolder, repoM)
		err = tester.Test(wave, sprint)
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al testear el sprint: %s", err))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&wave, "wave", "w", "", "Wave of the sprint")
	testCmd.Flags().StringVarP(&group, "group", "g", "", "Group of the sprint")
	testCmd.Flags().StringVarP(&sprint, "sprint", "s", "", "Sprint of the sprint")
	getCmd.Flags().StringVarP(&owner, "owner", "o", "", "Owner of the sprint's repo")

	testCmd.MarkFlagRequired("wave")
	testCmd.MarkFlagRequired("group")
	testCmd.MarkFlagRequired("sprint")
	testCmd.MarkFlagRequired("owner")
}
