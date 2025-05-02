package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Milestone struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Epic struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func ListGroupMilestone(cfg GitlabConfig, groupID string) ([]Milestone, error) {
	url := fmt.Sprintf("%s/api/v4/groups/%s/milestones?state=active&include_descendants=true&per_page=1000", cfg.BaseURL, groupID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	var milestones []Milestone
	if err != nil {
		return milestones, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&milestones); err != nil {
		return milestones, err
	}

	return milestones, nil
}

func ListGroupEpics(cfg GitlabConfig, groupID string) ([]Epic, error) {
	url := fmt.Sprintf("%s/api/v4/groups/%s/epics?include_descendants=true&per_page=1000", cfg.BaseURL, groupID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	var epics []Epic
	if err != nil {
		return epics, err
	}
	defer resp.Body.Close()

	fmt.Printf("Status Code %d", resp.StatusCode)

	if err := json.NewDecoder(resp.Body).Decode(&epics); err != nil {
		return epics, err
	}

	return epics, nil
}

func GetGroupMilestone(cfg GitlabConfig, groupID, milestoneTitle string) (Milestone, error) {
	// encode the milestoneTitle
	encodedMilestoneTitle := url.QueryEscape(milestoneTitle)
	url := fmt.Sprintf("%s/api/v4/groups/%s/milestones?search=%s&include_descendants=true", cfg.BaseURL, groupID, encodedMilestoneTitle)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Milestone{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Erro na busca da milestone %d - %s - %s", resp.StatusCode, groupID, milestoneTitle)
		return Milestone{}, errors.New("milestone não encontrada")
	}
	var milestones []Milestone
	if err := json.NewDecoder(resp.Body).Decode(&milestones); err != nil {
		return Milestone{}, err
	}
	for _, m := range milestones {
		if m.Title == milestoneTitle {
			return m, nil
		}
	}
	return Milestone{}, errors.New("milestone não encontrada")
}

func GetGroupEpic(cfg GitlabConfig, groupID, epicTitle string) (Epic, error) {
	encodedTitle := url.QueryEscape(epicTitle)
	url := fmt.Sprintf("%s/api/v4/groups/%s/epics?search=%s&include_descendant_groups=true", cfg.BaseURL, groupID, encodedTitle)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Epic{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Epic{}, errors.New("épico não encontrado")
	}
	var epics []Epic
	if err := json.NewDecoder(resp.Body).Decode(&epics); err != nil {
		return Epic{}, err
	}
	for _, e := range epics {
		if e.Title == epicTitle {
			return e, nil
		}
	}
	return Epic{}, errors.New("épico não encontrado")
}

func GetClosestIteration(cfg Config, projectID string) (int, error) {
	// Exemplo: busca a primeira iteration aberta
	url := fmt.Sprintf("%s/api/v4/projects/%s/iterations?state=opened", cfg.Gitlab.BaseURL, projectID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PRIVATE-TOKEN", cfg.Gitlab.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return 0, errors.New("falha ao buscar no gitlab")
	}
	var iterations []struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		StartDate time.Time `json:"start_date"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&iterations); err != nil {
		return 0, nil
	}
	if len(iterations) == 0 {
		return 0, nil
	}
	return iterations[0].ID, nil
}

func CreateIssue(cfg Config, projectID, title, description string, epic, milestone int, labels []string, userID string, iterationId int) (string, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%s/issues", cfg.Gitlab.BaseURL, projectID)
	payload := map[string]interface{}{
		"title":        title,
		"description":  description,
		"milestone_id": milestone,
		"epic_id":      epic,
		"labels":       labels,
		"assignee_ids": []string{userID},
	}

	if iterationId != 0 {
		payload["iteration_id"] = iterationId
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("PRIVATE-TOKEN", cfg.Gitlab.Token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("erro ao criar issue: %s", string(b))
	}

	// Parse the response to get the web_url
	var result struct {
		WebURL string `json:"web_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.WebURL == "" {
		return "", errors.New("issue criada, mas URL não encontrada na resposta")
	}
	return result.WebURL, nil
}
