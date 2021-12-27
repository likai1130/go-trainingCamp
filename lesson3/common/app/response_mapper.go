package app

import (
	"github.com/gin-gonic/gin"
	"go-trainingCamp/lesson3/common/e"
	"go-trainingCamp/lesson3/i18n"
)

type Response struct {
	Code      string      `json:"code"`
	Msg       string      `json:"msg"`
	ErrDetail string      `json:"err_detail"`
	Data      interface{} `json:"data"`
}

type ResponseI18nMsgParams struct {
	C            *gin.Context `json:"c"`
	Code         string       `json:"code"`
	Data         interface{}  `json:"data"`
	Err          error        `json:"err"`
	TemplateData interface{}  `json:"template_data"`
	PluralCount  interface{}  `json:"plural_count"`
}

func NewResponse(respi18n ResponseI18nMsgParams,c *gin.Context) {
	code := respi18n.Code
	httpCode := getHttpCode(code)
	templateData := respi18n.TemplateData
	pluralCount := respi18n.PluralCount
	data := respi18n.Data
	errDetail := ""
	if respi18n.Err != nil {
		errDetail = respi18n.Err.Error()
	}

	if code == e.ScodeOK {
		c.JSON(httpCode, Response{
			Code:      code,
			Msg:       "",
			ErrDetail: "",
			Data:      data,
		})
		return
	}

	lang := c.Request.Header.Get("lang")
	accept := c.Request.Header.Get("Accept-Language")
	msg := i18n.MustLocalize(lang, accept, code, templateData, pluralCount)
	c.JSON(httpCode, Response{
		Code:      code,
		Msg:       msg,
		ErrDetail: errDetail,
		Data:      data,
	})
	return
}

func getHttpCode(code string) int {
	msg, ok := e.MsgFlags[code]
	if ok {
		return msg
	}

	return e.MsgFlags[e.InternalError]
}

func NewI18nResponse(code string, data interface{}, err error) ResponseI18nMsgParams {
	return ResponseI18nMsgParams{
		Code: code,
		Err:  err,
		Data: data,
	}
}
