package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

var openaiToken = os.Getenv("OPENAI_API_KEY")

func GenerateDescription(title, milestoneDesc, complemento string) (string, error) {
	if openaiToken == "" {
		return "", errors.New("OPENAI_API_KEY não configurada")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}

	// validate if complemento is not null or empty
	if complemento != "" {
		complemento = "\n Considere também esse complemento de contexto: " + complemento
	}

	prompt := "Você é um analista de sistemas e seu objetivo é criar descritivos para tarefas para times de tecnologia. " +
		" \n Gere um descritivo detalhado para um card de tarefa com o título: '" +
		title +
		"' considerando o contexto: '" +
		milestoneDesc + "'." +
		complemento +
		"\n Não inclua o título na descrição a ser retornada, apenas use como fonte de informação." +
		"\n Retorne o conteúdo no format markdown."
	payload := map[string]interface{}{
		"model": cfg.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+openaiToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New(string(b))
	}
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("sem resposta da OpenAI")
	}
	return result.Choices[0].Message.Content, nil
}
