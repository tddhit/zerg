package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Parser map[string]Parser `json:"parser"`
	Writer map[string]string `json:"writer"`
	Seed   Seed              `json:"seed"`
}

type Parser struct {
	CssSelector string `json:"cssSelector"`
	Parser      string `json:"parser"`
	Writer      string `json:"writer"`
	Type        string `json:"type"`
}

type Seed struct {
	File   string `json:"file"`
	Parser string `json:"parser"`
}

func NewConf(path string) (*Conf, error) {
	c := &Conf{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
