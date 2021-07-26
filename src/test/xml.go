package test

import (
	"encoding/xml"
	"io/ioutil"
)

type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Value   string   `xml:",chardata"`
}
type InitParam struct {
	Props []Property `xml:"props>property"`
}

type Params struct {
	XMLName   xml.Name  `xml:"servlet"`
	InitParam InitParam `xml:"init-param"`
}

// func (p *Property) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	var s string
// 	if err := d.DecodeElement(&s, &start); err != nil {
// 		return err
// 	}
// 	fmt.Printf("d:%T\n", d)
// 	return nil
// }

func TestXml(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	params := Params{}
	err = xml.Unmarshal(content, &params)
	if err != nil {
		return err
	}
	return nil
}
