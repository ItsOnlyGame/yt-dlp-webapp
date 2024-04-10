package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"os/exec"
)

type DownloadRequest struct {
	URL    string    `json:"url"`
	Format MediaType `json:"format"`
}

type MediaType string

const (
	Audio MediaType = "audio"
	Video MediaType = "video"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/download/:filename", func(c *gin.Context) {
		c.File("./out/" + c.Param("filename"))
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/download", func(c *gin.Context) {
		var req DownloadRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate request fields (URL and format)
		if req.URL == "" || req.Format == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL and format fields are required"})
			return
		}

		filename := download(req.URL, req.Format)

		c.File("./out/" + filename)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func download(url string, mediaType MediaType) string {
	fmt.Printf("Downloading %s in %s format", url, mediaType)

	var executable = path.Join("yt-dlp")

	// Check if the OS is Windows
	os := runtime.GOOS
	if os == "windows" {
		executable += ".exe"
	}

	commandsArgs := []string{"-o", "./out/%(id)s.%(ext)s", "-f", "bestaudio"}
	if mediaType == Audio {
		commandsArgs = []string{"-o", "./out/%(id)s.%(ext)s", "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3"}
	}

	commandsArgs = append(commandsArgs, url)

	cmd := exec.Command(executable, commandsArgs...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

	output := out.String()
	outputArr := strings.Split(output, "\n")

	filename := ""

	for _, line := range outputArr {
		println(line)

		if strings.Contains(line, "[ExtractAudio] Destination:") {
			filename = strings.Split(line, "out\\")[1]
		}

		if strings.Contains(line, "file is already in target format") {
			filename = strings.Split(strings.Split(line, "out\\")[1], ";")[0]
		}

	}

	println("Filename: ", filename)
	return filename
}
