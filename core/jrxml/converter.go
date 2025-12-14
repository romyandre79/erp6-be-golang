package jrxml

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// JRXMLReport represents the root JRXML structure
type JRXMLReport struct {
	XMLName      xml.Name   `xml:"jasperReport"`
	Name         string     `xml:"name,attr"`
	PageWidth    int        `xml:"pageWidth,attr"`
	PageHeight   int        `xml:"pageHeight,attr"`
	Orientation  string     `xml:"orientation,attr,omitempty"`
	ColumnWidth  int        `xml:"columnWidth,attr,omitempty"`
	LeftMargin   int        `xml:"leftMargin,attr,omitempty"`
	RightMargin  int        `xml:"rightMargin,attr,omitempty"`
	TopMargin    int        `xml:"topMargin,attr,omitempty"`
	BottomMargin int        `xml:"bottomMargin,attr,omitempty"`
	Title        *JRXMLBand `xml:"title,omitempty"`
	PageHeader   *JRXMLBand `xml:"pageHeader,omitempty"`
	ColumnHeader *JRXMLBand `xml:"columnHeader,omitempty"`
	Detail       *JRXMLBand `xml:"detail,omitempty"`
	ColumnFooter *JRXMLBand `xml:"columnFooter,omitempty"`
	PageFooter   *JRXMLBand `xml:"pageFooter,omitempty"`
	Summary      *JRXMLBand `xml:"summary,omitempty"`
}

// JRXMLBand represents report bands (title, pageHeader, detail, etc.)
type JRXMLBand struct {
	Height      int               `xml:"height,attr"`
	StaticTexts []JRXMLStaticText `xml:"staticText,omitempty"`
	TextFields  []JRXMLTextField  `xml:"textField,omitempty"`
	Images      []JRXMLImage      `xml:"image,omitempty"`
	Lines       []JRXMLLine       `xml:"line,omitempty"`
	Rectangles  []JRXMLRectangle  `xml:"rectangle,omitempty"`
}

// JRXMLStaticText represents static text elements
type JRXMLStaticText struct {
	ReportElement JRXMLReportElement `xml:"reportElement"`
	Text          string             `xml:"text"`
}

// JRXMLTextField represents dynamic text field elements
type JRXMLTextField struct {
	ReportElement JRXMLReportElement `xml:"reportElement"`
	Expression    string             `xml:"textFieldExpression"`
}

// JRXMLImage represents image elements
type JRXMLImage struct {
	ReportElement JRXMLReportElement `xml:"reportElement"`
	Expression    string             `xml:"imageExpression"`
}

// JRXMLLine represents line elements
type JRXMLLine struct {
	ReportElement JRXMLReportElement `xml:"reportElement"`
}

// JRXMLRectangle represents rectangle elements
type JRXMLRectangle struct {
	ReportElement JRXMLReportElement `xml:"reportElement"`
}

// JRXMLReportElement represents the common reportElement tag
type JRXMLReportElement struct {
	X      int `xml:"x,attr"`
	Y      int `xml:"y,attr"`
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

// Report represents our internal JSON format
type Report struct {
	PageWidth   int           `json:"pageWidth"`
	PageHeight  int           `json:"pageHeight"`
	Orientation string        `json:"orientation"`
	Margins     ReportMargins `json:"margins"`
	Bands       []ReportBand  `json:"bands"`
}

// ReportMargins represents page margins
type ReportMargins struct {
	Left   int `json:"left"`
	Right  int `json:"right"`
	Top    int `json:"top"`
	Bottom int `json:"bottom"`
}

// ReportBand represents a report band
type ReportBand struct {
	Type     string          `json:"type"` // title, pageHeader, columnHeader, detail, etc.
	Height   int             `json:"height"`
	Elements []ReportElement `json:"elements"`
}

// ReportElement represents a report element
type ReportElement struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"` // staticText, textField, image, line, rectangle, etc.
	X          int                    `json:"x"`
	Y          int                    `json:"y"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Properties map[string]interface{} `json:"properties"`
}

// ConvertJRXMLToJSON converts JRXML to our internal JSON format
func ConvertJRXMLToJSON(jrxmlContent string) (string, error) {
	var jrxml JRXMLReport
	if err := xml.Unmarshal([]byte(jrxmlContent), &jrxml); err != nil {
		return "", fmt.Errorf("failed to parse JRXML: %w", err)
	}

	// Convert JRXML structure to our Report
	template := convertJRXMLToTemplate(&jrxml)

	jsonBytes, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// ConvertJSONToJRXML converts our JSON format to JRXML
func ConvertJSONToJRXML(jsonContent string) (string, error) {
	var template Report
	if err := json.Unmarshal([]byte(jsonContent), &template); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert our template to JRXML structure
	jrxml := convertTemplateToJRXML(&template)

	xmlBytes, err := xml.MarshalIndent(jrxml, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JRXML: %w", err)
	}

	return xml.Header + string(xmlBytes), nil
}

// convertJRXMLToTemplate converts JRXML structure to our internal format
func convertJRXMLToTemplate(jrxml *JRXMLReport) *Report {
	template := &Report{
		PageWidth:   jrxml.PageWidth,
		PageHeight:  jrxml.PageHeight,
		Orientation: jrxml.Orientation,
		Margins: ReportMargins{
			Left:   jrxml.LeftMargin,
			Right:  jrxml.RightMargin,
			Top:    jrxml.TopMargin,
			Bottom: jrxml.BottomMargin,
		},
		Bands: []ReportBand{},
	}

	// Convert each band
	if jrxml.Title != nil {
		template.Bands = append(template.Bands, convertBand("title", jrxml.Title))
	}
	if jrxml.PageHeader != nil {
		template.Bands = append(template.Bands, convertBand("pageHeader", jrxml.PageHeader))
	}
	if jrxml.ColumnHeader != nil {
		template.Bands = append(template.Bands, convertBand("columnHeader", jrxml.ColumnHeader))
	}
	if jrxml.Detail != nil {
		template.Bands = append(template.Bands, convertBand("detail", jrxml.Detail))
	}
	if jrxml.ColumnFooter != nil {
		template.Bands = append(template.Bands, convertBand("columnFooter", jrxml.ColumnFooter))
	}
	if jrxml.PageFooter != nil {
		template.Bands = append(template.Bands, convertBand("pageFooter", jrxml.PageFooter))
	}
	if jrxml.Summary != nil {
		template.Bands = append(template.Bands, convertBand("summary", jrxml.Summary))
	}

	return template
}

// convertBand converts a JRXML band to our internal format
func convertBand(bandType string, jrxmlBand *JRXMLBand) ReportBand {
	band := ReportBand{
		Type:     bandType,
		Height:   jrxmlBand.Height,
		Elements: []ReportElement{},
	}

	// Convert static texts
	for i, st := range jrxmlBand.StaticTexts {
		band.Elements = append(band.Elements, ReportElement{
			ID:     fmt.Sprintf("%s_staticText_%d", bandType, i),
			Type:   "staticText",
			X:      st.ReportElement.X,
			Y:      st.ReportElement.Y,
			Width:  st.ReportElement.Width,
			Height: st.ReportElement.Height,
			Properties: map[string]interface{}{
				"text": st.Text,
			},
		})
	}

	// Convert text fields
	for i, tf := range jrxmlBand.TextFields {
		band.Elements = append(band.Elements, ReportElement{
			ID:     fmt.Sprintf("%s_textField_%d", bandType, i),
			Type:   "textField",
			X:      tf.ReportElement.X,
			Y:      tf.ReportElement.Y,
			Width:  tf.ReportElement.Width,
			Height: tf.ReportElement.Height,
			Properties: map[string]interface{}{
				"expression": tf.Expression,
			},
		})
	}

	// Convert images
	for i, img := range jrxmlBand.Images {
		band.Elements = append(band.Elements, ReportElement{
			ID:     fmt.Sprintf("%s_image_%d", bandType, i),
			Type:   "image",
			X:      img.ReportElement.X,
			Y:      img.ReportElement.Y,
			Width:  img.ReportElement.Width,
			Height: img.ReportElement.Height,
			Properties: map[string]interface{}{
				"expression": img.Expression,
			},
		})
	}

	// Convert lines
	for i, line := range jrxmlBand.Lines {
		band.Elements = append(band.Elements, ReportElement{
			ID:         fmt.Sprintf("%s_line_%d", bandType, i),
			Type:       "line",
			X:          line.ReportElement.X,
			Y:          line.ReportElement.Y,
			Width:      line.ReportElement.Width,
			Height:     line.ReportElement.Height,
			Properties: map[string]interface{}{},
		})
	}

	// Convert rectangles
	for i, rect := range jrxmlBand.Rectangles {
		band.Elements = append(band.Elements, ReportElement{
			ID:         fmt.Sprintf("%s_rectangle_%d", bandType, i),
			Type:       "rectangle",
			X:          rect.ReportElement.X,
			Y:          rect.ReportElement.Y,
			Width:      rect.ReportElement.Width,
			Height:     rect.ReportElement.Height,
			Properties: map[string]interface{}{},
		})
	}

	return band
}

// convertTemplateToJRXML converts our internal format to JRXML structure
func convertTemplateToJRXML(template *Report) *JRXMLReport {
	jrxml := &JRXMLReport{
		Name:         "Report",
		PageWidth:    template.PageWidth,
		PageHeight:   template.PageHeight,
		Orientation:  template.Orientation,
		ColumnWidth:  template.PageWidth - template.Margins.Left - template.Margins.Right,
		LeftMargin:   template.Margins.Left,
		RightMargin:  template.Margins.Right,
		TopMargin:    template.Margins.Top,
		BottomMargin: template.Margins.Bottom,
	}

	// Convert each band
	for _, band := range template.Bands {
		jrxmlBand := convertBandToJRXML(&band)

		switch strings.ToLower(band.Type) {
		case "title":
			jrxml.Title = jrxmlBand
		case "pageheader":
			jrxml.PageHeader = jrxmlBand
		case "columnheader":
			jrxml.ColumnHeader = jrxmlBand
		case "detail":
			jrxml.Detail = jrxmlBand
		case "columnfooter":
			jrxml.ColumnFooter = jrxmlBand
		case "pagefooter":
			jrxml.PageFooter = jrxmlBand
		case "summary":
			jrxml.Summary = jrxmlBand
		}
	}

	return jrxml
}

// convertBandToJRXML converts our internal band to JRXML band
func convertBandToJRXML(band *ReportBand) *JRXMLBand {
	jrxmlBand := &JRXMLBand{
		Height: band.Height,
	}

	for _, elem := range band.Elements {
		switch elem.Type {
		case "staticText":
			text, _ := elem.Properties["text"].(string)
			jrxmlBand.StaticTexts = append(jrxmlBand.StaticTexts, JRXMLStaticText{
				ReportElement: JRXMLReportElement{
					X:      elem.X,
					Y:      elem.Y,
					Width:  elem.Width,
					Height: elem.Height,
				},
				Text: text,
			})
		case "textField":
			expr, _ := elem.Properties["expression"].(string)
			jrxmlBand.TextFields = append(jrxmlBand.TextFields, JRXMLTextField{
				ReportElement: JRXMLReportElement{
					X:      elem.X,
					Y:      elem.Y,
					Width:  elem.Width,
					Height: elem.Height,
				},
				Expression: expr,
			})
		case "image":
			expr, _ := elem.Properties["expression"].(string)
			jrxmlBand.Images = append(jrxmlBand.Images, JRXMLImage{
				ReportElement: JRXMLReportElement{
					X:      elem.X,
					Y:      elem.Y,
					Width:  elem.Width,
					Height: elem.Height,
				},
				Expression: expr,
			})
		case "line":
			jrxmlBand.Lines = append(jrxmlBand.Lines, JRXMLLine{
				ReportElement: JRXMLReportElement{
					X:      elem.X,
					Y:      elem.Y,
					Width:  elem.Width,
					Height: elem.Height,
				},
			})
		case "rectangle":
			jrxmlBand.Rectangles = append(jrxmlBand.Rectangles, JRXMLRectangle{
				ReportElement: JRXMLReportElement{
					X:      elem.X,
					Y:      elem.Y,
					Width:  elem.Width,
					Height: elem.Height,
				},
			})
		}
	}

	return jrxmlBand
}
