package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" `
	ResourceID string    `json:"resource_id" valid:"notnull" `
	FilePath   string    `json:"file_path" valid:"notnull" `
	CreatedAt  time.Time `json:"-" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo() *Video {
	return &Video{}
}

func (video *Video) Validate() error {

	_, err := govalidator.ValidateStruct(video)

	if err != nil {
		return err
	}

	return nil
}
