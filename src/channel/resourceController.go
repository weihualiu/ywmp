package channel

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

import "github.com/weihualiu/ywmp/src/lib"

//import resource "channel/resource"
import . "github.com/weihualiu/ywmp/src/services"

type Sizer interface {
	Size() int64
}

/*
 * 离线资源生成、下载
 */
type resourceGenController struct {
}

func (ic *resourceGenController) Process(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//fmt.Println("urlpath: ", path)
	switch path {
	case "/resource/":
		//上传页面
		uploadPath := "public/www/resource/upload.html"
		http.ServeFile(w, r, uploadPath)
	case "/resource/upload":
		uploadAndGen(w, r)
	case "/resource/pack":
		io.WriteString(w, packPage())
	case "/resource/packget":
		packGet(w, r)
	case "/resource/getdesc":
		getDesc(w, r)
	case "/resource/getfile":
		getFile(w, r)
	default:
		fmt.Printf("%s\n", path)
		http.ServeFile(w, r, "public/www"+path)
		//http.NotFound(w, r)
	}
}

//zip上传
//离线资源生成
func uploadAndGen(w http.ResponseWriter, r *http.Request) {
	var aj ajaxJson
	var status string = "false"

	//在defer中处理返回值
	var result string
	defer func() {
		aj = ajaxJson{Status: status, Msg: result}
		io.WriteString(w, aj.String())
	}()
	r.ParseMultipartForm(32 << 20)
	file, fileheader, err := r.FormFile("zipfile")
	if fileheader == nil {
		result = "错误：未选择上传文件"
		return
	}
	//fmt.Println("filename: ", fileheader)
	envtag := r.FormValue("envsel")
	//fmt.Println("selected: ", envtag)
	if err != nil {
		log.Println("FromFile: ", err.Error())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Close: ", err.Error())
			return
		}

	}()
	//判定文件是否是zip后缀
	if !strings.HasSuffix(fileheader.Filename, ".zip") {
		result = "警告：上传文件不是zip！"
		return
	}
	if len(fileheader.Filename) != 14 {
		result = "警告：上传文件名不符合规则！"
		return
	}
	//检查目录名是不是小写
	if lib.IsContainUpper(fileheader.Filename) {
		result = "警告：文件名包含大写字母！"
		return
	}
	//判定文件大小
	if file.(Sizer).Size() > 1024*1024*10 {
		err = errors.New("upload zip size overflow!")
		result = "警告：上传文件超过10MB！"
		return
	}

	// 直接落地文件到指定位置
	tmpPath, _ := Config.Get("RESOURCE_TMP").(string)
	tmpFile := tmpPath + "/" + envtag + "/" + fileheader.Filename
	fzip, err := lib.FileCreate(tmpFile)
	if err != nil {
		result = err.Error()
		return
	}
	defer fzip.Close()
	io.Copy(fzip, file)
	//zip文件内容目录格式分析
	zipr, err := zip.OpenReader(tmpFile)
	if err != nil {
		result = err.Error()
		return
	}
	defer zipr.Close()
	fname := strings.Split(path.Base(tmpFile), ".")[0]
	//fmt.Println("fname :", fname)
	for _, f := range zipr.File {
		//fmt.Println("zip file name: ", f.Name)
		if !strings.HasPrefix(f.Name, fname) {
			os.Remove(tmpFile)
			result = "zip文件内容目录结构错误！"
			return
		}
	}
	//copy file to remote
	//shell
	//scp -r tmpFile mobsewp@182.207.129.67:/home/mobsewp/work/sit/offline/tmp/
	if envtag == "SIT" {
		cmd := exec.Command("/bin/sh", "-c", "scp -r "+tmpFile+" mobsewp@182.207.129.67:/home/mobsewp/work/sit/offline/tmp/")
		cmd.Run()
	} else if envtag == "UAT" {
		cmd := exec.Command("/bin/sh", "-c", "scp -r "+tmpFile+" mobsewp@182.119.81.69:/mnt/data/offline/tmp/")
		cmd.Run()
	} else if envtag == "PRD" {
		cmd := exec.Command("/bin/sh", "-c", "scp -r "+tmpFile+" regress@182.207.129.67:/home/regress/offline/tmp/")
		cmd.Run()
	} else {
		os.Remove(tmpFile)
		result = "选择部署环境异常！"
		return
	}
	os.Remove(tmpFile)

	status = "true"
	result = "上传成功！5秒后离线资源生效！"
}

type ajaxJson struct {
	Status string `json:"status"` //JSON此处要求名称必须是对外开放的字段
	Msg    string `json:"msg"`
}

func (aj ajaxJson) String() string {
	data, err := json.Marshal(aj)
	if err != nil {
		return "json format error"
	}
	return string(data)
}

//打包主页
func packPage() string {
	basePath, _ := Config.Get("RESOURCE_BASE").(string)
	//获取功能码
	url := basePath + "h5_resources/"
	names, err := readFuncCodeNames(url)
	if err != nil {
		return "execute failed!"
	}
	//names = []string{"bmosyy0101","bmosyy0201"}
	//组装页面
	page := "<html><head></head><body>"
	page += "<form name=\"input\" action=\"/resource/packget\" method=\"post\">"
	page += "投产包名：<input type=\"text\" name=\"packageName\"/ style=\"width:400px\"><br /><br />"
	page += "功能码选择：<br />"
	for _, val := range names {
		page = page + val + "<input type=\"checkbox\" name=\"funcCode\" value=\"" + val + "\" />  "
	}
	page += "<br /><br /><input type=\"submit\" value=\"获取HTML5投产包\" />"
	page += "</form>"
	page += "</body></html>"
	return page
}

//下载
func packGet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	funcCodes := r.Form["funcCode"]
	packageName := r.Form["packageName"]

	if len(funcCodes) == 0 {
		fmt.Fprintf(w, "not selected function code!")
		return
	}
	if len(packageName[0]) == 0 {
		fmt.Fprintf(w, "not input the pacakge name!")
		return
	}
	onlinebag, _ := Config.Get("RESOURCE_PACKBAG").(string)
	dir, err := mkdirBag(onlinebag)
	if err != nil {
		fmt.Fprintf(w, "mkdirbag() failed!")
		return
	}

	basePath, _ := Config.Get("RESOURCE_BASE").(string)
	tarFileName := dir + "/" + packageName[0] + ".tar"
	err = tarFiles(funcCodes, tarFileName, basePath+"h5_resources/")
	if err != nil {
		fmt.Fprintf(w, "tarFiles() failed!")
		return
	}

	//fmt.Println("tarFileName:", tarFileName)
	file, err := os.Open(tarFileName)
	defer file.Close()
	if err != nil {
		fmt.Fprintf(w, "failed!")
		return
	}
	data := bufio.NewReader(file)

	//fmt.Fprintln(w, data)
	w.Header().Set("Content-Type", "application/x-tar; charset=utf-8") // normal header
	w.Header().Set("Content-Disposition", "attachment;filename="+packageName[0]+".tar")
	w.WriteHeader(http.StatusOK)
	//io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
	data.WriteTo(w)
}

//从指定目录读取文件名
func readFuncCodeNames(dirpath string) (names []string, err error) {
	file, err := os.Open(dirpath)
	if err != nil {
		//fmt.Println("err:", err)
		return nil, err
	}
	files, err := file.Readdir(0)
	if err != nil {
		//fmt.Println("err:", err)
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if len(filename) == 10 {
			names = append(names, filename)
		}
	}
	return names, nil
}

//创建封包目录
func mkdirBag(current string) (string, error) {
	//如果不存在则创建
	t := time.Now().UTC()
	currentd := t.Format("20060102")
	all := current + currentd
	if !lib.IsExists(all) {
		if err := os.MkdirAll(all, os.ModePerm); err != nil {
			return "", err
		}
	}
	return all, nil
}

//封包投产功能码
func tarFiles(funcCodes []string, dstTar, srcDir string) (err error) {
	src := path.Clean(srcDir)
	if !lib.IsExists(src) {
		return errors.New("要打包的文件或目录不存在")
	}
	if lib.FileExists(dstTar) {
		//删除已有文件
		if er := os.Remove(dstTar); er != nil {
			return er
		}
	}

	//创建tar文件
	fw, er := os.Create(dstTar)
	if er != nil {
		return er
	}
	defer fw.Close()

	tw := tar.NewWriter(fw)
	defer func() {
		if er := tw.Close(); er != nil {
			err = er
		}
	}()
	//读取目录或者文件内容
	fi, er := os.Stat(srcDir)
	if er != nil {
		return er
	}
	srcBase, srcRelative := path.Split(path.Clean(srcDir))
	tarDir(srcBase, srcRelative, tw, fi, funcCodes)

	return nil
}

func tarDir(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo, funcCodes []string) (err error) {
	srcFull := srcBase + srcRelative
	last := len(srcRelative) - 1
	if srcRelative[last] != os.PathSeparator {
		srcRelative += string(os.PathSeparator)
	}

	fis, er := ioutil.ReadDir(srcFull)
	if er != nil {
		return er
	}

	for _, fi := range fis {
		if fi.IsDir() {
			//目录是否是功能码
			//desc不过滤
			if srcRelative == "h5_resources/" {
				flag := false
				for _, val := range funcCodes {
					if strings.Contains(fi.Name(), val) {
						flag = true
					}
				}
				if fi.Name() == "desc" {
					flag = true
				}
				if flag {
					tarDir(srcBase, srcRelative+fi.Name(), tw, fi, funcCodes)
				}
			} else {
				tarDir(srcBase, srcRelative+fi.Name(), tw, fi, funcCodes)
			}

			//tarDir(srcBase, srcRelative + fi.Name(), tw, fi, funcCodes)
		} else {
			//fmt.Println("srcRelative:", srcRelative)
			if srcRelative == "h5_resources/desc/" {
				flag := false
				for _, val := range funcCodes {
					if strings.Contains(fi.Name(), val) {
						flag = true
					}
				}
				if flag {
					tarFile(srcBase, srcRelative+fi.Name(), tw, fi)
				}
			} else {
				tarFile(srcBase, srcRelative+fi.Name(), tw, fi)
			}
		}
	}

	if len(srcRelative) > 0 {
		hdr, er := tar.FileInfoHeader(fi, "")
		if er != nil {
			return er
		}
		hdr.Name = srcRelative
		if er = tw.WriteHeader(hdr); er != nil {
			return er
		}
	}

	return nil
}

func tarFile(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) (err error) {
	srcFull := srcBase + srcRelative

	hdr, er := tar.FileInfoHeader(fi, "")
	if er != nil {
		return er
	}
	hdr.Name = srcRelative

	if er = tw.WriteHeader(hdr); er != nil {
		return er
	}

	fr, er := os.Open(srcFull)
	if er != nil {
		return er
	}
	defer fr.Close()

	if _, er = io.Copy(tw, fr); er != nil {
		return er
	}

	return nil
}

func getDesc(w http.ResponseWriter, r *http.Request) {
	// read file
	r.ParseForm()
	name := r.Form["n"]
	basePath, _ := Config.Get("RESOURCE_BASE").(string)
	descName := basePath + "h5_resources/desc/" + name[0] + ".desc"
	http.ServeFile(w, r, descName)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	// read file
	r.ParseForm()
	name := r.Form["n"]
	basePath, _ := Config.Get("RESOURCE_BASE").(string)
	fileName := basePath + "h5_resources/" + name[0]
	http.ServeFile(w, r, fileName)
}
