// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package kb

import (
	"net/url"
	"time"

	kibana "github.com/ggsood/go-kibana-rest/v7"
	"github.com/ggsood/go-kibana-rest/v7/kbapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Provider define kibana provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_URL", nil),
				Description: "Kibana URL",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_USERNAME", nil),
				Description: "Username to use to connect to Kibana using basic auth",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_PASSWORD", nil),
				Description: "Password to use to connect to Kibana using basic auth",
			},
			"cacert_files": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A Custom CA certificates path",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable SSL verification of API calls",
			},
			"retry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     6,
				Description: "Nummber time it retry connexion before failed",
			},
			"wait_before_retry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Wait time in second before retry connexion",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable debug log level in provider",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kibana_user_space":        resourceKibanaUserSpace(),
			"kibana_role":              resourceKibanaRole(),
			"kibana_object":            resourceKibanaObject(),
			"kibana_logstash_pipeline": resourceKibanaLogstashPipeline(),
			"kibana_copy_object":       resourceKibanaCopyObject(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var (
		relevantClient interface{}
	)

	URL := d.Get("url").(string)
	insecure := d.Get("insecure").(bool)
	cacertFiles := convertArrayInterfaceToArrayString(d.Get("cacert_files").(*schema.Set).List())
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	retry := d.Get("retry").(int)
	waitBeforeRetry := d.Get("wait_before_retry").(int)
	debug := d.Get("debug").(bool)

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// Checks is valid URL
	if _, err := url.Parse(URL); err != nil {
		return nil, err
	}

	// Intialise connexion
	cfg := kibana.Config{
		Address: URL,
		CAs:     cacertFiles,
	}
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}
	if insecure {
		cfg.DisableVerifySSL = true
	}

	client, err := kibana.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Test connexion and check kibana version
	nbFailed := 0
	isOnline := false
	var kibanaStatus kbapi.KibanaStatus
	for isOnline {
		kibanaStatus, err = client.API.KibanaStatus.Get()
		if err == nil {
			isOnline = true
		} else {
			if nbFailed == retry {
				return nil, err
			}
			nbFailed++
			time.Sleep(time.Duration(waitBeforeRetry) * time.Second)
		}
	}

	if kibanaStatus == nil {
		return nil, errors.New("Status is empty, somethink wrong with Kibana ?")
	}

	version := kibanaStatus["version"].(map[string]interface{})["number"].(string)
	log.Debugf("Server: %s", version)

	if version < "8.0.0" && version >= "7.0.0" {
		log.Printf("[INFO] Using Kibana 7")
		relevantClient = client
	} else {
		return nil, errors.New("Kibana is older than 7.0.0")
	}

	return relevantClient, nil
}
