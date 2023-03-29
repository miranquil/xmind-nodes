package xmind_nodes

import (
	"encoding/xml"
	"errors"
)

var ErrInvalidXmindFile = errors.New("invalid xmind file")

type XmindFile struct {
	Sheets []Sheet
}

type Sheet struct {
	Id        string      `json:"id"`
	Title     string      `json:"title"`
	RootTopic *XmindTopic `json:"rootTopic"`
}

type XmindTopic struct {
	Id       string    `json:"id"`
	Title    string    `json:"title"`
	Children *Children `json:"children"`
}

type Children struct {
	Attached []*XmindTopic `json:"attached"`
	Detached []*XmindTopic `json:"detached"`
}

type Topic struct {
	Id     string   `json:"id"`
	Title  string   `json:"title"`
	Topics []*Topic `json:"topics"`
}

type xmlContent struct {
	XMLName xml.Name   `xml:"xmap-content"`
	Sheets  []XMLSheet `xml:"sheet"`
}

type XMLSheet struct {
	XMLName xml.Name  `xml:"sheet"`
	Title   string    `xml:"title"`
	Topic   *XMLTopic `xml:"topic"`
	Id      string    `xml:"id,attr"`
}

type XMLTopic struct {
	XMLName  xml.Name     `xml:"topic"`
	Title    string       `xml:"title"`
	Children *XMLChildren `xml:"children"`
	Id       string       `xml:"id,attr"`
}

type XMLChildren struct {
	XMLName        xml.Name            `xml:"children"`
	TypeTopicsList []*XMLChildrenTopic `xml:"topics"`
}

type XMLChildrenTopic struct {
	XMLName xml.Name    `xml:"topics"`
	Type    string      `xml:"type,attr"`
	Topics  []*XMLTopic `xml:"topic"`
}
