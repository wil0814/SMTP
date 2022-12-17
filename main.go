package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var Val Config

type Config struct {
	Sender     string `mapstructure:"SENDER"`
	PassWord   string `mapstructure:"PASSWORD"`
	SMTPServer string `mapstructure:"SMTP_SERVER"`
	Port       string `mapstructure:"PORT"`
	Recipient  string `mapstructure:"RECIPIENT"`
}

func main() {
	// 使用密碼驗證，也可以用OAth2驗證
	auth := smtp.PlainAuth(
		"",
		Val.Sender,     //你的gmail
		Val.PassWord,   //密碼
		Val.SMTPServer, //伺服器
	)
	//取得email body
	email, err := ioutil.ReadFile("email.txt")
	if err != nil {
		log.Fatal(err)
	}
	body := string(email)

	imageByte, err := ioutil.ReadFile("./波及王子.jpeg")

	if err != nil {
		log.Fatal(err)
	}
	msg := []string{
		"Subject: 來自Golang的祝福",
		"From:" + Val.Sender,
		"Content-Type: multipart/mixed; boundary=\"frontier\"",
		"",
		"--frontier",
		"Content-Type: text/html",
		"",
		"<h1 align=\"center\" style=\"color:red;\">Wish You a Merry Christmas</h1>",
		body,
		"--frontier",
		"Content-Type: image/jpeg;",
		"Content-Transfer-Encoding: base64",
		"",
		string(base64.StdEncoding.EncodeToString(imageByte)),
		"--frontier",
	}
	//發送郵件
	err = smtp.SendMail(
		Val.SMTPServer+Val.Port,           //伺服器＋port
		auth,                              //認證
		Val.Sender,                        //發件人
		[]string{Val.Recipient},           //收件人
		[]byte(strings.Join(msg, "\r\n")), //發送訊息
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sent email")
}
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("讀取設定檔出現錯誤: %v", err))
	}
	if err := viper.Unmarshal(&Val); err != nil {
		panic(fmt.Errorf("找不到Struct, %v", err))
	}
	log.WithFields(log.Fields{
		"val": Val,
	}).Info("config loaded")
}
