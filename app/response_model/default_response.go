package response_model

import (
	"os"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
	"github.com/MDAkramSiddiqui/sf-covid-api/app/schema"
)

// default response builder for all messages
func DefaultResponse(status int, data interface{}, isForcedVisible bool) (int, *schema.TDefaultResponse) {
	var response *schema.TDefaultResponse
	if status >= 200 && status < 300 {
		response = &schema.TDefaultResponse{
			Status: "success",
			Data:   data,
		}
	} else {
		response = &schema.TDefaultResponse{
			Status:  "failure",
			Message: data,
		}

		if os.Getenv(constants.Env) == constants.Production && !isForcedVisible {
			response.Message = "Oops, something bad happened!"
		}
	}
	return status, response
}
