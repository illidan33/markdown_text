package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"fmt"
	"flag"
	"io/ioutil"
	"regexp"
)

var (
	RootPath string
)

func main() {
	RootPath = fmt.Sprintf("%s/src/github.com/illidan33/markdown_text/", os.Getenv("GOPATH"))
	port := flag.Int("port", 8001, "listen port")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.LoadHTMLGlob(RootPath + "html/*")
	router.Static("/js", RootPath+"js")
	router.Static("/css", RootPath+"css")
	router.Static("/files", RootPath+"create_files")
	router.GET("/", IndexRouter)
	router.GET("/Detail/:name", DetailRouter)
	router.POST("/save", SaveRouter)

	router.Run(fmt.Sprintf(":%d", *port))
}

func IndexRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
func DetailRouter(c *gin.Context) {
	content := []byte("")
	name := c.Param("name")
	if name != "" {
		filePath := fmt.Sprintf("%screate_files/%s.tpl", RootPath, name)

		if IsExists(filePath) {
			handle, err := os.Open(filePath)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "500", gin.H{})
				return
			}
			defer handle.Close()

			content, _ = ioutil.ReadAll(handle)
		}
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"content": string(content),
		"name":    name,
	})
}

func SaveRouter(c *gin.Context) {
	c.Request.ParseForm()
	content := c.PostForm("content")
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "name empty",
		})
		return
	} else {
		reg := regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`)
		if ok := reg.MatchString(name); !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "name error",
			})
			return
		}
	}

	filePath := fmt.Sprintf("%screate_files/%s.tpl", RootPath, name)

	var handle *os.File
	var err error
	if IsExists(filePath) {
		os.Remove(filePath)
	}

	handle, err = os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "name error",
		})
		return
	}
	defer handle.Close()

	handle.Write([]byte(content))
	c.JSON(http.StatusOK, gin.H{
		"path": fmt.Sprintf("http://%s/Detail/%s", c.Request.Host, name),
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
