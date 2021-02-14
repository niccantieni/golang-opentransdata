package opentransdata

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	apiURL = "https://api.opentransportdata.swiss/trias2020"
	//Used for timestamp formatting
	ShortRFC3339 = "2006-01-02T15:04:05Z"
)

//OTDRequest is a struct to store the needed data for the request creation
type OTDRequest struct {
	Timestamp    string
	StopPointRef string
	DepArrTime   string
	Parameters   OTDParameters
}

//OTDParameters stores the request-parameters of an OTDRequest
type OTDParameters struct {
	NumberOfResults      string
	StopEventType        string
	IncludePreviousCalls bool
	IncludeOnwardCalls   bool
	IncludeRealtimeData  bool
}

//NewOTDRequest is a constructor for OTDRequest
func NewOTDRequest(timestamp string, stopPointRef string, depArrTime string, NumberOfResults string, StopEventType string,
	IncludePreviousCalls bool, IncludeOnwardCalls bool, IncludeRealtimeData bool) (request OTDRequest) {
	request = OTDRequest{
		Timestamp:    timestamp,
		StopPointRef: stopPointRef,
		DepArrTime:   depArrTime,
		Parameters: OTDParameters{
			NumberOfResults:      NumberOfResults,
			StopEventType:        StopEventType,
			IncludePreviousCalls: IncludePreviousCalls,
			IncludeOnwardCalls:   IncludeOnwardCalls,
			IncludeRealtimeData:  IncludeRealtimeData,
		},
	}
	return request
}

//TemplateOTDRequestNow creates a OTDRequest for general use, with somewhat useful settings;
//depArrTime = now, no previous and onward calls, with live timetable,
//stationID needs to be set afterwards
func TemplateOTDRequestNow() (request OTDRequest) {
	//Load Location in Zürich, everything else does not make much sense
	zurich, _ := time.LoadLocation("Europe/Zurich")

	//current time in Zurich
	now := time.Now().In(zurich)

	//format timestamp as timestamp correctly; OTD interprets this as localtime (somehow, but not really?!?)
	//whatever, it works like this to get the current (as in right now, instant) events.
	depArrTime := now.Format(ShortRFC3339)

	request = NewOTDRequest("", "", depArrTime, "1", "departure", false, false, true)

	return request
}

//CreateRequest creates a HTTP-Request to OpenTransportData with the given input and API-Key
func CreateRequest(OTDApiKey string, request OTDRequest) (data []byte, err error) {
	//get time
	now := time.Now()

	//create Timestamp format
	request.Timestamp = now.Format(ShortRFC3339)

	//create XML Request
	XMLRequest := CreateXML(request)

	//create reader for body
	body := strings.NewReader(XMLRequest)

	//Create request
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return data, err
	}

	//add content-type and authorization headers
	req.Header.Add("Content-Type", "text/XML")
	req.Header.Add("Authorization", OTDApiKey)

	//execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}

	//read from response
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	//close reader
	err = res.Body.Close()

	//return
	return data, err
}

//CreateXML creates a string containing the appropriate XML from the given input
func CreateXML(request OTDRequest) (XMLRequest string) {

	//puzzle together the XML
	XMLRequest = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Trias version=\"1.1\" xmlns=\"http://www.vdv.de/trias\" xmlns:siri=\"http://www.siri.org.uk/siri\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">\n<ServiceRequest>\n<siri:RequestTimestamp>" + request.Timestamp +
		"</siri:RequestTimestamp>\n<siri:RequestorRef>NicCantieni</siri:RequestorRef>\n<RequestPayload>\n<StopEventRequest>\n<Location>\n<LocationRef>\n<StopPointRef>" +
		request.StopPointRef + "</StopPointRef>\n</LocationRef>\n<DepArrTime>" +
		request.DepArrTime + "</DepArrTime>\n</Location>\n<Params>\n<NumberOfResults>" +
		request.Parameters.NumberOfResults + "</NumberOfResults>\n<StopEventType>" +
		request.Parameters.StopEventType + "</StopEventType>\n<IncludePreviousCalls>" +
		strconv.FormatBool(request.Parameters.IncludePreviousCalls) + "</IncludePreviousCalls>\n<IncludeOnwardCalls>" +
		strconv.FormatBool(request.Parameters.IncludeOnwardCalls) + "</IncludeOnwardCalls>\n<IncludeRealtimeData>" +
		strconv.FormatBool(request.Parameters.IncludeRealtimeData) + "</IncludeRealtimeData>\n</Params>\n</StopEventRequest>\n</RequestPayload>\n</ServiceRequest>\n</Trias>"

	//return
	return XMLRequest
}
