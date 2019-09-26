package artifact

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(ctx context.Context, jobUrl string, artifactPath string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	bodyBytes := bytes.NewBuffer([]byte{})
	req, err := http.NewRequestWithContext(ctx, "GET", getStorageUrl(getArtifactsUrl(jobUrl))+"/"+strings.TrimPrefix(artifactPath, "/"), bodyBytes)
	if err != nil {
		return bodyBytes.Bytes(), err
	}
	response, err := tr.RoundTrip(req)
	if err != nil {
		return bodyBytes.Bytes(), err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func getArtifactsUrl(jobUrl string) string {
	return "https://gcsweb-ci.svc.ci.openshift.org" + strings.TrimSuffix(strings.TrimPrefix(jobUrl, "https://prow.svc.ci.openshift.org/view"), "/")
}

func getStorageUrl(artifactsUrl string) string {
	return "https://storage.googleapis.com" + strings.TrimSuffix(strings.TrimPrefix(artifactsUrl, "https://gcsweb-ci.svc.ci.openshift.org/gcs"), "/")
}
