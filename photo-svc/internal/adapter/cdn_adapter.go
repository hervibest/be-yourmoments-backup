package adapter

import (
	"fmt"

	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/utils"
)

type CDNAdapter interface {
	GenerateCDN(fileKey string) string
}

type CDNAdapterImpl struct {
	cdnBaseURL string
}

func NewCDNadapter() CDNAdapter {
	cdnBaseURL := utils.GetEnv("CDN_BASE_URL")
	return &CDNAdapterImpl{
		cdnBaseURL: cdnBaseURL,
	}
}

func (u *CDNAdapterImpl) GenerateCDN(fileKey string) string {
	return fmt.Sprintf("%s/?fileKey=%s", u.cdnBaseURL, fileKey)
}
