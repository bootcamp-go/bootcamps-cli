package repo

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/ezedh/bootcamps/pkg/color"
	"github.com/ezedh/bootcamps/pkg/http"
)

type createReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type inviteReq struct {
	Permission string `json:"permission"`
}

type RepoManager interface {
	// SetName sets the name of the repo
	SetName(name string)
	// Clone clones the repository with the given name in the current directory.
	Clone(name string) error
	// CloneFromBranch clones the repository with the given name in the current directory from the given branch.
	CloneFromBranch(name, branch string) error
	// Clone clones the repository with the given name in the current directory.
	CloneDH(name string) error
	// CloneFromBranch clones the repository with the given name in the current directory from the given branch.
	CloneFromBranchDH(name, branch string) error
	// PushChanges pushes changes to the repository given a commit message.
	PushChanges(message string) error
	// CreateRepo creates a new repository given a description.
	CreateRepo(desc string) error
	// InviteUsers invites the given users to the repository.
	InviteUsers(users []string) error
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

func (r *repoManager) Clone(name string) error {
	fmt.Println("Clonando repo...")

	repo := r.getRepoURL()

	if name == "" {
		name = r.name
	}

	return r.clone(repo, name)
}

func (r *repoManager) CloneFromBranch(name, branch string) error {
	fmt.Println("Clonando repo...")

	repo := r.getRepoURL()
	// Clonar usando git clone repo shell

	if name == "" {
		name = r.name
	}

	return r.cloneFromBranch(repo, name, branch)
}

func (r *repoManager) CloneDH(name string) error {
	fmt.Println("Clonando repo...")

	repo := r.getRepoURLFromDH()

	if name == "" {
		name = r.name
	}

	return r.clone(repo, name)
}

func (r *repoManager) CloneFromBranchDH(name, branch string) error {
	repo := r.getRepoURLFromDH()

	// Clonar usando git clone repo shell
	if name == "" {
		name = r.name
	}
	// git clone
	return r.cloneFromBranch(repo, name, branch)
}

func (r *repoManager) PushChanges(message string) error {
	fmt.Println("Agregando cambios...")
	err := r.execRepoGitCommand("add", ".")
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar cambios: %s", err.Error())
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

func (r *repoManager) InviteUsers(users []string) error {
	var req inviteReq
	req.Permission = "push"
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("Invitando a %s\n", user)
		url := fmt.Sprintf("/repos/%s/%s/collaborators/%s", r.username, r.name, user)
		err := r.apiManager.Put(url, body, nil)
		if err != nil {
			color.Print("red", fmt.Sprintf("Error al invitar al usuario %s", user))
		}
	}

	return nil
}

func (r *repoManager) initializeRepo(desc string) error {
	err := r.Clone("")
	if err != nil {
		return err
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

	repo := r.getRepoURL()

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

func (r *repoManager) clone(repo, name string) error {
	fmt.Println("Clonando repo...")
	// git clone
	err := exec.Command("git", "clone", repo, name).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al clonar el repositorio")
	}

	return nil
}

func (r *repoManager) cloneFromBranch(repo, name, branch string) error {
	fmt.Println("Clonando repo...")
	// git clone
	err := exec.Command("git", "clone", "--single-branch", "--branch", branch, repo, name).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al clonar el repositorio")
	}

	return nil
}

func (r *repoManager) getRepoURL() string {
	return fmt.Sprintf("git@github.com:%s/%s.git", r.username, r.name)
}

func (r *repoManager) getRepoURLFromDH() string {
	return fmt.Sprintf("git@github.com:ezedh/%s.git", r.name)
}
