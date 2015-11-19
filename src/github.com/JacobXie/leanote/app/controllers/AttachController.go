package controllers

import (
	"github.com/revel/revel"
	//	"encoding/json"
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/JacobXie/leanote/app/info"
	. "github.com/JacobXie/leanote/app/lea"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// 附件
type Attach struct {
	BaseController
}

// 上传附件
func (c Attach) UploadAttach(noteId string) revel.Result {
	re := c.uploadAttach(noteId)
	return c.RenderJson(re)
}
func (c Attach) uploadAttach(noteId string) (re info.Re) {
	var fileId = ""
	var resultMsg = "error" // 错误信息
	var Ok = false
	var fileInfo info.Attach

	re = info.NewRe()

	defer func() {
		re.Id = fileId // 只是id, 没有其它信息
		re.Msg = resultMsg
		re.Ok = Ok
		re.Item = fileInfo
	}()

	// 判断是否有权限为笔记添加附件
	if !shareService.HasUpdateNotePerm(noteId, c.GetUserId()) {
		return re
	}

	file, handel, err := c.Request.FormFile("file")
	if err != nil {
		return re
	}
	defer file.Close()

	// 生成上传路径
	filePath := "files/" + c.GetUserId() + "/attachs"

	// 生成新的文件名
	filename := handel.Filename
	_, ext := SplitFilename(filename) // .doc
	filename = NewGuid() + ext

	// > 5M?
	maxFileSize := configService.GetUploadSize("uploadAttachSize")
	if maxFileSize <= 0 {
		maxFileSize = 1000
	}
	//modified by JacobXie
	//filesize = int64(len(data))
	var filesize int64
	// 获取文件大小的接口
	type Sizer interface {
			Size() int64
	}

	if sizeInterface, ok := file.(Sizer); ok {
			filesize = sizeInterface.Size()
	}

	if float64(filesize) > maxFileSize*float64(1024*1024) {
		resultMsg = fmt.Sprintf("The file's size is bigger than %vM", maxFileSize)
		return re
	}

	if qiniuService.IsUseQiniu(){
		toPath := filePath + "/" + filename
		err := qiniuService.Upload2Qiniu(toPath,file,filesize)
		if err != nil {
			return re
		}
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return re
		}
		dir := revel.BasePath + "/" + filePath
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return re
		}

		toPath := dir + "/" + filename
		err = ioutil.WriteFile(toPath, data, 0777)
		if err != nil {
			return re
		}
		filesize = GetFilesize(toPath)
	}

	// add File to db
	fileType := ""
	if ext != "" {
		fileType = strings.ToLower(ext[1:])
	}

	fileInfo = info.Attach{Name: filename,
		Title:        handel.Filename,
		NoteId:       bson.ObjectIdHex(noteId),
		UploadUserId: c.GetObjectUserId(),
		Path:         filePath + "/" + filename,
		Type:         fileType,
		Size:         filesize}

	id := bson.NewObjectId()
	fileInfo.AttachId = id
	fileId = id.Hex()
	Ok, resultMsg = attachService.AddAttach(fileInfo, false)
	if resultMsg != "" {
		resultMsg = c.Message(resultMsg)
	}

	fileInfo.Path = "" // 不要返回
	if Ok {
		resultMsg = "success"
	}
	return re
}

// 删除附件
func (c Attach) DeleteAttach(attachId string) revel.Result {
	re := info.NewRe()
	re.Ok, re.Msg = attachService.DeleteAttach(attachId, c.GetUserId())
	return c.RenderJson(re)
}

// get all attachs by noteId
func (c Attach) GetAttachs(noteId string) revel.Result {
	re := info.NewRe()
	re.Ok = true
	re.List = attachService.ListAttachs(noteId, c.GetUserId())
	return c.RenderJson(re)
}

// 下载附件
// 权限判断
func (c Attach) Download(attachId string) revel.Result {
	attach := attachService.GetAttach(attachId, c.GetUserId()) // 得到路径
	path := attach.Path
	if path == "" {
		return c.RenderText("")
	}
	//modified by JacobXie
	if qiniuService.IsUseQiniu(){
		return c.Redirect(qiniuService.GetUrlOnQiniu(strings.TrimLeft(path,"/")))
	} else {
		fn := revel.BasePath + "/" + strings.TrimLeft(path, "/")
		file, _ := os.Open(fn)
		return c.RenderBinary(file, attach.Title, revel.Attachment, time.Now()) // revel.Attachment
	}
	// return c.RenderFile(file, revel.Attachment) // revel.Attachment
}

func (c Attach) DownloadAll(noteId string) revel.Result {
	//modified by JacobXie
	//如果使用七牛，不支持打包下载
	if qiniuService.IsUseQiniu(){
		return c.RenderText("")
	}
	note := noteService.GetNoteById(noteId)
	if note.NoteId == "" {
		return c.RenderText("")
	}

	// 得到文件列表
	attachs := attachService.ListAttachs(noteId, c.GetUserId())
	if attachs == nil || len(attachs) == 0 {
		return c.RenderText("")
	}

	/*
		dir := revel.BasePath + "/files/tmp"
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return c.RenderText("")
		}
	*/

	filename := note.Title + ".tar.gz"
	if note.Title == "" {
		filename = "all.tar.gz"
	}

	// file write
	fw, err := os.Create(revel.BasePath + "/files/" + filename)
	if err != nil {
		return c.RenderText("")
	}
	// defer fw.Close() // 不需要关闭, 还要读取给用户下载

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 遍历文件列表
	for _, attach := range attachs {
		fn := revel.BasePath + "/" + strings.TrimLeft(attach.Path, "/")
		fr, err := os.Open(fn)
		fileInfo, _ := fr.Stat()
		if err != nil {
			return c.RenderText("")
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = attach.Title
		h.Size = fileInfo.Size()
		h.Mode = int64(fileInfo.Mode())
		h.ModTime = fileInfo.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			panic(err)
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			panic(err)
		}
	} // for

	//    tw.Close()
	//    gw.Close()
	//    fw.Close()
	//    file, _ := os.Open(dir + "/" + filename)
	// fw.Seek(0, 0)
	return c.RenderBinary(fw, filename, revel.Attachment, time.Now()) // revel.Attachment
}
