package gendocx

import (
	"bytes"
	"github.com/lukasjarosch/go-docx"
)

func GenDocFromBytes(file []byte, replaceMap map[string]any) ([]byte, error) {
	doc, err := docx.OpenBytes(file)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	err = doc.Write(&buff)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
