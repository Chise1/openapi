package openapi

type ContentType string

const (
	//Html     ContentType = "text/html"
	//Text     ContentType = "text/plain"
	//TextXml  ContentType = "text/xml"
	//Gif      ContentType = "image/gif"
	//Jpeg     ContentType = "image/jpeg"
	//Png      ContentType = "image/png"
	//AppXml   ContentType = "application/xml"
	Json ContentType = "application/json"
	//Pdf      ContentType = "application/pdf"
	//Stream   ContentType = "application/octet-stream"
	Form     ContentType = "application/x-www-form-urlencoded"
	FormData ContentType = "multipart/form-data"
)

type RouteStruct interface {
	GetReqPara() interface{}
	GetReqBody() interface{}         //interface可以为IContentType
	GetResBody() map[int]interface{} //interface可以为IContentType
	GetResPara() interface{}
}

// IBody body返回的数据结构
type IBody interface {
	Marshal() ([]byte, error)
}

// IContentType 返回的结构体在请求或返回参数的类型
type IContentType interface {
	GetContentType() ContentType
}
