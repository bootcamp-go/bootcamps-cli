package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ezedh/bootcamps/internal/repo"
	"github.com/ezedh/bootcamps/pkg/color"
)

const (
	TemplateRepo = "https://github.com/ezedh/bootcamps-templates.git"
)

type TemplateManager interface {
	SetName(name string)
	PlaceTemplateInRepo() error
	ReplaceImportsInRepo() error
	RemoveRepoFolder()
	Clean()
}

type templateManager struct {
	name     string
	username string
	path     string
	company  string
	folder   string
	repoM    repo.RepoManager
}

func NewTemplateManager(username, company, uuid string, repoM repo.RepoManager) TemplateManager {
	folder := fmt.Sprintf("%s-%s", uuid, "template")
	return &templateManager{
		username: username,
		company:  company,
		folder:   folder,
		repoM:    repoM,
	}
}

func (tm *templateManager) SetName(name string) {
	tm.name = name
	tm.path = "./" + name
}

func (tm *templateManager) PlaceTemplateInRepo() error {
	fmt.Println("Buscando template...")
	if err := tm.findTemplateFolder(); err != nil {
		color.Print("red", err.Error())
		return err
	}

	// copy template folder content into repo folder
	fmt.Println("Copiando template...")
	cmd := fmt.Sprintf("cp -r ./%s/* %s", tm.folder, tm.path)
	cmdgithub := fmt.Sprintf("cp -r ./%s/.github %s", tm.folder, tm.path)
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
	fmt.Println("Reemplazando imports...")
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

func (tm *templateManager) Clean() {
	os.RemoveAll(tm.folder)
}

// findTemplateFolder finds the template folder in the current directory
func (tm *templateManager) findTemplateFolder() error {
	// check if a "template" folder exists
	// if not, create one

	if _, err := os.Stat(fmt.Sprintf("./%s", tm.folder)); os.IsNotExist(err) {
		// clone template folder from github TemplateRepo from meli branch
		tm.repoM.SetName("bootcamps-templates")
		return tm.repoM.CloneFromBranchDH(tm.folder, tm.company)
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
