package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ezedh/bootcamps/pkg/color"
)

type TemplateManager interface {
	PlaceTemplateInRepo() error
	ReplaceImportsInRepo() error
	RemoveRepoFolder()
}

type templateManager struct {
	name     string
	username string
	path     string
}

func NewTemplateManager(name string, username string) TemplateManager {
	return &templateManager{
		name:     name,
		username: username,
		path:     "./" + name,
	}
}

func (tm *templateManager) PlaceTemplateInRepo() error {
	fmt.Println("Buscando template...")
	if err := findTemplateFolder(); err != nil {
		color.Print("red", err.Error())
		return err
	}

	// copy template folder content into repo folder
	fmt.Println("Copiando template...")
	cmd := fmt.Sprintf("cp -r ./template/* %s", tm.path)
	cmdgithub := fmt.Sprintf("cp -r ./template/.github %s", tm.path)
	err := exec.Command("bash", "-c", cmd).Run()
	if err != nil {
		color.Print("red", fmt.Sprintf("Couldn't copy template folder: %s", err.Error()))
		return err
	}

	err = exec.Command("bash", "-c", cmdgithub).Run()
	if err != nil {
		color.Print("red", fmt.Sprintf("Couldn't copy .github template folder: %s", err.Error()))
		return err
	}
	return nil
}

func (tm *templateManager) ReplaceImportsInRepo() error {
	err := filepath.Walk(tm.path, tm.visit)
	if err != nil {
		return err
	}

	return nil
}

func (tm *templateManager) RemoveRepoFolder() {
	fmt.Println("Limpiando repo...")
	err := os.RemoveAll(tm.path)
	if err != nil {
		color.Print("red", fmt.Sprintf("Couldn't remove repo folder: %s", err.Error()))
	}
}

// findTemplateFolder finds the template folder in the current directory
func findTemplateFolder() error {
	// check if a "template" folder exists
	// if not, create one

	if _, err := os.Stat("./template"); os.IsNotExist(err) {
		// clone template folder from github repo https://github.com/ezedh/bootcamp-template.git
		fmt.Println("No se encontró el template, clonando desde github")
		cmd := exec.Command("git", "clone", "https://github.com/ezedh/bootcamp-template.git", "template")
		err := cmd.Run()
		if err != nil {
			color.Print("red", fmt.Sprintf("Couldn't clone template folder: %s", err.Error()))
			return err
		}
	}

	return nil
}

func (tm *templateManager) visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil //
	}

	matched, _ := filepath.Match("*go*", fi.Name())

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), "usuario/repositorio", tm.username+"/"+tm.name, -1)

		err = os.WriteFile(path, []byte(newContents), 0644)
		if err != nil {
			return err
		}

	}

	return nil
}
