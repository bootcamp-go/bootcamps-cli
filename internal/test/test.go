package test

import (
	"fmt"
	"os/exec"

	"github.com/ezedh/bootcamps/internal/repo"
)

type Tester interface {
	Test(wave string, sprint string) error
}

type tester struct {
	repoM        repo.RepoManager
	folder       string
	company      string
	sprintFolder string
}

func NewTester(uuid, company, sprintFolder string, repoM repo.RepoManager) Tester {
	folder := fmt.Sprintf("/tmp/%s-test", uuid)
	return &tester{
		repoM:        repoM,
		folder:       folder,
		company:      company,
		sprintFolder: sprintFolder,
	}
}

func (t *tester) Test(wave string, sprint string) error {
	err := t.clone()
	if err != nil {
		return err
	}

	fmt.Println("Restaurando base de datos")

	var sqlPassword string
	fmt.Print("Contraseña de mysql: ")
	fmt.Scan(&sqlPassword)

	// use command line tool to restore database melisprint
	passString := fmt.Sprintf("p=%s", sqlPassword)
	err = exec.Command("make", "-C", t.sprintFolder, "rebuild-database-with-password", passString).Run()
	if err != nil {
		return fmt.Errorf("error al restaurar la base de datos")
	}

	fmt.Println("Ejecutando tests")
	waveF := fmt.Sprintf("wave=%s", wave)
	sprintF := fmt.Sprintf("sprint=%s", sprint)
	folderF := fmt.Sprintf("folder=%s", t.sprintFolder)
	fmt.Println(exec.Command("make", "-C", t.folder, "test", waveF, sprintF, folderF).String())
	o, err := exec.Command("make", "-C", t.folder, "test", waveF, sprintF, folderF).Output()
	if err != nil {
		return fmt.Errorf("error al ejecutar los tests: %s", err)
	}

	fmt.Println(string(o))
	return nil
}

func (t *tester) clone() error {
	t.repoM.SetName("bootcamps-tests")
	return t.repoM.CloneFromBranchDH(t.folder, t.company)
}
