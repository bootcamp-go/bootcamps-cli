package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type createReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

type RepoManager interface {
	// CreateRepo creates a new repository given a name and description.
	CreateRepo(name, desc string) error
	// PushChanges pushes changes to the repository.
	PushChanges(name string, commit string) error
}

type repoManager struct {
	token    string
	username string
}

func NewRepoManager(token, username string) RepoManager {
	return &repoManager{
		token:    token,
		username: username,
	}
}

func (r *repoManager) CreateRepo(name, desc string) error {
	creq := &createReq{
		Name:        name,
		Description: desc,
		Private:     true,
	}

	postBody, err := json.Marshal(creq)
	if err != nil {
		return err
	}

	url := "https://api.github.com/user/repos"

	client := &http.Client{}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(postBody))

	req.Header.Set("Authorization", fmt.Sprintf("token %s", r.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type request struct{}

	// decore resp body into req
	var rr request
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		return fmt.Errorf("error decoding response: %s", err)
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("bad response: %d", resp.StatusCode)
	}

	return r.initializeRepo(name, desc)
}

func (r *repoManager) PushChanges(name string, commit string) error {
	fmt.Println("Agregando cambios...")
	err := exec.Command("git", "-C", name, "add", ".").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar el readme")
	}

	fmt.Println("Haciendo commits de cambios...")
	err = exec.Command("git", "-C", name, "commit", "-m", commit).Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("ocurrió un error al hacer commit de los cambios: %s", err))
	}

	fmt.Println("Subiendo cambios...")
	err = exec.Command("git", "-C", name, "push", "-u", "origin", "main").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al pushear a main")
	}

	return nil
}

func (r *repoManager) initializeRepo(name, desc string) error {
	fmt.Println("Clonando repo...")

	repo := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s/%s.git", r.token, r.username, name)

	// Clonar usando git clone repo shell
	// git clone
	err := exec.Command("git", "clone", repo).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al clonar el repositorio")
	}

	fmt.Println("Creando readme...")
	// Create file README.md inside the repo folder
	err = exec.Command("touch", fmt.Sprintf("./%s/README.md", name)).Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al crear el readme")
	}

	fmt.Println("Inicializando repo...")
	err = exec.Command("git", "-C", name, "init", ".").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al incializar el repositorio")
	}

	fmt.Println("Agregando readme...")
	err = exec.Command("git", "-C", name, "add", ".").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al agregar el readme")
	}

	fmt.Println("Haciendo commits de cambios...")
	err = exec.Command("git", "-C", name, "commit", "-m", "initial commit").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al hacer commit de los cambios")
	}

	fmt.Println("Creando branch main...")
	err = exec.Command("git", "-C", name, "branch", "-M", "main").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al crear branch main")
	}

	fmt.Println("Configurando origin...")
	_ = exec.Command("git", "-C", name, "remote", "add", "origin", repo).Run()

	err = exec.Command("git", "-C", name, "push", "-u", "origin", "main").Run()
	if err != nil {
		return fmt.Errorf("ocurrió un error al pushear a main")
	}

	return nil
}
