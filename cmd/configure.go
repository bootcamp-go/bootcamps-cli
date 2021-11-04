package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ezedh/bootcamps/pkg/color"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configurar el CLI con el token y usuario de github",
	Long:  `Configurar el CLI con el token y usuario de github`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configurar el token de DH
		fmt.Print("Introduzca el token de DH: ")
		var tokenDH string
		fmt.Scanln(&tokenDH)
		viper.Set("tokenDH", tokenDH)
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Token de DH guardado correctamente")

		// Configurar el token
		fmt.Print("Introduzca el token de github: ")
		var token string
		fmt.Scanln(&token)
		viper.Set("token", token)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Token actualizado")

		// Configurar el usuario
		fmt.Print("Introduzca el usuario de github: ")
		var user string
		fmt.Scanln(&user)
		viper.Set("username", user)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Usuario actualizado")

		// Configurar la empresa
		fmt.Print("Introduzca el nombre de la empresa (meli): ")
		var company string
		fmt.Scanln(&company)
		viper.Set("company", company)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Empresa actualizada")
		return nil
	},
}

var configureTokenDH = &cobra.Command{
	Use:   "dh",
	Short: "Configurar el CLI con el token de DH",
	Long:  `Configurar el CLI con el token de DH`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configurar el token
		fmt.Print("Introduzca el token de DH: ")
		var token string
		fmt.Scanln(&token)
		viper.Set("tokendh", token)
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Token de DH actualizado")
		return nil
	},
}

var configureToken = &cobra.Command{
	Use:   "token",
	Short: "Configurar el CLI con el token",
	Long:  `Configurar el CLI con el token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configurar el token
		fmt.Print("Introduzca el token de github: ")
		var token string
		fmt.Scanln(&token)
		viper.Set("token", token)
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Token actualizado")
		return nil
	},
}

var configureUsername = &cobra.Command{
	Use:   "username",
	Short: "Configurar el usuario de github",
	Long:  `Configurar el usuario de github`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configurar el usuario
		fmt.Print("Introduzca el usuario de github: ")
		var user string
		fmt.Scanln(&user)
		viper.Set("username", user)
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Usuario actualizado")
		return nil
	},
}

var configureCompany = &cobra.Command{
	Use:   "company",
	Short: "Configurar la empresa",
	Long:  `Configurar la empresa`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configurar la empresa
		fmt.Print("Introduzca el nombre de la empresa: ")
		var company string
		fmt.Scanln(&company)
		viper.Set("company", company)
		err := viper.WriteConfig()
		if err != nil {
			return err
		}
		color.Print("green", "Empresa actualizada")
		return nil
	},
}

func init() {
	// bootcamps configure
	rootCmd.AddCommand(configureCmd)
	// bootcamps configure dh
	configureCmd.AddCommand(configureTokenDH)
	// bootcamps configure token
	configureCmd.AddCommand(configureToken)
	// bootcamps configure username
	configureCmd.AddCommand(configureUsername)
	// bootcamps configure company
	configureCmd.AddCommand(configureCompany)
}
