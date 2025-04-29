# Issue Writer CLI

CLI para gerar cards com descrição automática via OpenAI.

## Pré-requisitos

- Go 1.18+
- Ter um token de acesso pessoal do GitLab
- Ter uma chave de API da OpenAI (`OPENAI_API_KEY` no ambiente)

## Compilação

### Linux/MacOS

```sh
go build -o dist/issue-writer
```

### Windows

```sh
go build -o issue-writer.exe
```

## Uso

### Configuração inicial:

`./issue-writer setup`

**Parâmetros:**

- URL base do GitLab ( Sem barra no final )
- Personal Access Token
- ID do usuário do GitLab ( {URL_GIT}/-/user_settings/profile - Campo User ID )
- ID do projeto (opcional - Projeto onde a Issue será registrada )
- ID do grupo do GitLab ( Id do grupo onde reside a Milestone e o épico - Funciona a hierarquia de cima para baixo, logo o grupo mais alto enxerga os épicos e milestones dos grupos abaixo )

### Gerar card:

`./issue-writer new-issue --titulo "Título do Card" --epico "Epico" --milestone "Milestone" --labels "bug,backend"`

**Parâmetros opcionais:**

- --user
- --project

### Variáveis de ambiente

- OPENAI_API_KEY: sua chave da OpenAI

## Observações

- O projeto salva as configurações em ~/.issue_writer_cli_config.json
- Por enquanto, apenas integração com GitLab.

## Adicionar no path do S.O

Copie o binário gerado na pasta dist para /usr/local/bin

```sh
sudo cp dist/issue-writer /usr/local/bin
```
