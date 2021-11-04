package invitation

import (
	"fmt"
	"os"

	"github.com/ezedh/bootcamps/internal/repo"
	"gopkg.in/yaml.v2"
)

type Inviter interface {
	Wave() string
	Company() string
	GetCreationConfig() (*Invitations, error)
}

type inviter struct {
	wave    string
	repoM   repo.RepoManager
	folder  string
	company string
}

type InvitationsConfig struct {
	Teachers map[string][]string `yaml:"teachers"`
	Groups   map[string][]string `yaml:"groups"`
}

type Invitations struct {
	Amount int
	Groups map[string][]string
}

func NewInviter(repoM repo.RepoManager, company, uuid string) Inviter {
	var wave string
	fmt.Printf("Wave N°: ")
	fmt.Scan(&wave)

	folder := fmt.Sprintf("%s-users", uuid)

	return &inviter{
		wave:    wave,
		repoM:   repoM,
		folder:  folder,
		company: company,
	}
}

func (i *inviter) Wave() string {
	return i.wave
}

func (i *inviter) Company() string {
	return i.company
}

func (i *inviter) GetCreationConfig() (*Invitations, error) {
	i.repoM.SetName("bootcamps-users")
	err := i.repoM.CloneFromBranchDH(i.folder, i.company)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(i.folder)

	ok := i.waveConfigFileExists()
	if !ok {
		return nil, fmt.Errorf("no se encontró el archivo de configuración para la wave %s de %s", i.wave, i.company)
	}

	config, err := i.getInvitationsConfig()
	if err != nil {
		return nil, err
	}

	return i.buildInvitations(config), nil
}

func (i *inviter) waveConfigFileExists() bool {
	_, err := os.Stat(fmt.Sprintf("./%s/%s", i.folder, "wave-"+i.wave+".yaml"))

	return !os.IsNotExist(err)
}

func (i *inviter) getInvitationsConfig() (*InvitationsConfig, error) {
	yamlData, err := os.ReadFile(fmt.Sprintf("./%s/%s", i.folder, "wave-"+i.wave+".yaml"))
	if err != nil {
		return nil, err
	}

	var config InvitationsConfig
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (i *inviter) buildInvitations(config *InvitationsConfig) *Invitations {
	inv := new(Invitations)
	inv.Amount = 0
	inv.Groups = config.Groups

	for t, g := range config.Teachers {
		for _, gg := range g {
			inv.Groups[gg] = append(inv.Groups[gg], t)
		}
	}

	for range inv.Groups {
		inv.Amount++
	}

	return inv
}
