package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

		ytID := strings.Split(strings.Split(req.URL, "v=")[1], "&")[0]
		var ext string
		if req.Format == Video {
			ext = ".mp4"
		} else {
			ext = ".mp3"
		}

		fileExists, _ := FileExists("./out/" + ytID + ext)
		filename := ytID + ext

		println(fileExists, filename)

		if !fileExists {
			filename = download(req.URL, req.Format)
		}

		c.File("./out/" + filename)
	})

	r.Run("127.0.0.1:8080")
}

func download(url string, mediaType MediaType) string {
	fmt.Printf("Downloading %s in %s format", url, mediaType)

	var executable = path.Join("yt-dlp")

	// Check if the OS is Windows
	os := runtime.GOOS
	if os == "windows" {
		executable += ".exe"
	}

	commandsArgs := []string{"--no-playlist", "-o", "./out/%(id)s.%(ext)s"}

	if mediaType == Video {
		commandsArgs = append(commandsArgs, "-f", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/bestvideo[height<=720]+bestaudio", "--merge-output-format", "mp4")
	} else if mediaType == Audio {
		commandsArgs = append(commandsArgs, "-f", "bestaudio", "--extract-audio", "--audio-format", "mp3")
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

		if mediaType == Audio {
			if strings.Contains(line, "[ExtractAudio] Destination:") {
				filename = strings.Split(line, "out\\")[1]
				break
			}

			if strings.Contains(line, "file is already in target format") {
				filename = strings.Split(strings.Split(line, "out\\")[1], ";")[0]
				break
			}
		}

		if mediaType == Video {
			if strings.Contains(line, "[download] Destination:") {
				filename = strings.Split(line, "out\\")[1]
			}

			if strings.Contains(line, "has already been downloaded") {
				filename = strings.Split(strings.Split(line, "out\\")[1], " ")[0]
			}
		}

	}

	println("Filename: ", filename)
	return filename
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
