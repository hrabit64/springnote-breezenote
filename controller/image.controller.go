package controller

import (
	"encoding/json"
	"errors"
	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/database"
	"github.com/hrabit64/springnote-breezenote/dto"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	"github.com/hrabit64/springnote-breezenote/pkg/utils/img"
	validationUtil "github.com/hrabit64/springnote-breezenote/pkg/utils/validation"
	"github.com/hrabit64/springnote-breezenote/service"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

type ImageController interface {
	UploadImage(c *gin.Context)
	GetImage(c *gin.Context)
}

type imageController struct {
	itemService service.ItemService
}

func NewImageController(itemService service.ItemService) ImageController {
	return &imageController{itemService: itemService}
}

// UploadImage 이미지를 업로드합니다.
func (i *imageController) UploadImage(c *gin.Context) {
	var request dto.ImageCreateRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		i.handleError(err, c)

		return
	}

	if !validationUtil.CheckUploadImageExt(request.Name) {
		i.handleError(errors.New("not supported image type"), c)
		return
	}

	convName := uuid.New().String()

	if utils.GetFileExt(request.Name) == "gif" {

		err := i.saveGif(request.Image, convName)
		if err != nil {
			i.handleError(err, c)
			return
		}
	} else {
		//gif 가 아닌 경우 webp, jpeg로 변환
		err := i.saveNormalImage(request.Image, convName)
		if err != nil {
			i.handleError(err, c)
			return
		}
	}

	// 이미지 정보를 DB에 저장
	err = i.saveImgToDb(request.Name, convName)
	if err != nil {
		i.handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": convName})

}

// GetImage 이미지를 다운로드합니다.
func (i *imageController) GetImage(c *gin.Context) {
	id := c.Param("id")

	if !validationUtil.CheckDownloadImageExt(id) {
		i.handleError(errors.New("not supported download image type"), c)
		return
	}

	ext := utils.GetFileExt(id)

	// jpg -> jpeg
	if ext == "jpg" {
		id = utils.GetFileName(id) + ".jpeg"
		ext = "jpeg"
	}

	switch ext {
	case "webp":
		err := i.getWebp(id, c)
		if err != nil {
			i.handleError(err, c)
			return
		}
		break

	case "jpeg":
		err := i.getJpeg(id, c)
		if err != nil {
			i.handleError(err, c)
			return
		}
		break

	case "gif":
		err := i.getGif(id, c)
		if err != nil {
			i.handleError(err, c)
			return
		}
		break
	}

	c.Status(http.StatusOK)
}

// saveGif gif 이미지를 저장합니다.
func (i *imageController) saveGif(base64Img, convName string) error {
	err := img.SaveGif(base64Img, convName)
	if err != nil {
		return err
	}
	return nil
}

// saveNormalImage 이미지를 webp, jpeg로 변환하여 저장합니다.
func (i *imageController) saveNormalImage(base64Img, convName string) error {
	decodedImg := utils.DecodeBase64(base64Img)
	err := img.ConvertToWebp(decodedImg, config.RootConfig.MaxImageLen, convName)
	if err != nil {
		return err
	}

	decodedImg = utils.DecodeBase64(base64Img)
	err = img.ConvertToJpeg(decodedImg, config.RootConfig.MaxImageLen, convName)
	if err != nil {
		return err
	}

	return nil
}

// saveImgToDb 이미지 정보를 DB에 저장합니다.
func (i *imageController) saveImgToDb(originName, convName string) error {
	newItem := &database.Item{
		Id:         convName,
		OriginName: originName,
		CreateAt:   time.Now(),
	}

	err := i.itemService.CreateItem(newItem)
	if err != nil {
		return err
	}

	return nil
}

// getWebp webp 이미지를 반환합니다.
func (i *imageController) getWebp(name string, c *gin.Context) error {
	width, height, err := i.getWidthAndHeight(c)
	if err != nil {
		return err
	}

	webpImg, err := img.LoadImageWithResize(name, width, height)
	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "image/webp")
	err = webp.Encode(c.Writer, webpImg, nil)
	if err != nil {
		return err
	}

	return nil
}

// getJpeg jpeg 이미지를 반환합니다.
func (i *imageController) getJpeg(name string, c *gin.Context) error {
	width, height, err := i.getWidthAndHeight(c)
	if err != nil {
		return err
	}

	jpgImg, err := img.LoadImageWithResize(name, width, height)
	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(c.Writer, jpgImg, nil)
	if err != nil {
		return err
	}

	return nil
}

// getGif gif 이미지를 반환합니다.
func (i *imageController) getGif(name string, c *gin.Context) error {

	img, err := img.LoadGif(name)
	if err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Type", "image/gif")
	err = gif.EncodeAll(c.Writer, img)
	if err != nil {
		return err
	}

	return nil
}

// getWidthAndHeight width, height를 반환합니다. 이때 width, height가 없는 경우 0으로 반환합니다.
// 또한, width, height가 0보다 작거나 config.RootConfig.MaxImageLen보다 큰 경우 에러를 반환합니다.
func (i *imageController) getWidthAndHeight(c *gin.Context) (int, int, error) {
	qWidth := c.Query("width")
	qHeight := c.Query("height")

	if qWidth == "" {
		qWidth = "0"
	}

	if qHeight == "" {
		qHeight = "0"
	}

	width, err := strconv.Atoi(qWidth)
	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(qHeight)
	if err != nil {
		return 0, 0, err
	}

	if width > config.RootConfig.MaxImageLen || height > config.RootConfig.MaxImageLen {
		return 0, 0, errors.New("width and height is too big")
	}

	if width < 0 || height < 0 {
		return 0, 0, errors.New("width and height must be greater than 0")
	}

	return width, height, nil
}

func (i *imageController) handleError(err error, c *gin.Context) {

	if err == nil {
		return
	}

	//logger, _ := utils.GetLogger()
	path := c.Request.Host + c.Request.URL.Path

	var status int
	var title, detail string

	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	var file *os.PathError
	switch {

	case errors.As(err, &syntaxErr):
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "Request body contains badly-formed JSON"
		break

	case errors.As(err, &unmarshalErr):
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "Request body contains invalid JSON"
		break

	case errors.Is(err, image.ErrFormat):
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "Unsupported image type.  Supported types are jpeg, png, gif, webp"

	case err == io.EOF:
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "Request body must not be empty"
		break

	case utils.IsValidationError(err):
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = validationUtil.CreateValidationErrorDetail(err.(validator.ValidationErrors))
		break

	case errors.As(err, &file):
		status = http.StatusNotFound
		title = "Bad Request"
		detail = "Not Exist File"

	case err.Error() == "not supported image type":
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "Not supported image type (supported types are jpeg, gif, webp, png, jpg)"
		break

	case err.Error() == "not supported download image type":
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "Not supported download image type (supported types are jpeg(jpg), gif, webp)"
		break
	case err.Error() == "width and height is too big":
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "width and height must be less than" + strconv.Itoa(config.RootConfig.MaxImageLen)

	case err.Error() == "width and height must be greater than 0":
		status = http.StatusBadRequest
		title = "Bad Request"
		detail = "width and height must be greater than 0"

	default:
		log.Printf("unhandled error: %v %s", err, reflect.TypeOf(err))

		status = http.StatusInternalServerError
		title = "Internal Server Error"
		detail = "An internal server error occurred"
	}

	c.JSON(status, dto.ErrorResponse{
		Status:   status,
		Title:    title,
		Detail:   detail,
		Instance: path,
	})

	c.AbortWithStatus(status)

}
