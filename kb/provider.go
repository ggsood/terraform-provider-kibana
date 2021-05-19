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

type ProviderConf struct {
	rawUrl          string
	insecure        bool
	caCertFiles     []string
	username        string
	password        string
	parsedUrl       *url.URL
	maxRetry        int
	waitBeforeRetry int
	debug           bool
}

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
	rawUrl := d.Get("url").(string)
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	return &ProviderConf{
		rawUrl:          rawUrl,
		insecure:        d.Get("insecure").(bool),
		caCertFiles:     convertArrayInterfaceToArrayString(d.Get("cacert_files").(*schema.Set).List()),
		username:        d.Get("username").(string),
		password:        d.Get("password").(string),
		parsedUrl:       parsedUrl,
		maxRetry:        d.Get("retry").(int),
		waitBeforeRetry: d.Get("wait_before_retry").(int),
		debug:           d.Get("debug").(bool),
	}, nil
}

func getClient(conf *ProviderConf) (*kibana.Client, error) {
	if conf.debug {
		log.SetLevel(log.DebugLevel)
	}

	// Intialise connection
	cfg := kibana.Config{
		Address:          conf.rawUrl,
		CAs:              conf.caCertFiles,
		DisableVerifySSL: conf.insecure,
	}
	if conf.username != "" && conf.password != "" {
		cfg.Username = conf.username
		cfg.Password = conf.password
	}

	client, err := kibana.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Test connection and check kibana version
	nbFailed := 0
	isOnline := false
	var kibanaStatus kbapi.KibanaStatus
	for isOnline == false {
		kibanaStatus, err = client.API.KibanaStatus.Get()
		if err == nil {
			isOnline = true
		} else {
			if nbFailed == conf.maxRetry {
				return nil, err
			}
			nbFailed++
			time.Sleep(time.Duration(conf.waitBeforeRetry) * time.Second)
		}
	}

	if kibanaStatus == nil {
		return nil, errors.New("Status is empty, something wrong with Kibana?")
	}

	version := kibanaStatus["version"].(map[string]interface{})["number"].(string)
	log.Debugf("Server: %s", version)

	if version < "7.0.0" {
		return nil, errors.New("Kibana is older than 7.0.0")
	} else if version >= "8.0.0" {
		return nil, errors.New("Kibana is version 8.0.0 or newer and support has not been tested")
	}

	log.Printf("[INFO] Using Kibana 7")
	return client, nil
}
