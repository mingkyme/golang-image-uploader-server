package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	staticPath := os.Getenv("STATIC_PATH")
	// serverURL := os.Getenv("SERVER_URL") // Request 의 URL 따라 응답.
	port := os.Getenv("PORT")

	httpMode := os.Getenv("HTTPS_MODE") == "ON"
	var certPath string
	var keyPath string
	if httpMode {
		certPath = os.Getenv("HTTPS_CERT_PATH")
		keyPath = os.Getenv("HTTPS_KEY_PATH")
	}
	app := fiber.New()
	os.MkdirAll(staticPath, os.ModePerm)
	app.Static("/static", staticPath)
	app.Post("/upload", func(c *fiber.Ctx) error {
		mpf, err := c.MultipartForm()
		if err != nil {
			fmt.Println(err)
			return c.Status(400).SendString("send multipart/form-data")
		}
		if len(mpf.File["file"]) == 0 {
			fmt.Println("file is empty")
			return c.Status(400).SendString("send with key 'file'")
		}
		fileHeader := mpf.File["file"][0]
		fileName := fileHeader.Filename                              // image.png
		fileExtension := fileName[strings.LastIndex(fileName, "."):] // .png
		randomFileName := randSeq(30) + fileExtension
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		}
		// Date
		t := time.Now()
		timeStr := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
		os.Mkdir(path.Join(staticPath, timeStr), 0777)

		f, err := os.OpenFile(path.Join(staticPath, timeStr, randomFileName), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		defer f.Close()
		io.Copy(f, file)
		fmt.Println(c.BaseURL())
		return c.SendString(fmt.Sprintf("%s/static", c.BaseURL()) + fmt.Sprintf("%s/%s", timeStr, randomFileName))
	})

	if httpMode {
		err = app.ListenTLS(port, certPath, keyPath)
	} else {
		err = app.Listen(port)
	}
	if err != nil {
		fmt.Println(err)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
