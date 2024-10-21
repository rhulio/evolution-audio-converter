package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func convertAudioToOpusWithDuration(inputData []byte) ([]byte, int, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-ac", "1", "-ar", "16000", "-c:a", "libopus", "-f", "ogg", "pipe:1")

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	cmd.Stdin = bytes.NewReader(inputData)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err := cmd.Run()
	if err != nil {
		return nil, 0, fmt.Errorf("error during conversion: %v, details: %s", err, errBuffer.String())
	}

	convertedData := outBuffer.Bytes()

	outputText := errBuffer.String()

	splitTime := strings.Split(outputText, "time=")

	if len(splitTime) < 2 {
		return nil, 0, errors.New("duração não encontrada")
	}

	re := regexp.MustCompile(`(\d+):(\d+):(\d+\.\d+)`)
	matches := re.FindStringSubmatch(string(splitTime[2]))
	if len(matches) != 4 {
		return nil, 0, errors.New("formato de duração não encontrado")
	}

	hours, _ := strconv.ParseFloat(matches[1], 64)
	minutes, _ := strconv.ParseFloat(matches[2], 64)
	seconds, _ := strconv.ParseFloat(matches[3], 64)
	duration := int(hours*3600 + minutes*60 + seconds)

	return convertedData, duration, nil
}

func processAudio(c *gin.Context) {
	var inputData []byte
	var err error

	// Verifica se o arquivo foi enviado como form-data
	file, _, err := c.Request.FormFile("file")
	if err == nil {
		inputData, err = ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler o arquivo"})
			return
		}
	} else {
		// Verifica se foi enviado um base64
		base64Data := c.PostForm("base64")
		if base64Data != "" {
			inputData, err = base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar base64"})
				return
			}
		} else {
			// Verifica se foi enviada uma URL
			url := c.PostForm("url")
			if url != "" {
				resp, err := http.Get(url)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao obter o arquivo da URL"})
					return
				}
				defer resp.Body.Close()
				inputData, err = io.ReadAll(resp.Body)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler o arquivo da URL"})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum arquivo, base64 ou URL fornecido"})
				return
			}
		}
	}

	// Chama a função para converter o áudio e obter a duração
	convertedData, duration, err := convertAudioToOpusWithDuration(inputData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Codifica o áudio convertido para base64
	base64Audio := base64.StdEncoding.EncodeToString(convertedData)

	// Retorna a resposta em JSON
	c.JSON(http.StatusOK, gin.H{
		"duration": duration,
		"audio":    base64Audio,
	})
}

func main() {
	// Carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
	}

	// Obtém a porta do arquivo .env ou usa a porta 8080 por padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	// Endpoint para processar o áudio
	router.POST("/process-audio", processAudio)

	// Inicia o servidor na porta especificada
	router.Run(":" + port)
}
