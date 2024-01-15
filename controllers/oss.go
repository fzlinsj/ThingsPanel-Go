package controllers

import (
	response "ThingsPanel-Go/utils"

	beego "github.com/beego/beego/v2/server/web"
	context2 "github.com/beego/beego/v2/server/web/context"
)

type OssInfo struct {
	BucketName string

	bucketURL string
}

type OssController struct {
	beego.Controller
}

func (this *OssController) GetOssInfo() {

	// authorization := this.Ctx.Request.Header["Authorization"][0]
	// userToken := authorization[7:]
	// _, err := jwt.ParseCliamsToken(userToken)
	// if err != nil {
	// 	response.SuccessWithMessage(400, "token异常", (*context2.Context)(this.Ctx))
	// 	return
	// }

	host := this.Ctx.Request.Host

	http := "http"

	if this.Ctx.Request.TLS != nil {

		http = "https"

	}

	oosurl := http + "://" + host + "/api/goview/project/getImages/"

	d := OssInfo{
		BucketName: "getuserphoto",
		bucketURL:  oosurl,
	}

	response.SuccessWithDetailed(200, "返回成功", d, map[string]string{}, (*context2.Context)(this.Ctx))

}
