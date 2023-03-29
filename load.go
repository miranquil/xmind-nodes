package xmind_nodes

import (
	"archive/zip"
	"encoding/xml"
	"io"

	jsonAPI "github.com/json-iterator/go"
)

const (
	TargetFileZen = "content.json"
	TargetFilePro = "content.xml"
)

const (
	fileTypeZen = iota
	fileTypePro
)

func Load(filePath string) (*XmindFile, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}

	var fileReader io.ReadCloser
	var fileType int
	for _, file := range reader.File {
		if file.Name == TargetFileZen || file.Name == TargetFilePro {
			if file.Name == TargetFileZen {
				fileType = fileTypeZen
			} else {
				fileType = fileTypePro
			}
			fileReader, err = file.Open()
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if fileReader == nil {
		return nil, ErrInvalidXmindFile
	}

	defer func() {
		_ = fileReader.Close()
	}()

	bytes, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	sheets, err := unmarshal(bytes, fileType)
	if err != nil {
		return nil, err
	}

	return &XmindFile{Sheets: sheets}, nil
}

func unmarshal(bytes []byte, fileType int) ([]Sheet, error) {
	switch fileType {
	case fileTypeZen:
		return unmarshalJSON(bytes)
	case fileTypePro:
		return unmarshalXML(bytes)
	}
	return nil, nil
}

func unmarshalJSON(bytes []byte) ([]Sheet, error) {
	var content []Sheet
	err := jsonAPI.Unmarshal(bytes, &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func unmarshalXML(bytes []byte) ([]Sheet, error) {
	var xmlContent xmlContent
	err := xml.Unmarshal(bytes, &xmlContent)
	if err != nil {
		return nil, err
	}

	sheets := make([]Sheet, 0, len(xmlContent.Sheets))
	for i := 0; i != len(xmlContent.Sheets); i++ {
		sheets = append(sheets, Sheet{
			Id:        xmlContent.Sheets[i].Id,
			Title:     xmlContent.Sheets[i].Title,
			RootTopic: xmlTopicToXmindTopic(xmlContent.Sheets[i].Topic),
		})
	}

	return sheets, nil
}

func xmlTopicToXmindTopic(z *XMLTopic) *XmindTopic {
	if z == nil {
		return nil
	}
	return &XmindTopic{
		Id:       z.Id,
		Title:    z.Title,
		Children: xmlChildrenToChildren(z.Children),
	}
}

func xmlChildrenToChildren(xmlChildren *XMLChildren) *Children {
	if xmlChildren == nil {
		return nil
	}

	attachedCount := 0
	detailsCount := 0

	for i := 0; i != len(xmlChildren.TypeTopicsList); i++ {
		if xmlChildren.TypeTopicsList[i].Type == "attached" {
			attachedCount++
		} else {
			detailsCount++
		}
	}

	children := &Children{
		Attached: make([]*XmindTopic, 0, attachedCount),
		Detached: make([]*XmindTopic, 0, detailsCount),
	}

	for i := 0; i != len(xmlChildren.TypeTopicsList); i++ {
		if xmlChildren.TypeTopicsList[i].Type == "attached" {
			children.Attached = append(children.Attached, xmlTypeTopicsToTopics(xmlChildren.TypeTopicsList[i].Topics)...)
		} else {
			children.Detached = append(children.Detached, xmlTypeTopicsToTopics(xmlChildren.TypeTopicsList[i].Topics)...)
		}
	}

	return children
}

func xmlTypeTopicsToTopics(typeTopics []*XMLTopic) []*XmindTopic {
	if len(typeTopics) == 0 {
		return nil
	}

	results := make([]*XmindTopic, 0, len(typeTopics))
	for i := 0; i != len(typeTopics); i++ {
		results = append(results, xmlTopicToXmindTopic(typeTopics[i]))
	}

	return results
}

func (s *XmindFile) ExtractAttached() []Topic {
	results := make([]Topic, 0, len(s.Sheets))
	for i := 0; i != len(s.Sheets); i++ {
		results = append(results, Topic{
			Id:     s.Sheets[i].Id,
			Title:  s.Sheets[i].Title,
			Topics: s.rootToTopic(s.Sheets[i].RootTopic),
		})
	}

	return results
}

func (s *XmindFile) rootToTopic(z *XmindTopic) []*Topic {
	if z == nil {
		return make([]*Topic, 0)
	}
	results := []*Topic{
		{
			Id:     z.Id,
			Title:  z.Title,
			Topics: s.childrenToTopic(z.Children),
		},
	}

	return results
}

func (s *XmindFile) childrenToTopic(z *Children) []*Topic {
	if z == nil {
		return make([]*Topic, 0)
	}
	results := make([]*Topic, 0, len(z.Attached))
	for i := 0; i != len(z.Attached); i++ {
		results = append(results, &Topic{
			Id:     z.Attached[i].Id,
			Title:  z.Attached[i].Title,
			Topics: s.childrenToTopic(z.Attached[i].Children),
		})
	}

	return results
}
