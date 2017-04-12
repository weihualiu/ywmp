package resource

import (
	//"fmt"
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

import "github.com/weihualiu/ywmp/src/lib"
import . "github.com/weihualiu/ywmp/src/services"

type GenDescFile struct {
	name   string
	detail *GenDescFileDetail
}

type GenDescFileList []*GenDescFile

type GenDescFileDetail struct {
	DescType   string `json:"p"`
	Filever    string `json:"r"`
	MimeType   string `json:"mime"`
	EnsureType int64  `json:"u"`
}

// Main process
func GenH5Desc(filePath string) error {
	//fmt.Println("GenH5Desc() :", filePath)
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer r.Close()
	fname := strings.Split(path.Base(filePath), ".")[0]
	//fmt.Println("fname :", fname)
	for _, f := range r.File {
		//fmt.Println("zip file name: ", f.Name)
		if !strings.HasPrefix(f.Name, fname) {
			os.Remove(filePath)
			return errors.New("zip struct error!")
		}
	}

	basePath, _ := Config.Get("RESOURCE_BASE").(string)
	basePathH5 := basePath + "h5_resources/"
	descFileName := basePathH5 + "desc/h5-" + fname + "-common-3.desc"
	// delete exist funcCode resources
	os.RemoveAll(basePathH5 + fname)
	os.Remove(descFileName)
	// write funcCode to store and funcCode offline description
	var descList GenDescFileList
	err = descList.parseAndStoreFile(basePathH5, fname, r.File)
	if err != nil {
		return err
	}

	err = descList.storeDesc(descFileName)
	if err != nil {
		return err
	}
	return nil
}

// uncompress zip file
func (this *GenDescFileList) parseAndStoreFile(funcPath, fname string, files []*zip.File) error {
	for _, f := range files {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return errors.New("read zip's file failed")
		}

		b, err := ioutil.ReadAll(rc)
		if err != nil {
			rc.Close()
			return err
		}
		rc.Close()

		// local store zip files
		storeFile, err := lib.FileCreate(funcPath + f.Name)
		if err != nil {
			return err
		}
		io.Copy(storeFile, rc)
		storeFile.Close()

		filename := strings.Replace(f.Name, fname+"/", "", -1)
		fileDetail := &GenDescFileDetail{"h5", lib.SHA128_Buf(b), lib.MimeType(f.Name), 1}
		*this = append(*this, &GenDescFile{filename, fileDetail})
		b = nil
	}

	return nil
}

func (this GenDescFileList) storeDesc(path string) error {
	desc, err := this.String()
	if err != nil {
		return err
	}
	f, err := lib.FileCreate(path)
	if err != nil {
		return err
	}
	defer f.Close()
	io.WriteString(f, desc)
	return nil
}

func (this GenDescFileList) String() (res string, err error) {
	sort.Sort(GenDescFileList(this))
	for _, k := range this {
		b, err := json.Marshal(k.detail)
		if err != nil {
			return "", err
		}
		//fmt.Println("json to string: ", lib.ByteToString(b))
		res += "\"" + k.name + "\":" + lib.ByteToString(b) + ","
	}
	//fmt.Println("result:", strings.TrimSuffix(res, ","))
	descLists := "{" + strings.TrimSuffix(res, ",") + "}"
	res = "{\"hash\":\"" + lib.SHA128(descLists) + ",\"h5\":" + descLists + "}"
	//fmt.Println("desc:", res)
	return res, nil
}

func (this GenDescFileList) Len() int {
	return len(this)
}

func (this GenDescFileList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this GenDescFileList) Less(i, j int) bool {
	return strings.ToLower(this[i].name) > strings.ToLower(this[j].name)
}

func NewGenDescFileDetail(descType, filever, mimeType string, ensureType int64) *GenDescFileDetail {
	return &GenDescFileDetail{descType, filever, mimeType, ensureType}
}

func NewGenDescFile(name string, detail *GenDescFileDetail) *GenDescFile {
	return &GenDescFile{name, detail}
}
