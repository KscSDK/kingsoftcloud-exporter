package iam

import (
	"encoding/json"
	"fmt"

	"github.com/KscSDK/kingsoftcloud-exporter/config"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/iam"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var IAMProjectIDs = []int64{}

var OnlyIncludeProjectIDs = map[string][]int64{}

type ProjectMeta struct {
	AccountId   string `json:"AccountId"`
	ProjectId   int64  `json:"ProjectId"`
	ProjectName string `json:"ProjectName"`
	Status      int    `json:"Status"`
}

type IAMClient struct {
	credential config.Credential
	client     *iam.Iam
	logger     log.Logger
}

type ProjectList struct {
	Projects []*ProjectMeta `json:"ProjectList"`
	Total    int64          `json:"Total"`
}

type GetAccountAllProjectListResponse struct {
	ListProjectResult ProjectList `json:"ListProjectResult"`
	RequestId         string      `json:"RequestId"`
}

//GetAccountAllProjectList
func (c *IAMClient) GetAccountAllProjectList() ([]int64, error) {

	level.Info(c.logger).Log("msg", "IAM Project 开始加载")

	resp, err := c.client.GetAccountAllProjectList(nil)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("get no iam project list.")
	}

	respBytes, _ := json.Marshal(resp)

	var response GetAccountAllProjectListResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, fmt.Errorf("parse iam project list err, %+v", err)
	}

	iamProjectIDs := make([]int64, 0, len(response.ListProjectResult.Projects))
	for _, v := range response.ListProjectResult.Projects {
		iamProjectIDs = append(iamProjectIDs, v.ProjectId)
	}

	level.Info(c.logger).Log("msg", "IAM Project 资源加载完毕", "project_num", len(iamProjectIDs))

	return iamProjectIDs, nil
}

func ReloadIAMProjects(conf *config.KscExporterConfig, logger log.Logger) (err error) {
	c, err := NewKscIAMClient(conf, logger)
	if err != nil {
		return err
	}

	for i := 0; i < len(conf.Products); i++ {
		if len(conf.Products[i].OnlyIncludeProjects) > 0 {
			if _, isOK := OnlyIncludeProjectIDs[conf.Products[i].Namespace]; !isOK {
				OnlyIncludeProjectIDs[conf.Products[i].Namespace] = conf.Products[i].OnlyIncludeProjects
			}
		}
	}

	IAMProjectIDs, err = c.GetAccountAllProjectList()
	if err != nil {
		return err
	}

	return nil
}

//NewKscIAMClient
func NewKscIAMClient(conf *config.KscExporterConfig, logger log.Logger) (*IAMClient, error) {
	svc := iam.SdkNew(
		ksc.NewClient(conf.Credential.AccessKey, conf.Credential.SecretKey),
		&ksc.Config{Region: &conf.Credential.Region},
		&utils.UrlInfo{
			UseSSL: conf.Credential.UseSSL,
		},
	)

	repo := &IAMClient{
		credential: conf.Credential,
		client:     svc,
		logger:     logger,
	}

	return repo, nil
}
