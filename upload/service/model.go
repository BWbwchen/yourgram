package upload_svc

type UploadRequest struct {
	UserID string `json:"userid"`
	Data   []byte `json:"data"`
	ImgID  string `json:"imgid"`
}

type UploadResponse struct {
	StatusCode int     `json:"status"`
	Info       ImgInfo `json:"info"`
}

type ImgInfo struct {
	User   string `json:"user"`
	ImgID  string `json:"imgid"`
	ImgURL string `json:"imgurl"`
}
