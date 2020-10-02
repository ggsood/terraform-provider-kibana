terraform {
  required_version = ">= 0.12.29"

  required_providers {
    kibana = {
      source = "ggsood/kibana"
      version = "1.0.3"
    }
  }
}

provider "kibana" {
  url     = ""
  username = ""
  password = ""
}

data "template_file" "sample_pipeline" {
   template = file("${path.cwd}/pipelines/sample.conf")
}

resource kibana_logstash_pipeline "test" {
  name 				= "terraform-test"
  description 		= "test"
  pipeline			= data.template_file.sample_pipeline.rendered
  settings = {
    "pipeline.batch.delay": 50,
    "pipeline.batch.size": 125,
    "pipeline.workers": 1,
    "queue.checkpoint.writes": 10,
    "queue.max_bytes": "2gb",
    "queue.type": "persisted",
  }
}