package cmd

import (
	"fmt"
	"os/exec"

	"github.com/ezedh/bootcamps/internal/config"
	"github.com/ezedh/bootcamps/internal/repo"
	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/spf13/cobra"
)

var (
	wave   string
	group  string
	sprint string
	owner  string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get -w [WAVE] -g [GROUP] -s [SPRINT] -o [REPO OWNER]",
	Short: "Get the sprint repository",
	Long:  `Get the sprint repository of a group of a wabe`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetConfiguration()
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al obtener la configuración: %s", err))
			return
		}

		repoName := fmt.Sprintf(config.GoRepoNameFormat, c.Company, wave, group)
		repoM := repo.NewRepoManager(c.Token, owner)

		_ = exec.Command("rm", "-rf", repoName).Run()

		sprintTag := fmt.Sprintf("sprint_%s.0.0", sprint)

		err = exec.Command("rm", "-rf", repoName).Run()
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al borrar el repositorio: %s", err))
			return
		}

		repoM.SetName(repoName)

		err = repoM.CloneFromBranch(repoName, sprintTag)
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al obtener el sprint: %s", err))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&wave, "wave", "w", "", "Wave of the sprint")
	getCmd.Flags().StringVarP(&group, "group", "g", "", "Group of the sprint")
	getCmd.Flags().StringVarP(&sprint, "sprint", "s", "", "Sprint of the sprint")
	getCmd.Flags().StringVarP(&owner, "owner", "o", "", "Owner of the sprint's repo")

	getCmd.MarkFlagRequired("wave")
	getCmd.MarkFlagRequired("group")
	getCmd.MarkFlagRequired("sprint")
	getCmd.MarkFlagRequired("owner")
}
