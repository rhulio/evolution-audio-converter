# Evolution Audio Converter

Este projeto é um micro serviço em Go que processa arquivos de áudio para whatsapp, os converte para o formato **opus** e retorna tanto a duração do áudio quanto o arquivo convertido em base64. O serviço aceita arquivos de áudio enviados como **form-data**, **base64** ou **URL**.

## Requisitos

Antes de começar, você precisará ter os seguintes itens instalados:

- [Go](https://golang.org/doc/install) (versão 1.21 ou superior)
- [Docker](https://docs.docker.com/get-docker/) (para rodar o projeto em um container)
- [FFmpeg](https://ffmpeg.org/download.html) (para processamento de áudio)

## Instalação

### Clonar o Repositório

Clone este repositório em sua máquina local:

```bash
git clone https://github.com/EvolutionAPI/evolution-audio-converter.git
cd evolution-audio-converter
```

### Instalar Dependências

Instale as dependências do projeto:

```bash
go mod tidy
```

### Instalar o FFmpeg

O serviço depende do **FFmpeg** para converter o áudio. Certifique-se de que o FFmpeg está instalado no seu sistema.

- No Ubuntu:

  ```bash
  sudo apt update
  sudo apt install ffmpeg
  ```

- No MacOS (via Homebrew):

  ```bash
  brew install ffmpeg
  ```

- No Windows, baixe o FFmpeg [aqui](https://ffmpeg.org/download.html) e adicione-o ao `PATH` do sistema.

### Configuração

Crie um arquivo `.env` no diretório raiz do projeto com a seguinte configuração:

```env
PORT=4040
```

Isso define a porta onde o serviço será executado.

## Rodando o Projeto

### Localmente

Para rodar o serviço localmente, use o seguinte comando:

```bash
go run main.go
```

O servidor estará disponível em `http://localhost:4040`.

### Usando Docker

Se preferir rodar o serviço em um container Docker, siga os passos abaixo:

1. **Buildar a imagem Docker**:

   ```bash
   docker build -t audio-service .
   ```

2. **Rodar o container**:

   ```bash
   docker run -p 4040:4040 --env-file=.env audio-service
   ```

   Isso irá iniciar o container na porta especificada no arquivo `.env`.

## Como Usar

Você pode enviar requisições `POST` para o endpoint `/process-audio` com um arquivo de áudio nos seguintes formatos:

- **Form-data** (para enviar arquivos)
- **Base64** (para enviar o áudio codificado em base64)
- **URL** (para enviar o link do arquivo de áudio)

### Exemplo de Requisição via cURL

#### Envio como Form-data

```bash
curl -X POST -F "file=@caminho/do/audio.mp3" http://localhost:4040/process-audio
```

#### Envio como Base64

```bash
curl -X POST -d "base64=$(base64 caminho/do/audio.mp3)" http://localhost:4040/process-audio
```

#### Envio como URL

```bash
curl -X POST -d "url=https://exemplo.com/caminho/para/audio.mp3" http://localhost:4040/process-audio
```

### Resposta

A resposta será um JSON contendo a duração do áudio e o arquivo convertido em base64:

```json
{
  "duration": 120,
  "audio": "UklGR... (base64 do arquivo)"
}
```

- `duration`: A duração do áudio em segundos.
- `audio`: O arquivo de áudio convertido para o formato opus, codificado em base64.

## Licença

Este projeto está sob a licença [MIT](LICENSE).
```

### Resumo do Conteúdo:
- **Requisitos**: Dependências necessárias (Go, Docker, FFmpeg).
- **Instalação**: Como clonar o repositório, instalar dependências e configurar o arquivo `.env`.
- **Rodando o projeto**: Instruções para rodar o projeto localmente ou via Docker.
- **Como usar**: Exemplos de como fazer requisições para o serviço com `form-data`, `base64` ou URL, e o formato de resposta.
- **Licença**: Incluir uma seção para a licença, se aplicável.