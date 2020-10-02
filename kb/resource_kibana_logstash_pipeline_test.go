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
	"fmt"
	"testing"

	kibana "github.com/ggsood/go-kibana-rest/v7"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
)

func TestAccKibanaLogstashPipeline(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaLogstashPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKibanaLogstashPipeline,
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaLogstashPipelineExists("kibana_logstash_pipeline.test"),
				),
			},
			{
				ResourceName:            "kibana_logstash_pipeline.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testCheckKibanaLogstashPipelineExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No logstash pipeline ID is set")
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		logstashPipeline, err := client.API.KibanaLogstashPipeline.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if logstashPipeline == nil {
			return errors.Errorf("Logstash pipeline %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaLogstashPipelineDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_logstash_pipeline" {
			continue
		}

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		logstashPipeline, err := client.API.KibanaLogstashPipeline.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		if logstashPipeline == nil {
			return nil
		}

		return fmt.Errorf("Logstash pipeline %q still exists", rs.Primary.ID)
	}

	return nil
}

var testKibanaLogstashPipeline = `
resource "kibana_logstash_pipeline" "test" {
  name 				= "terraform-test"
  description 		= "test"
  pipeline			= "input { stdin {} } output { stdout {} }"
  settings = {
	  "queue.type" = "persisted"
  }
}
`
