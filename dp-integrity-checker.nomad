job "dp-integrity-checker" {
  datacenters = ["eu-west-2"]
  region      = "eu"
  type        = "batch"

  periodic {
    cron             = "0 0 1 * * * *"
    time_zone        = "UTC"
    prohibit_overlap = true
  }

  group "publishing" {
    count = "{{PUBLISHING_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "publishing-mount"
    }

    restart {
      attempts = 3
      delay    = "15s"
      interval = "1m"
      mode     = "delay"
    }

    task "dp-integrity-checker" {
      driver = "docker"

      artifact {
        source = "s3::https://s3-eu-west-2.amazonaws.com/{{DEPLOYMENT_BUCKET}}/dp-integrity-checker/{{PROFILE}}/{{RELEASE}}.tar.gz"
      }

      config {
        command = "${NOMAD_TASK_DIR}/start-task"

        args = ["./dp-integrity-checker"]

        image = "{{ECR_URL}}:concourse-{{REVISION}}"

        mounts = [
          {
            type     = "bind"
            target   = "/content"
            source   = "/var/florence"
            readonly = true
          }
        ]
      }

      resources {
        cpu    = "{{PUBLISHING_RESOURCE_CPU}}"
        memory = "{{PUBLISHING_RESOURCE_MEM}}"

        network {
          port "http" {}
        }
      }

      template {
        source      = "${NOMAD_TASK_DIR}/vars-template"
        destination = "${NOMAD_TASK_DIR}/vars"
      }

      vault {
        policies = ["dp-integrity-checker"]
      }
    }
  }
}
