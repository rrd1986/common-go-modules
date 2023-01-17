package eemi

type Response struct {
	MessageID      string `json:"messageId"`
	Message        string `json:"message"`
	ResponseAction string `json:"responseAction"`
	Category       string `json:"category"`
	SeverityLevel  string `json:"severity"`
}

// CreateEEMIError create a new instance of Error
func NewEemiResponse(messageID2 string, message2 string, responseAction2 string, category2 string, severityLevel2 string) Response {
	return Response{MessageID: messageID2, Message: message2, ResponseAction: responseAction2, Category: category2, SeverityLevel: severityLevel2}
}
