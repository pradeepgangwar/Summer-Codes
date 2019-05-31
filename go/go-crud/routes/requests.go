package route

import (
	"gopkg.in/go-playground/validator.v9"
)

// POSTRequest :
type POSTRequest struct {
	Name  string `bson:"name" json:"name"`
	Email string `json:"email" bson:"email" validate:"required"`
	About string `json:"about" bson:"about" validate:"required"`
}

// IsPostRequestValid ;
func (p *POSTRequest) IsPostRequestValid() (bool, error) {
	validate := validator.New()
	err := validate.Struct(p)

	if err != nil {
		return false, err
	}
	return true, nil
}
