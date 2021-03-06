package main

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConf(t *testing.T) {
	c := &Conf{
		Parser: make(map[string]Parser),
		Writer: make(map[string]string),
	}
	c.Parser["baiduHref"] = Parser{
		CssSelector: ".result .c-title a",
		Type:        "href",
		Parser:      "newsText",
	}
	c.Parser["newsText"] = Parser{
		CssSelector: "p",
		Writer:      "tv",
		Type:        "text",
	}
	c.Writer["tv"] = "data/tv.txt"
	c.Seed = Seed{
		File:   "data/seed.txt",
		Parser: "baiduHref",
	}
	out, _ := yaml.Marshal(c)
	fmt.Println(string(out))
}
