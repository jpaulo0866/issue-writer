// cmd/setup.go
package cmd

import (
	"bufio"
	"fmt"
	"issue-writer/internal"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configura a CLI",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("URL base do GitLab: ")
		baseURL, _ := reader.ReadString('\n')
		fmt.Print("Personal Access Token: ")
		token, _ := reader.ReadString('\n')
		fmt.Print("ID do usuário do GitLab: ")
		userID, _ := reader.ReadString('\n')
		fmt.Print("ID do projeto (opcional): ")
		projectID, _ := reader.ReadString('\n')
		fmt.Print("ID do grupo do GitLab: ")
		groupID, _ := reader.ReadString('\n')

		cfg := internal.Config{
			Gitlab: internal.GitlabConfig{
				BaseURL:   strings.TrimSpace(baseURL),
				Token:     strings.TrimSpace(token),
				UserID:    strings.TrimSpace(userID),
				ProjectID: strings.TrimSpace(projectID),
				GroupID:   strings.TrimSpace(groupID),
			},
		}
		if err := internal.SaveConfig(cfg); err != nil {
			fmt.Println("Erro ao salvar configuração:", err)
			return
		}
		fmt.Println("Configuração salva com sucesso!")
	},
}
