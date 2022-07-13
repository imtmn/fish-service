package common

import (
	"encoding/json"
	"log"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
	"mtmn.top/fish-service/entity"
)

var (
	APPLY_URL string = "https://www.mtmn.top/api/v1/apply"
)

func SendMailApply(form *entity.ApplicationForm) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "xxx<xxx@126.com>"
	// 设置接收方的邮箱
	e.To = []string{"xxx@qq.com"}
	//设置主题
	e.Subject = "申请请求：" + form.Name
	//设置文件发送的内容
	e.HTML = []byte(`
		<li>名称：<span>` + form.Name + `</span></li>
		<li>说明：<span>` + form.Remark + `</span></li>
		<li>姓名：<span>` + form.PersonName + `</span></li>
		<li>电话：<span>` + form.Phone + `</span></li>
		<li>备注：<span>` + form.Info + `</span></li>
		<li>图片：<span>` + strings.Join(form.Imgs, "<br/>") + `</span></li>
		<li><img src="` + form.Imgs[0] + `"></img></li>
		<br/>
		<a  href="` + APPLY_URL + "/pass/" + form.Id + `" >审批通过</a>
		<br/><br/>
		<a  href="` + APPLY_URL + "/reject/" + form.Id + `" >审批不通过</a>
	`)
	//设置服务器相关的配置
	e.Send("smtp.126.com:25", smtp.PlainAuth("", "xxx@126.com", "RIELLOLXKJEKBKKE", "smtp.126.com"))

}

func SendMailObj(suject string, params interface{}) {
	paramsStr, err := json.Marshal(params)
	if err != nil {
		log.Println("发送邮件失败，参数为:", params)
	}
	SendMail(suject, string(paramsStr))
}

func SendMail(suject string, msg string) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "xxx<xxx@126.com>"
	// 设置接收方的邮箱
	e.To = []string{"xxx@qq.com"}
	//设置主题
	e.Subject = suject
	//设置文件发送的内容
	e.Text = []byte(msg)
	//设置服务器相关的配置
	err := e.Send("smtp.126.com:25", smtp.PlainAuth("", "xxxx@126.com", "RIELLOLXKJEKBKKE", "smtp.126.com"))
	if err != nil {
		log.Fatal(err)
	}
}
