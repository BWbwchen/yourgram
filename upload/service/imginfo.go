package upload_svc

func (ii ImgInfo) Valid() bool {
	return ii.ImgID != "" && ii.ImgURL != "" && ii.User != ""
}
