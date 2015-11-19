package service

import (
	"qiniupkg.com/api.v7/kodo"
	"golang.org/x/net/context"
	. "github.com/JacobXie/leanote/app/lea"
	"github.com/JacobXie/leanote/app/lea/netutil"
	"io"
	"github.com/revel/revel"
)

type QiniuService struct {
	qiniu_BaseUrl string
	qiniu_bucket	string
	bucket kodo.Bucket
	ctx context.Context
	isUsed	bool
}

//初始化七牛服务
func (this *QiniuService) InitQiniu(){
		config := revel.Config
		storage_type, ok := config.String("storage.type")
		this.isUsed = false
		if ok && storage_type == "qiniu"{
			var qiniu_AK string
			var qiniu_SK string
			qiniu_AK,ok = config.String("storage.qiniu.AK")
			if !ok || qiniu_AK == "" {
				return
			}
			qiniu_SK,ok = config.String("storage.qiniu.SK")
			if !ok || qiniu_SK == "" {
				return
			}
			this.qiniu_BaseUrl,ok = config.String("storage.qiniu.BaseUrl")
			if !ok || this.qiniu_BaseUrl == "" {
				return
			}
			this.qiniu_bucket,ok = config.String("storage.qiniu.bucket")
			if !ok || this.qiniu_bucket == "" {
				return
			}
			this.isUsed = true
			kodo.SetMac(qiniu_AK, qiniu_SK)
			zone := 0
			c := kodo.New(zone, nil)
			this.bucket = c.Bucket(this.qiniu_bucket)
			this.ctx = context.Background()
		}
}

func (this *QiniuService) IsUseQiniu() bool{
	return this.isUsed
}
//上传到七牛
//toPath 目标路径
//data 文件数据
//size 文件大小
func (this *QiniuService) Upload2Qiniu(toPath string, data io.Reader, size int64) (err error){
	return this.bucket.Put(this.ctx, nil, toPath,data,size, nil)
}

//七牛上复制文件
//pathSrc 原路径
//pathDest 目标路径
func (this *QiniuService) CopyOnQiniu(pathSrc, pathDest string)(err error){
	return this.bucket.Copy(this.ctx, pathSrc, pathDest)
}

//七牛上删除文件
func (this *QiniuService) DeleteOnQiniu(path string)(err error){
		return this.bucket.Delete(this.ctx, path)
}

//七牛上移动文件
func (this *QiniuService) MoveOnQiniu(pathSrc , pathDest string)(err error){
	return this.bucket.Move(this.ctx, pathSrc, pathDest)
}

//获取完整url
func (this *QiniuService) GetUrlOnQiniu(path string)(url string){
	return kodo.MakeBaseUrl(this.qiniu_BaseUrl,path)
}

//获取文件属性
func (this *QiniuService) GetFSizeOnQiniu(path string)(filesize int64, err error){
		entry,errT := this.bucket.Stat(this.ctx,path)
		if errT != nil{
			err = errT
			return
		}
		filesize = entry.Fsize
		err = nil
		return
}

//上传url对应文件到七牛
func (this *QiniuService) Upload2QiniuWithUrl(url, toPath string)(filesize int64,newFilename , path string,ok bool){
	if url == "" {
		return
	}
	// a.html?a=a11&xxx
	url = netutil.TrimQueryParams(url)
	_, ext := SplitFilename(url)

	newFilename = NewGuid() + ext
	path = toPath + "/" + newFilename

	err := this.bucket.Fetch(this.ctx, path, url)
	if err != nil{
		return
	}
	filesize,err = this.GetFSizeOnQiniu(path)
	if err != nil{
		return
	}
	ok = true
	return
}

//获取七牛上文件
func (this *QiniuService) GetFileOnQiniu(filePath string)(content []byte,err error){
	fullurl := this.GetUrlOnQiniu(filePath)
	content, err = netutil.GetContent(fullurl)
	return
}
