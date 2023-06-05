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
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"os"
)

var config *Config

type Config struct {
	MailAddress2DingToken *MailAddress2DingToken `yaml:"mailAddressDingTokenMap"`
	Password              string                 `yaml:"password"`
	Port                  int                    `yaml:"port"`
	SmtpClient            *SmtpClient            `yaml:"smtpClient"`
}

type SmtpClient struct {
	SmtpServer     string `yaml:"smtpServer"`
	SmtpPort       int    `yaml:"smtpPort"`
	Tls            bool   `yaml:"tls"`
	SenderEmail    string `yaml:"senderEmail"`
	SenderPassword string `yaml:"senderPassword"`
}

type MailAddress2DingToken map[string]string

func LoadConfiguration(file string) {
	var cfg Config
	readFile(&cfg, file)
	glog.Infof("config content %+v", cfg)
	config = &cfg
}

func readFile(cfg *Config, file string) {
	f, err := os.Open(file)
	if err != nil {
		glog.Fatalf("fail to read config file: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		glog.Fatalf("fail to decode config file yaml: %w", err)
	}
}
