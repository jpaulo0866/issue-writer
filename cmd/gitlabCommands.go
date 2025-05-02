package cmd

import (
	"fmt"
	"issue-writer/internal"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	titulo      string
	epico       string
	milestone   string
	labels      string
	userID      string
	projectID   string
	complemento string
)

var issueCmd = &cobra.Command{
	Use:   "new-issue",
	Short: "Gera um card no GitLab",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Println("Erro ao carregar configuração:", err)
			os.Exit(1)
		}

		// Usa valores do setup se não informados
		if userID == "" {
			userID = cfg.Gitlab.UserID
		}
		if projectID == "" {
			projectID = cfg.Gitlab.ProjectID
		}

		labelList := []string{}
		if labels != "" {
			labelList = strings.Split(labels, ",")
		}

		gitlabCfg := cfg.Gitlab

		// Busca milestone do GitLab
		milestoneInfo, err := internal.GetGroupMilestone(gitlabCfg, gitlabCfg.GroupID, milestone)
		if err != nil {
			fmt.Println("Erro ao buscar milestone:", err)
			os.Exit(1)
		}

		epicInfo, err := internal.GetGroupEpic(gitlabCfg, gitlabCfg.GroupID, epico)
		if err != nil {
			fmt.Println("Erro ao buscar epico:", err)
			os.Exit(1)
		}

		// Gera descrição via OpenAI
		description, err := internal.GenerateDescription(titulo, milestoneInfo.Description, complemento)
		if err != nil {
			fmt.Println("Erro ao gerar descrição:", err)
			os.Exit(1)
		}

		// Busca iteration mais próxima
		iterationID, err := internal.GetClosestIteration(cfg, projectID)
		if err != nil {
			fmt.Println("Erro ao buscar iteration:", err)
			os.Exit(1)
		}

		// Cria o card no GitLab
		url, err := internal.CreateIssue(cfg, projectID, titulo, description, epicInfo.ID, milestoneInfo.ID, labelList, userID, iterationID)
		if err != nil {
			fmt.Println("Erro ao criar card:", err)
			os.Exit(1)
		}
		fmt.Println("Card criado com sucesso!")
		fmt.Println(url)
	},
}

var milestoneCmd = &cobra.Command{
	Use:   "list-milestones",
	Short: "Lista as milestones do projeto configurado",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Println("Erro ao carregar configuração:", err)
			os.Exit(1)
		}
		gitlabCfg := cfg.Gitlab
		milestones, err := internal.ListGroupMilestone(gitlabCfg, gitlabCfg.GroupID)
		if err != nil {
			fmt.Println("Erro ao buscar milestones:", err)
			os.Exit(1)
		}

		fmt.Println("Milestones encontradas:")
		for _, milestone := range milestones {
			fmt.Println(" - ", milestone.Title)
		}
	},
}

var epicsCmd = &cobra.Command{
	Use:   "list-epics",
	Short: "Lista os épicos do projeto configurado",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := internal.LoadConfig()
		if err != nil {
			fmt.Println("Erro ao carregar configuração:", err)
			os.Exit(1)
		}
		gitlabCfg := cfg.Gitlab
		epics, err := internal.ListGroupEpics(gitlabCfg, gitlabCfg.GroupID)
		if err != nil {
			fmt.Println("Erro ao buscar epicos:", err)
			os.Exit(1)
		}

		fmt.Println("Epicos encontrados:")
		for _, epic := range epics {
			fmt.Println(" - ", epic.Title)
		}
	},
}

var gitlabRootCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Operações no GITLAB",
}

func init() {
	issueCmd.Flags().StringVarP(&titulo, "titulo", "t", "", "Título do card")
	issueCmd.Flags().StringVarP(&epico, "epico", "e", "", "Epico")
	issueCmd.Flags().StringVarP(&milestone, "milestone", "m", "", "Milestone")
	issueCmd.Flags().StringVarP(&labels, "labels", "l", "", "Lista de labels separadas por vírgula")
	issueCmd.Flags().StringVarP(&userID, "user", "u", "", "ID do usuário (opcional)")
	issueCmd.Flags().StringVarP(&projectID, "project", "p", "", "ID do projeto (opcional)")
	issueCmd.Flags().StringVarP(&complemento, "complemento", "c", "", "Complemento de prompt (opcional)")
	issueCmd.MarkFlagRequired("titulo")
	issueCmd.MarkFlagRequired("epico")
	issueCmd.MarkFlagRequired("milestone")

	gitlabRootCmd.AddCommand(issueCmd)
	gitlabRootCmd.AddCommand(epicsCmd)
	gitlabRootCmd.AddCommand(milestoneCmd)
}
