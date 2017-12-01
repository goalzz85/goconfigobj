package goconfigobj

import (
	"strings"
	"testing"
)

func TestNoSection(t *testing.T) {
	txt := `
    aa    =aa
    bb =gg
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Value("aa") == "aa" && co.Value("bb") == "gg" {
		t.Log("NoSection test ok")
	} else {
		t.Error("NoSection test failed")
	}
}

func TestNoSectionSameKeyReplace(t *testing.T) {
	txt := `
    aa    =aa
    bb =gg
    "aa" = "gggg"
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Value("aa") == "gggg" && co.Value("bb") == "gg" {
		t.Log("NoSectionSameKeyReplace test ok")
	} else {
		t.Error("NoSectionSameKeyReplace test failed")
	}
}

func TestNoSectionSomeQuot(t *testing.T) {
	txt := `
    "aa"   ="aa"
    'bb'=   'gg'
    "bb3"  ='gg122'
    'bbc'=   "ggg"
    "bbcx"=   ggg
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Value("bbcx") == "ggg" && co.Value("aa") == "aa" {
		t.Log("NoSectionSomeQuot test ok")
	} else {
		t.Error("NoSectionSomeQuot test failed")
	}
}

func TestNoSectionMultiLine(t *testing.T) {
	txt := `
    "aa"   ="""
    ffff
    ffff
    """
    `

	c := `
    ffff
    ffff
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Value("aa") == c {
		t.Log("NoSectionMultiLine test ok")
	} else {
		t.Error("NoSectionMultiLine test failed")
	}
}

func TestSection(t *testing.T) {
	txt := `
    [fff]
    aa    =aa
    bb =gg
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Section("fff").Value("aa") == "aa" && co.Section("fff").Value("bb") == "gg" {
		t.Log("Section test ok")
	} else {
		t.Error("Section test failed")
	}
}

func TestSectionQuot(t *testing.T) {
	txt := `
    ['fff']
    aa    =aa
    ["sss"]
    bb =gg
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Section("fff").Value("aa") == "aa" && co.Section("sss").Value("bb") == "gg" {
		t.Log("Section test ok")
	} else {
		t.Error("Section test failed")
	}
}

func TestSectionDepth(t *testing.T) {
	txt := `
    [fff]

    bb =gg
    [[zz]]
    [[["ggg"]]]
    aa    =aa
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Section("fff").Section("zz").Section("ggg").Value("aa") == "aa" && co.Section("fff").Value("bb") == "gg" {
		t.Log("SectionDepth test ok")
	} else {
		t.Error("SectionDepth test failed")
	}
}

func TestSectionMultiDepth(t *testing.T) {
	txt := `
	aa = bb

	[fff]
		bb =gg
		[[zz]]
			[[["ggg"]]]
				aa    =aa

	[eee]
		bb =gg
		[[zz]]
			[[["ggg"]]]
			aa    =aa
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Value("aa") == "bb" && co.Section("fff").Section("zz").Section("ggg").Value("aa") == "aa" && co.Section("fff").Value("bb") == "gg" && co.Section("eee").Section("zz").Section("ggg").Value("aa") == "aa" && co.Section("eee").Value("bb") == "gg" {
		t.Log("SectionMultiDepth test ok")
	} else {
		t.Error("SectionMultiDepth test failed")
	}
}

func TestChinese(t *testing.T) {
	txt := `
    [fff]

    bb =gg
    [[zz]]
    [[["ggg"]]]
    aa    =开心猫
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Section("fff").Section("zz").Section("ggg").Value("aa") == "开心猫" && co.Section("fff").Value("bb") == "gg" {
		t.Log("Chinese test ok")
	} else {
		t.Error("Chinese test failed")
	}
}

func TestChineseQuot(t *testing.T) {
	txt := `
    [fff]

    bb ="你好"
    [[zz]]
    [[["ggg"]]]
    aa    ='开心猫'
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Section("fff").Section("zz").Section("ggg").Value("aa") == "开心猫" && co.Section("fff").Value("bb") == "你好" {
		t.Log("ChineseQuot test ok")
	} else {
		t.Error("ChineseQuot test failed")
	}
}

func TestHttpUrl(t *testing.T) {
	txt := `
    [fff]

    bb ="你好"
    [[zz]]
    [[["ggg"]]]
    aa    =https://www.google.co.jp/search?q=%E5%BC%80%E5%BF%83%E7%8C%AB&oq=%E5%BC%80%E5%BF%83%E7%8C%AB&aqs=chrome..69i57.1391j0j4&sourceid=chrome&ie=UTF-8
    `
	co := NewConfigObj(strings.NewReader(txt))
	if co.Section("fff").Section("zz").Section("ggg").Value("aa") == "https://www.google.co.jp/search?q=%E5%BC%80%E5%BF%83%E7%8C%AB&oq=%E5%BC%80%E5%BF%83%E7%8C%AB&aqs=chrome..69i57.1391j0j4&sourceid=chrome&ie=UTF-8" && co.Section("fff").Value("bb") == "你好" {
		t.Log("HttpUrl test ok")
	} else {
		t.Error("HttpUrl test failed")
	}
}
