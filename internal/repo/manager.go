package repo

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/ezedh/bootcamps/pkg/http"
)

type createReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type RepoManager interface {
	// SetName sets the name of the repo
	SetName(name string)
	// CreateRepo creates a new repository given a description.
	CreateRepo(desc string) error
	// PushChanges pushes changes to the repository given a commit message.
	PushChanges(message string) error
}

type repoManager struct {
	token      string
	username   string
	name       string
	apiManager http.ApiManager
}

func NewRepoManager(token, username string) RepoManager {
	apiManager := http.NewApiManager(token)
	return &repoManager{
		token:      token,
		username:   username,
		name:       "",
		apiManager: apiManager,
	}
}

func (r *repoManager) SetName(name string) {
	r.name = name
}

func (r *repoManager) CreateRepo(desc string) error {
	creq := &createReq{
		Name:        r.name,
		Description: desc,
		Private:     true,
	}

	postBody, err := json.Marshal(creq)
	if err != nil {
		return err
	}

	url := "/user/repos"

	err = r.apiManager.Post(url, postBody, nil)

	if err != nil {
		return err
	}

	return r.initializeRepo(desc)
}

func (r *repoManager) PushChanges(message string) error {
	fmt.Println("Agregando cambios...")
	err := r.execRepoGitCommand("add", ".")
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar el readme")
	}

	fmt.Println("Haciendo commits de cambios...")
	err = r.execRepoGitCommand("commit", "-m", message)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("ocurrió un error al hacer commit de los cambios: %s", err))
	}

	fmt.Println("Subiendo cambios...")
	err = r.execRepoGitCommand("push", "-u", "origin", "main")
	if err != nil {
		return fmt.Errorf("ocurrió un error al pushear a main")
	}

	return nil
}

func (r *repoManager) initializeRepo(desc string) error {
	fmt.Println("Clonando repo...")

	repo := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s/%s.git", r.token, r.username, r.name)

	// Clonar usando git clone repo shell
	// git clone
	err := exec.Command("git", "clone", repo).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al clonar el repositorio")
	}

	fmt.Println("Creando readme...")
	// Create file README.md inside the repo folder
	err = exec.Command("touch", fmt.Sprintf("./%s/README.md", r.name)).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al crear el readme")
	}

	fmt.Println("Inicializando repo...")
	err = r.execRepoGitCommand("init", ".")
	if err != nil {
		return fmt.Errorf("ocurrió un error al incializar el repositorio")
	}

	fmt.Println("Agregando readme...")
	err = r.execRepoGitCommand("add", ".")
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar el readme")
	}

	fmt.Println("Haciendo commits de cambios...")
	err = r.execRepoGitCommand("commit", "-m", "initial commit")
	if err != nil {
		return fmt.Errorf("ocurrió un error al hacer commit de los cambios")
	}

	fmt.Println("Creando branch main...")
	err = r.execRepoGitCommand("branch", "-M", "main")
	if err != nil {
		return fmt.Errorf("ocurrió un error al crear branch main")
	}

	fmt.Println("Configurando origin...")
	_ = r.execRepoGitCommand("remote", "add", "origin", repo)

	err = r.execRepoGitCommand("push", "-u", "origin", "main")
	if err != nil {
		return fmt.Errorf("ocurrió un error al pushear a main")
	}

	return nil
}

// execRepoGitCommand executes a git command in the repo
func (r *repoManager) execRepoGitCommand(args ...string) error {
	command := append([]string{"-C", r.name}, args...)
	return exec.Command("git", command...).Run()
}
