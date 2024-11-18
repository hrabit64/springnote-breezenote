package dto

type ImageCreateRequest struct {
	Image string `json:"image" binding:"required,base64,max=5592405"`
	Name  string `json:"name" binding:"required,min=1,max=255"`
}
