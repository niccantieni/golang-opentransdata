package opentransdata

import (
	"encoding/xml"
	"time"
)

type Text struct {
	Text     string `xml:"Text"`
	Language string `xml:"Language"`
}

//Timetable makes the mapping for ServiceArrival and -Departure to a somewhat useful time
type Timetable struct {
	TimetabledTime Time `xml:"TimetabledTime"`
	EstimatedTime  Time `xml:"EstimatedTime"`
}

//Time is used to map to a "real" golang time.Time. At the moment, just the string is parsed
//from the XML, the automatic conversion to something useful is not implemented yet.
type Time struct {
	Time   time.Time
	String string `xml:",chardata"`
}

type Trias struct {
	XMLName         xml.Name        `xml:"Trias"`
	SiriNS          string          `xml:"siri,attr"`
	TriasNS         string          `xml:"trias,attr"`
	AcsbURl         string          `xml:"acsb,attr"`
	IfoptURl        string          `xml:"ifopt,attr"`
	Datex2URl       string          `xml:"datex2,attr"`
	Version         string          `xml:"version,attr"`
	ServiceDelivery ServiceDelivery `xml:"ServiceDelivery"`
}

type ServiceDelivery struct {
	XMLName           xml.Name        `xml:"ServiceDelivery"`
	ResponseTimestamp Time            `xml:"ResponseTimestamp"`
	ProducerRef       string          `xml:"ProducerRef"`
	Status            bool            `xml:"Status"`
	Language          string          `xml:"Language"`
	CalcTime          int             `xml:"CalcTime"`
	DeliveryPayload   DeliveryPayload `xml:"DeliveryPayload"`
}

type DeliveryPayload struct {
	XMLName           xml.Name          `xml:"DeliveryPayload"`
	StopEventResponse StopEventResponse `xml:"StopEventResponse"`
}

type StopEventResponse struct {
	XMLName                  xml.Name                 `xml:"StopEventResponse"`
	ErrorMessage             ErrorMessage             `xml:"ErrorMessage"`
	StopEventResponseContext StopEventResponseContext `xml:"StopEventResponseContext"`
	StopEventResult          []StopEventResult        `xml:"StopEventResult"`
}

type ErrorMessage struct {
	XMLName xml.Name `xml:"ErrorMessage"`
	Code    string   `xml:"Code"`
	Text    Text     `xml:"Text"`
}

type StopEventResponseContext struct {
	XMLName    xml.Name   `xml:"StopEventResponseContext"`
	Situations Situations `xml:"Situations"`
}

type Situations struct {
	XMLName    xml.Name `xml:"Situations"`
	Situations string   `xml:",chardata"`
}

type StopEventResult struct {
	XMLName   xml.Name  `xml:"StopEventResult"`
	ResultId  string    `xml:"ResultId"`
	StopEvent StopEvent `xml:"StopEvent"`
}

type StopEvent struct {
	XMLName      xml.Name `xml:"StopEvent"`
	PreviousCall []Call   `xml:"PreviousCall"`
	ThisCall     Call     `xml:"ThisCall"`
	OnwardCall   []Call   `xml:"OnwardCall"`
	Service      Service  `xml:"Service"`
}

type Call struct {
	CallAtStop CallAtStop `xml:"CallAtStop"`
}

type Service struct {
	XMLName                 xml.Name    `xml:"Service"`
	OperatingDayRef         string      `xml:"OperatingDayRef"`
	JourneyRef              string      `xml:"JourneyRef"`
	LineRef                 string      `xml:"LineRef"`
	DirectionRef            string      `xml:"DirectionRef"`
	Mode                    Mode        `xml:"Mode"`
	PublishedLineName       Text        `xml:"PublishedLineName"`
	OperatorRef             string      `xml:"OperatorRef"`
	OriginStopPointRef      string      `xml:"OriginStopPointRef"`
	OriginText              Text        `xml:"OriginText"`
	DestinationStopPointRef string      `xml:"DestinationStopPointRef"`
	DestinationText         Text        `xml:"DestinationText"`
	Attribute               []Attribute `xml:"Attribute"`
}

type CallAtStop struct {
	XMLName          xml.Name  `xml:"CallAtStop"`
	StopPointRef     string    `xml:"StopPointRef"`
	StopPointName    Text      `xml:"StopPointName"`
	PlannedBay       Text      `xml:"PlannedBay"`
	EstimatedBay     Text      `xml:"EstimatedBay"`
	ServiceArrival   Timetable `xml:"ServiceArrival"`
	ServiceDeparture Timetable `xml:"ServiceDeparture"`
	StopSeqNumber    int       `xml:"StopSeqNumber"`
}

type Mode struct {
	XMLName     xml.Name `xml:"Mode"`
	PtMode      string   `xml:"PtMode"`
	RailSubmode string   `xml:"RailSubmode"`
	Name        Text     `xml:"Name"`
}

type Attribute struct {
	XMLName xml.Name `xml:"Attribute"`
	Text    Text     `xml:"Text"`
	Code    string   `xml:"Code"`
}

func ParseXML(input []byte) (out Trias, err error) {
	var parsed Trias
	err = xml.Unmarshal(input, &parsed)

	return parsed, err
}
