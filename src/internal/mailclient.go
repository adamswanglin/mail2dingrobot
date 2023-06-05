/*
 Copyright 2023 adamswanglin

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package internal

import (
	"crypto/tls"
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/smtp"
)

func sendEmail(smtpClient *SmtpClient, to, message string) {
	if smtpClient == nil {
		//未配置不发送邮件
		return
	}
	client := *smtpClient

	auth := smtp.PlainAuth("", client.SenderEmail, client.SenderPassword, client.SmtpServer)
	smtpAddress := fmt.Sprintf("%s:%d", client.SmtpServer, client.SmtpPort)
	receivers := []string{to}
	var err error
	if client.Tls {
		err = SendMailTls(smtpAddress, auth, client.SenderEmail, receivers, []byte(message))
	} else {
		err = smtp.SendMail(smtpAddress, auth, client.SenderEmail, receivers, []byte(message))
	}

	if err != nil {
		glog.Errorf("Email sent error to %s, error: %w\n", to, err)
	}
}

// SendMailTls tls端口发送邮件
func SendMailTls(addr string, a smtp.Auth, from string, to []string, msg []byte) error {

	host, _, _ := net.SplitHostPort(addr)
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	defer c.Quit()
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(a); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
