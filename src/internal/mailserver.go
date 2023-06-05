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
	"errors"
	"github.com/emersion/go-smtp"
	"github.com/golang/glog"
	"io"
	"strconv"
	"strings"
	"time"
)

// The Backend implements SMTP server methods.
type Backend struct {
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct {
	from string
	to   string
}

func (s *Session) AuthPlain(username, password string) error {
	if password != (*config).Password {
		glog.Errorf("Invalid password, username:%s\n", username)
		return errors.New("invalid password")
	}
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		return err
	} else {
		original := string(b)
		subject, content := getSubjectAndContent(original)
		dingToken := ""
		if token, ok := (*(*config).MailAddress2DingToken)[s.to]; ok {
			dingToken = token
		} else {
			glog.Warningf("fail to find ding token for email address %s\n", s.to)
		}
		glog.Infof("toMail: %s toDing: %s subject: %s", s.to, dingToken, subject)

		go sendToDingTalkRobot(dingToken, subject, content)
		go sendEmail((*config).SmtpClient, s.to, original)
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

// getSubjectAndContent split mail message to subject and content
func getSubjectAndContent(mailContent string) (string, string) {
	//https://www.rfc-editor.org/rfc/rfc822.html
	//     A message consists of header fields and, optionally, a body.
	//     The  body  is simply a sequence of lines containing ASCII charac-
	//     ters.  It is separated from the headers by a null line  (i.e.,  a
	//     line with nothing preceding the CRLF).
	strs := strings.Split(mailContent, "\n")
	str := ""
	subject := ""
	for _, s := range strs {
		if len(subject) == 0 {
			if strings.Index(s, "Subject:") == 0 {
				subject = s[8:]
			}
		}
		str += strings.TrimSpace(s) + "\n"
	}
	if len(strings.TrimSpace(subject)) == 0 {
		subject = "WARNING"
	}

	index := strings.Index(str, "\n\n")
	if index == -1 {
		glog.Warningf("fail to find null line: %s", str)
		return subject, str
	}

	return subject, str[index+2:]
}

// InitServer init mail server
func InitServer() {

	be := &Backend{}
	s := smtp.NewServer(be)

	s.Addr = ":" + strconv.Itoa((*config).Port)
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	glog.Infof("Starting server at %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		glog.Fatal(err)
	}
}
