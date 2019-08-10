package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gobuffalo/packr"

	"gopkg.in/gomail.v2"
	yaml "gopkg.in/yaml.v2"
)

type SmtpSetting struct {
	Host     string `yaml: "host"`
	Port     int    `yaml: "port"`
	UserName string `yaml: "username"`
	Password string `yaml: "password"`
}

func main() {
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	to := parser.String("t", "to", &argparse.Options{Required: true, Help: "recipent address (required)"})
	from := parser.String("f", "from", &argparse.Options{Required: true, Help: "sender address (required)"})
	subject := parser.String("s", "subject", &argparse.Options{Required: true, Help: "mail subject (required)"})
	message := parser.String("m", "message", &argparse.Options{Required: true, Help: "mail text_body message (required)"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	// Finally print the collected string

	sendmail(*to, *from, *subject, *message)
}

func yamlData() string {
	box := packr.NewBox("./config/data.yml.d")
	yamlStr, err := box.FindString("/smtp.yml")
	if err != nil {
		log.Fatalf("smtp.yml not found: %v", err)
	}
	return yamlStr
}

func smtp() SmtpSetting {
	var data SmtpSetting
	err := yaml.Unmarshal([]byte(yamlData()), &data)
	if err != nil {
		log.Fatalf("smtp.yaml can not parse: %v", err)
	}
	return data
}

func sendmail(to string, from string, subject string, message string) {
	m := gomail.NewMessage()

	//m.SetBody("text/html", HtmlBody)
	m.AddAlternative("text/plain", message)

	// if you use formataddress, then write bellow
	//"From":    {m.FormatAddress(Sender, SenderName)},
	//"To":      {m.FormatAddress(Recipent, RecipentName)},
	m.SetHeaders(map[string][]string{
		"From": {from},
		//"From":    {m.FormatAddress(Sender, SenderName)},
		"To":      {to},
		"Subject": {subject},
	})

	conf := smtp()
	d := gomail.NewPlainDialer(conf.Host, conf.Port, conf.UserName, conf.Password)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent!")
	}

}
