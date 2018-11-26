package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"fmt"
	"flag"
	"io/ioutil"
	"regexp"
	"strings"
	"path"
)

var (
	RootPath string
)

func main() {
	RootPath = fmt.Sprintf("%s/src/github.com/illidan33/markdown_text", os.Getenv("GOPATH"))
	port := flag.Int("port", 8001, "listen port")
	flag.Parse()

	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.LoadHTMLGlob(RootPath + "/html/*")
	router.Static("/js", RootPath+"/js")
	router.Static("/css", RootPath+"/css")
	router.Static("/files", RootPath+"/create_files")
	router.GET("/", IndexRouter)
	router.GET("/Detail/:name", DetailRouter)
	router.POST("/save", SaveRouter)

	router.Run(fmt.Sprintf(":%d", *port))
}

func IndexRouter(c *gin.Context) {
	fileNames := []string{}
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/create_files", RootPath))
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fileNames = append(fileNames, strings.TrimSuffix(file.Name(), path.Ext(file.Name())))
		}
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"files": fileNames,
	})
}
func DetailRouter(c *gin.Context) {
	content := []byte("")
	name := c.Param("name")
	if name != "" {
		filePath := fmt.Sprintf("%s/create_files/%s.tpl", RootPath, name)

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
	var err error
	var newPath string
	c.Request.ParseForm()
	content := c.PostForm("content")
	name := c.PostForm("name")
	oldName := c.PostForm("old_name")

	// 校验名称
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"msg": "name empty",
		})
		return
	} else {
		reg := regexp.MustCompile("^[\\-_0-9a-zA-Z\u4E00-\u9FA5]+$")
		if ok := reg.MatchString(name); !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": "名称只能为中文、数字、字母、-、_，不能含有特殊字符!",
			})
			return
		}
	}

	// 如果改名,不能覆盖同名文件
	if name != oldName {
		newPath = fmt.Sprintf("/Detail/%s", name)
		files, err := ioutil.ReadDir(fmt.Sprintf("%s/create_files/", RootPath))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		nameStr := fmt.Sprintf("%s.tpl", name)
		for _, file := range files {
			if file.Name() == nameStr {
				c.JSON(http.StatusOK, gin.H{
					"msg": "已有同名文件",
				})
				return
			}
		}
	}

	filePath := fmt.Sprintf("%s/create_files/%s.tpl", RootPath, name)

	var handle *os.File
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

	handle.Write([]byte(content))

	if newPath != "" {
		os.Remove(fmt.Sprintf("%s/create_files/%s.tpl", RootPath, oldName))
	}
	c.JSON(http.StatusOK, gin.H{
		"path": newPath,
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
