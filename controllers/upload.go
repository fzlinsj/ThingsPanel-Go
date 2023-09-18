package controllers

import (
	gvalid "ThingsPanel-Go/initialize/validate"
	"ThingsPanel-Go/services"
	"ThingsPanel-Go/utils"
	response "ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	context2 "github.com/beego/beego/v2/server/web/context"
)

type UploadController struct {
	beego.Controller
}

// func (uploadController *UploadController) UpForm() {
// 	uploadController.TplName = "upload.tpl"
// }

func (c *UploadController) List() {
	reqData := valid.TpVisPluginPaginationValidate{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqData)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(reqData)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(reqData, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(c.Ctx))
			break
		}
		return
	}
	//获取租户id
	tenantId, ok := c.Ctx.Input.GetData("tenant_id").(string)
	if !ok {
		response.SuccessWithMessage(400, "代码逻辑错误", (*context2.Context)(c.Ctx))
		return
	}
	var tpvisplugin services.TpVis
	isSuccess, d, t := tpvisplugin.GetBlackGroudImgList(reqData, tenantId)

	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(c.Ctx))
		return
	}
	dd := valid.RspTpOtaPaginationValidate{
		CurrentPage: reqData.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     reqData.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(c.Ctx))

}

//多文件上传
func (c *UploadController) UpImgFiles() {

	fileType := c.GetString("type")
	if fileType == "" {
		response.SuccessWithMessage(1000, "类型为空", (*context2.Context)(c.Ctx))
	} else {
		err := utils.CheckPath(fileType)
		if err != nil {
			response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		}
	}
	//获取租户id
	tenantId, ok := c.Ctx.Input.GetData("tenant_id").(string)
	if !ok {
		response.SuccessWithMessage(400, "代码逻辑错误", (*context2.Context)(c.Ctx))
		return
	}

	files, err := c.GetFiles("files")
	if err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var visfiles []map[string]string
	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
			return
		}
		defer file.Close()
		//创建目录
		uploadDir := "./files/" + fileType + "/" + time.Now().Format("2006-01-02/")
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
			return
		}
		//构造文件名称
		rand.Seed(time.Now().UnixNano())
		randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
		hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))
		ext := path.Ext(files[i].Filename)
		fileName := fmt.Sprintf("%x", hashName) + ext
		err = utils.CheckFilename(fileName)
		if err != nil {
			response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
			return
		}
		fpath := uploadDir + fileName

		dst, err := os.Create(fpath)
		if err != nil {
			response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
			return
		}

		visfiles = append(visfiles, map[string]string{
			"file_name":   files[i].Filename,
			"file_url":    fpath,
			"file_size":   fmt.Sprintf("%d", files[i].Size),
			"file_remark": fileType + "_" + tenantId,
		})
	}
	var tpvisplugin services.TpVis
	isSuccess := tpvisplugin.UploadBlackGroundImg(tenantId, visfiles)
	if !isSuccess {
		utils.SuccessWithMessage(1000, "上传失败", (*context2.Context)(c.Ctx))
		return
	}
	utils.SuccessWithMessage(200, "上传成功", (*context2.Context)(c.Ctx))

}

func (uploadController *UploadController) UpFile() {
	fileType := uploadController.GetString("type")
	if fileType == "" {
		response.SuccessWithMessage(1000, "类型为空", (*context2.Context)(uploadController.Ctx))
	} else {
		err := utils.CheckPath(fileType)
		if err != nil {
			response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(uploadController.Ctx))
		}
	}
	f, h, _ := uploadController.GetFile("file") //获取上传的文件
	ext := path.Ext(h.Filename)
	//验证后缀名是否符合要求
	var AllowExtMap map[string]bool = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".svg":  true,
		".ico":  true,
		".gif":  true,
	}
	var AllowUpgradePackageMap map[string]bool = map[string]bool{
		".bin":  true,
		".tar":  true,
		".gz":   true,
		".zip":  true,
		".gzip": true,
		".apk":  true,
		".dav":  true,
		".pack": true,
	}
	var AllowimportBatchMap map[string]bool = map[string]bool{
		".xlsx": true,
	}
	switch fileType {
	case "upgradePackage":
		if _, ok := AllowUpgradePackageMap[ext]; !ok {
			response.SuccessWithMessage(1000, "文件类型不正确", (*context2.Context)(uploadController.Ctx))
			return
		}
	case "importBatch":
		if _, ok := AllowimportBatchMap[ext]; !ok {
			response.SuccessWithMessage(1000, "文件类型不正确", (*context2.Context)(uploadController.Ctx))
			return
		}
	case "imporBackground":
		if _, ok := AllowExtMap[ext]; !ok {
			response.SuccessWithMessage(1000, "文件类型不正确", (*context2.Context)(uploadController.Ctx))
			return
		}

	case "d_plugin":
		// 不做限制
	default:
		if _, ok := AllowExtMap[ext]; !ok {
			response.SuccessWithMessage(1000, "文件类型不正确", (*context2.Context)(uploadController.Ctx))
			return
		}
	}

	//创建目录
	uploadDir := "./files/" + fileType + "/" + time.Now().Format("2006-01-02/")
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(uploadController.Ctx))
	}
	//构造文件名称
	rand.Seed(time.Now().UnixNano())
	randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
	hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))
	fileName := fmt.Sprintf("%x", hashName) + ext
	err = utils.CheckFilename(fileName)
	if err != nil {
		response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(uploadController.Ctx))
	}
	fpath := uploadDir + fileName

	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	err = uploadController.SaveToFile("file", fpath)
	if fileType == "upgradePackage" {
		fpath = "./api/ota/download/files/" + fileType + "/" + time.Now().Format("2006-01-02/") + fileName
	}
	if err != nil {
		response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(uploadController.Ctx))
	}

	if fileType == "imporBackground" {

		//获取租户id
		tenantId, ok := uploadController.Ctx.Input.GetData("tenant_id").(string)
		if !ok {
			response.SuccessWithMessage(400, "代码逻辑错误", (*context2.Context)(uploadController.Ctx))
			return
		}

		var visfiles []map[string]string
		visfiles = append(visfiles, map[string]string{
			"file_name":   fileName,
			"file_url":    fpath,
			"file_size":   fmt.Sprintf("%d", h.Size),
			"file_remark": "imporBackground",
		})

		var tpvisplugin services.TpVis
		isSuccess := tpvisplugin.UploadBlackGroundImg(tenantId, visfiles)
		if !isSuccess {
			response.SuccessWithMessage(1000, "上传失败", (*context2.Context)(uploadController.Ctx))
			return
		}

	}

	response.SuccessWithDetailed(200, "success", fpath, map[string]string{}, (*context2.Context)(uploadController.Ctx))
}

func (uploadController *UploadController) DeleteBackGroundImgFile() {

	id := uploadController.GetString("id")
	if id == "" {
		response.SuccessWithMessage(1000, "类型为空", (*context2.Context)(uploadController.Ctx))
	}

}
