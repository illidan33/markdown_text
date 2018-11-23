package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"flag"
)

var (
	RootPath string
)

func main() {
	RootPath = fmt.Sprintf("%s/src/gotest/markdown_test/", os.Getenv("GOPATH"))
	port := flag.Int("port", 8001, "listen port")
	flag.Parse()

	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.LoadHTMLGlob(RootPath + "html/*")
	router.Static("/js", RootPath+"js")
	router.Static("/css", RootPath+"css")
	router.Static("/files", RootPath+"create_files")
	router.GET("/", IndexRouter)
	router.POST("/save", SaveRouter)

	router.Run(fmt.Sprintf(":%d", *port))
}

func IndexRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func SaveRouter(c *gin.Context) {
	c.Request.ParseForm()
	content := c.PostForm("content");
	name := c.PostForm("name");
	fileName := fmt.Sprintf("%s.html", name)
	filePath := fmt.Sprintf("%screate_files/%s", RootPath, fileName)

	var handle *os.File
	var err error
	if IsExists(filePath) {
		os.Remove(filePath)
	}

	handle, err = os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer handle.Close()

	headerHandle, _ := os.Open(fmt.Sprintf("%shtml/header.tpl", RootPath))
	defer headerHandle.Close()
	headerHtml, _ := ioutil.ReadAll(headerHandle)
	handle.Write(headerHtml)

	_, err = handle.Write([]byte(fmt.Sprintf("%s</div> </div> </body> </html>", content)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"path": fmt.Sprintf("http://%s/files/%s", c.Request.Host, fileName),
	})
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
