// Copyright 2020 FairwindsOps Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package insights

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fairwindsops/insights-plugins/opa/pkg/opa"
	"github.com/imroc/req"
	"github.com/sirupsen/logrus"

	"github.com/fairwindsops/insights-cli/pkg/models"
)

const opaURLFormat = "%s/v0/organizations/%s/opa/customChecks"

const opaCheckURLFormat = opaURLFormat + "/%s"
const opaCheckInstancesURLFormat = opaCheckURLFormat + "/instances"

const opaInstanceURLFormat = opaCheckInstancesURLFormat + "/%s"

// GetChecks queries Fairwinds Insights to retrieve all of the Checks for an organization
func GetChecks(org, token, hostName string) ([]opa.OPACustomCheck, error) {
	url := fmt.Sprintf(opaURLFormat, hostName, org)
	logrus.Infof("Url: %s", url)
	resp, err := req.Get(url, getHeaders(token))
	if err != nil {
		return nil, err
	}
	var checks []opa.OPACustomCheck
	if resp.Response().StatusCode != http.StatusOK {
		logrus.Errorf("Invalid response code: %s %v", string(resp.Bytes()), resp.Response().StatusCode)
		return nil, errors.New("invalid response code")
	}
	err = resp.ToJSON(&checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

// GetInstances queries Fairwinds Insights to retrieve all of the instances for a given check
func GetInstances(org, checkName, token, hostName string) ([]opa.CheckSetting, error) {
	url := fmt.Sprintf(opaCheckInstancesURLFormat, hostName, org, checkName)
	resp, err := req.Get(url, getHeaders(token))
	if err != nil {
		return nil, err
	}
	if resp.Response().StatusCode != http.StatusOK {
		logrus.Errorf("Invalid response code: %s %v", string(resp.Bytes()), resp.Response().StatusCode)
		return nil, errors.New("invalid response code")
	}
	var instances []opa.CheckSetting
	err = resp.ToJSON(&instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// DeleteCheck deletes an OPA Check from Fairwinds Insights
func DeleteCheck(check models.CustomCheckModel, org, token, hostName string) error {
	url := fmt.Sprintf(opaCheckURLFormat, hostName, org, check.CheckName)
	resp, err := req.Delete(url, getHeaders(token))
	if err != nil {
		return err
	}
	if resp.Response().StatusCode != http.StatusOK {
		logrus.Errorf("Invalid response code: %s %v", string(resp.Bytes()), resp.Response().StatusCode)
		return errors.New("invalid response code")
	}
	return nil
}

// PutCheck upserts an OPA Check to Fairwinds Insights
func PutCheck(check models.CustomCheckModel, org, token, hostName string) error {
	url := fmt.Sprintf(opaCheckURLFormat, hostName, org, check.CheckName)
	resp, err := req.Put(url, getHeaders(token), req.BodyJSON(&check))
	if err != nil {
		return err
	}
	if resp.Response().StatusCode != http.StatusOK {
		logrus.Errorf("Invalid response code: %s %v", string(resp.Bytes()), resp.Response().StatusCode)
		return errors.New("invalid response code")
	}
	return nil
}

// DeleteInstance deletes an Instance from Fairwinds Insights
func DeleteInstance(instance models.CustomCheckInstanceModel, org, token, hostName string) error {
	url := fmt.Sprintf(opaInstanceURLFormat, hostName, org, instance.CheckName, instance.InstanceName)
	resp, err := req.Delete(url, getHeaders(token))
	if err != nil {
		return err
	}
	if resp.Response().StatusCode != http.StatusOK {
		logrus.Errorf("Invalid response code: %s %v", string(resp.Bytes()), resp.Response().StatusCode)
		return errors.New("invalid response code")
	}
	return nil
}

// PutInstance upserts an Instance to Fairwinds Insights
func PutInstance(instance models.CustomCheckInstanceModel, org, token, hostName string) error {
	url := fmt.Sprintf(opaInstanceURLFormat, hostName, org, instance.CheckName, instance.InstanceName)
	resp, err := req.Put(url, getHeaders(token), req.BodyJSON(&instance))
	if err != nil {
		return err
	}
	if resp.Response().StatusCode != http.StatusOK {
		logrus.Errorf("Invalid response code: %s %v", string(resp.Bytes()), resp.Response().StatusCode)
		return errors.New("invalid response code")
	}
	return nil
}

func getHeaders(token string) req.Header {
	return req.Header{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Accept":        "application/json",
	}
}
