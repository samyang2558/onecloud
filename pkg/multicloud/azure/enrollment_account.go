// Copyright 2019 Yunion
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

package azure

import (
	"fmt"
	"net/url"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/cloudprovider"
)

type SEnrollmentAccountProperties struct {
	PrincipalName string
	OfferTypes    []string
}

type SEnrollmentAccount struct {
	Id         string
	Name       string
	Type       string
	Properties SEnrollmentAccountProperties
}

func (cli *SAzureClient) GetEnrollmentAccounts() ([]cloudprovider.SEnrollmentAccount, error) {
	accounts := struct {
		Value []SEnrollmentAccount
	}{}
	err := cli.Get("providers/Microsoft.Billing/enrollmentAccounts", nil, &accounts)
	if err != nil {
		return nil, err
	}
	eas := []cloudprovider.SEnrollmentAccount{}
	for i := range accounts.Value {
		ea := cloudprovider.SEnrollmentAccount{
			Id:   accounts.Value[i].Name,
			Name: accounts.Value[i].Properties.PrincipalName,
			Type: accounts.Value[i].Type,
		}
		eas = append(eas, ea)
	}
	return eas, nil
}

func (cli *SAzureClient) CreateSubscription(name string, eaId string, offerType string) error {
	appId, err := cli.GetAppObjectId()
	if err != nil {
		log.Errorf("GetAppObjectId error: %v", err)
	}
	owners := []map[string]string{}
	if len(appId) > 0 {
		owners = append(owners, map[string]string{"objectId": appId})
	}
	body := map[string]interface{}{
		"displayName": name,
		"offerType":   offerType,
		"owners":      owners,
	}
	resource := fmt.Sprintf("providers/Microsoft.Billing/enrollmentAccounts/%s/providers/Microsoft.Subscription/createSubscription", eaId)
	return cli.POST(resource, jsonutils.Marshal(body))
}

type SServicePrincipal struct {
	AppId    string
	ObjectId string
}

func (cli *SAzureClient) GetAppObjectId() (string, error) {
	result, err := cli.ListServicePrincipal(cli.clientId)
	if err != nil {
		return "", errors.Wrap(err, "ListServicePrincipal")
	}
	if len(result) == 1 {
		return result[0].ObjectId, nil
	}
	return "", nil
}

func (cli *SAzureClient) ListServicePrincipal(appId string) ([]SServicePrincipal, error) {
	params := url.Values{}
	if len(appId) > 0 {
		params.Set("$filter", fmt.Sprintf(`appId eq '%s'`, cli.clientId))
	}
	result := []SServicePrincipal{}
	err := cli.ListGraphResource("servicePrincipals", params, &result)
	if err != nil {
		return result, errors.Wrap(err, "ListGraphResource.servicePrincipals")
	}
	return result, nil
}