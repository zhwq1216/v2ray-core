package admin

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"
)

const GfwlistUrl = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"

func init() {
	RegisterController("pac", &PacController{})
}

type PacController struct {
	admin *Server
}
type PacConfig struct {
	Proxy    string `json:"proxy"`
	UserRule string `json:"userRule"`
}

func (ctl *PacController) InitRouter(admin *Server, httpRouter gin.IRouter) {
	ctl.admin = admin
	httpRouter.POST("/pac/gfwlist/download", ctl.UpdatePac)
	httpRouter.POST("/pac/save", ctl.SavePac)
	httpRouter.GET("/pac", ctl.GetPac)
}
func (ctl *PacController) GetPac(gCtx *gin.Context) {

	pacFile := getPacV2rayFile()
	if !fileExists(pacFile) {
		gCtx.Status(404)
		gCtx.Writer.WriteString("还没有生成pac文件")
		return
	}
	pacContent, err := ioutil.ReadFile(pacFile)
	if err != nil {
		gCtx.Status(500)
		gCtx.Writer.WriteString("读取pac文件失败" + err.Error())
		return
	}
	gCtx.Header("Content-Type", "application/x-ns-proxy-autoconfig; charset=utf-8")
	gCtx.Status(200)
	gCtx.Writer.Write(pacContent)
}
func generateV2rayPac(proxy, userRule string) error {
	pacGfwFile := getPacGfwFile()
	if !fileExists(pacGfwFile) {
		if err := downloadGfwPac(); err != nil {
			return err
		}
	}
	gfwContent, _ := ioutil.ReadFile(pacGfwFile)

	gfwContent = bytes.ReplaceAll(gfwContent, Str2Bytes("__PROXY__"), Str2Bytes(proxy))
	if userRule != "" {
		userRule = regexp.MustCompile("[,;\r\n]+").ReplaceAllString(userRule, "\",\n  \"")

		userRule = "var rules = [\n  \"" + userRule + "\","
		gfwContent = bytes.ReplaceAll(gfwContent, Str2Bytes("var rules = ["), Str2Bytes(userRule))
	}
	err := ioutil.WriteFile(getPacV2rayFile(), gfwContent, 0644)
	return err
}
func downloadGfwPac() error {
	resp, err := http.Get(GfwlistUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	gfwListContent, err := ioutil.ReadAll(resp.Body)
	// bytes 转化为string的高效方法
	gfwListBytes, err := base64.StdEncoding.DecodeString(*(*string)(unsafe.Pointer(&gfwListContent)))
	gwfScanner := bufio.NewScanner(bytes.NewReader(gfwListBytes[0:]))
	lines := make([]string, 0, 4096)
	lastModifyLine := ""
	// 跳过第一行 [AutoProxy 0.2.9]
	gwfScanner.Scan()
	for gwfScanner.Scan() {
		line := gwfScanner.Text()
		if strings.Contains(line, "! Last Modified:") {
			lastModifyLine = "// GFWList"+ line[1:] + "\n"
		}
		// 跳过空行和注释行
		if line == "" || line[0] == '!' {
			continue
		}
		lines = append(lines, line)
	}
	rules, err := json.MarshalIndent(lines, "", "  ")
	rulesJson := *(*string)(unsafe.Pointer(&rules))
	tplContent, err := ioutil.ReadFile(getTemplateFile())
	if err != nil {
		return err
	}
	pacString := strings.Replace(*(*string)(unsafe.Pointer(&tplContent)), "__RULES__", rulesJson, 1)
	err = ioutil.WriteFile(getPacGfwFile(),  Str2Bytes(lastModifyLine + pacString), 0644)
	return err
}
func (ctl *PacController) UpdatePac(gCtx *gin.Context) {

	err := downloadGfwPac()
	if err != nil {
		gCtx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	config := PacConfig{}
	err = gCtx.ShouldBindJSON(&config)
	if err != nil && config.Proxy != "" {
		generateV2rayPac(config.Proxy, config.UserRule)
	}
	gCtx.JSON(200, gin.H{"msg": "下载gfwlist成功"})
	return

}
func (ctl *PacController) SavePac(gCtx *gin.Context) {
	config := PacConfig{}
	gCtx.BindJSON(&config)
	err := generateV2rayPac(config.Proxy, config.UserRule)
	if err != nil {
		gCtx.Status(500)
		gCtx.Writer.WriteString(err.Error())
		return
	}
	gCtx.JSON(200, gin.H{"msg": "生成pac成功"})
	return

}
func GetExecutableDir() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}
func getTemplateFile() string {
	return GetExecutableDir() + "/pac_template.js"
}
func getPacGfwFile() string {
	return GetExecutableDir() + "/pac_gfw.js"
}
func getPacV2rayFile() string {
	return GetExecutableDir() + "/pac_v2ray.js"
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
