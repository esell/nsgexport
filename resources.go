package main

import "encoding/json"

type AuthTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	ExtExpiresIn string `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

type ARMTemplate struct {
	Schema         string               `json:"$schema"`
	ContentVersion string               `json:"contentVersion"`
	Parameters     []TemplateParamsVars `json:"parameters"`
	Variables      []TemplateParamsVars `json:"variables"`
	Resources      []json.RawMessage    `json:"resources"`
}
type TemplateParamsVars struct {
	Key   string
	Value string
}

type NSGResponse struct {
	Value []NSG `json:"value"`
}

type NSG struct {
	Etag       string        `json:"-"`
	ID         string        `json:"-"`
	Location   string        `json:"location"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Properties NSGProperties `json:"properties"`
}

type NSGProperties struct {
	DefaultSecurityRules []SecurityRule `json:"-"`
	SecurityRules        []SecurityRule `json:"securityRules"`
}

type SecurityRule struct {
	Etag       string                 `json:"-"`
	ID         string                 `json:"-"`
	Name       string                 `json:"name"`
	Properties SecurityRuleProperties `json:"properties"`
}

type SecurityRuleProperties struct {
	Access                   string `json:"access"`
	Description              string `json:"description"`
	DestinationAddressPrefix string `json:"destinationAddressPrefix"`
	DestinationPortRange     string `json:"destinationPortRange"`
	Direction                string `json:"direction"`
	Priority                 int64  `json:"priority"`
	Protocol                 string `json:"protocol"`
	ProvisioningState        string `json:"-"`
	SourceAddressPrefix      string `json:"sourceAddressPrefix"`
	SourcePortRange          string `json:"sourcePortRange"`
}

func (n NSG) Transform() TemplateNSG {
	formattedNSG := TemplateNSG{
		Location:   n.Location,
		Name:       n.Name,
		APIVersion: "2016-03-30",
		Type:       "Microsoft.Network/networkSecurityGroups",
		Properties: n.Properties,
	}

	return formattedNSG
}

type TemplateNSG struct {
	ResourceType string        `json:"-"`
	Comments     string        `json:"comments"`
	Location     string        `json:"location"`
	Name         string        `json:"name"`
	APIVersion   string        `json:"apiVersion"`
	Properties   NSGProperties `json:"properties"`
	Type         string        `json:"type"`
}
