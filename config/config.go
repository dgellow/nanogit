package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/qrclabs/nanogit/log"
)

type ConfigInfo struct {
	ConfigFile string
	Conf       Config
}

type ServerConfig struct {
	Root  string
	User  string
	Group string
}

type TeamConfig struct {
	Name  string
	Write bool
	Read  bool
}

type OrgConfig struct {
	Id          string
	Description string
	Teams       []TeamConfig
}

type PubKeyConfig struct {
	Type string
	Val  string
}

type UserOrgConfig struct {
	Id    string
	Teams []string
}

type UserConfig struct {
	Name    string
	SSHKeys []PubKeyConfig
	Orgs    []UserOrgConfig
}

type Config struct {
	Server ServerConfig
	Orgs   []OrgConfig
	Users  []UserConfig
}

func (ci *ConfigInfo) ReadFile() {
	data, err := ioutil.ReadFile(ci.ConfigFile)
	if err != nil {
		panic(err)
	}

	t := Config{}
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		log.Fatal("config: cannot deserialize config file: %s, error: %v", ci.ConfigFile, err)
	}
	ci.Conf = t
}

func (ci *ConfigInfo) LookupUserByKey(k string) (UserConfig, error) {
	log.Trace("config: LookupUserByKey")
	for _, user := range ci.Conf.Users {
		for _, key := range user.SSHKeys {
			log.Trace("config: key: %s", key)
			if key.Val == k {
				return user, nil
			}
		}
	}
	return UserConfig{}, fmt.Errorf("Cannot find given key in config")
}

func (ci *ConfigInfo) LookupOrgById(orgId string) (OrgConfig, error) {
	log.Trace("config: LookupOrgById, orgId: %v", orgId)
	for _, org := range ci.Conf.Orgs {
		log.Trace("config: %v", org.Id)
		log.Trace("config: length: %t", org.Id == orgId)
		if org.Id == orgId {
			return org, nil
		}
	}
	return OrgConfig{}, fmt.Errorf("Cannot find org in config: %s", orgId)
}
