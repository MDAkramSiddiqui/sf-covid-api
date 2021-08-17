package response_model

import "github.com/MDAkramSiddiqui/sf-covid-api/app/schema"

// default response builder for all messages
func DefaultResponse(status int, data interface{}) *schema.TDefaultResponse {
	var response *schema.TDefaultResponse
	if status >= 200 && status < 300 {
		response = &schema.TDefaultResponse{
			Status: "success",
			Data:   data,
		}
	} else {
		response = &schema.TDefaultResponse{
			Status: "failure",
			Data:   data,
		}
	}
	return response
}
