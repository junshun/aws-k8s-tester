// gen generates eksconfig documentation.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	k8s_tester "github.com/aws/aws-k8s-tester/k8s-tester"
	cloudwatch_agent "github.com/aws/aws-k8s-tester/k8s-tester/cloudwatch-agent"
	"github.com/aws/aws-k8s-tester/k8s-tester/clusterloader"
	"github.com/aws/aws-k8s-tester/k8s-tester/configmaps"
	"github.com/aws/aws-k8s-tester/k8s-tester/conformance"
	csi_ebs "github.com/aws/aws-k8s-tester/k8s-tester/csi-ebs"
	"github.com/aws/aws-k8s-tester/k8s-tester/csrs"
	falco "github.com/aws/aws-k8s-tester/k8s-tester/falco"
	fluent_bit "github.com/aws/aws-k8s-tester/k8s-tester/fluent-bit"
	jobs_echo "github.com/aws/aws-k8s-tester/k8s-tester/jobs-echo"
	jobs_pi "github.com/aws/aws-k8s-tester/k8s-tester/jobs-pi"
	kubernetes_dashboard "github.com/aws/aws-k8s-tester/k8s-tester/kubernetes-dashboard"
	metrics_server "github.com/aws/aws-k8s-tester/k8s-tester/metrics-server"
	nlb_hello_world "github.com/aws/aws-k8s-tester/k8s-tester/nlb-hello-world"
	php_apache "github.com/aws/aws-k8s-tester/k8s-tester/php-apache"
	"github.com/aws/aws-k8s-tester/k8s-tester/secrets"
	"github.com/aws/aws-k8s-tester/k8s-tester/stress"
	stress_in_cluster "github.com/aws/aws-k8s-tester/k8s-tester/stress/in-cluster"
	aws_v1_ecr "github.com/aws/aws-k8s-tester/utils/aws/v1/ecr"
	"github.com/olekukonko/tablewriter"
)

func main() {
	doc := createDoc()
	if err := ioutil.WriteFile("../../README.config.md", []byte("\n```\n"+doc+"```\n"), 0666); err != nil {
		panic(err)
	}
	fmt.Println("generated")
}

func createDoc() string {
	es := &enableEnvVars{}
	b := strings.Builder{}

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX, &k8s_tester.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+cloudwatch_agent.Env()+"_", &cloudwatch_agent.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+fluent_bit.Env()+"_", &fluent_bit.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+metrics_server.Env()+"_", &metrics_server.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+conformance.Env()+"_", &conformance.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+csi_ebs.Env()+"_", &csi_ebs.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+kubernetes_dashboard.Env()+"_", &kubernetes_dashboard.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+falco.Env()+"_", &falco.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+php_apache.Env()+"_", &php_apache.Config{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+php_apache.EnvRepository()+"_", &aws_v1_ecr.Repository{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+nlb_hello_world.Env()+"_", &nlb_hello_world.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+jobs_pi.Env()+"_", &jobs_pi.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+jobs_echo.Env("Job")+"_", &jobs_echo.Config{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+jobs_echo.EnvRepository("Job")+"_", &aws_v1_ecr.Repository{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+jobs_echo.Env("CronJob")+"_", &jobs_echo.Config{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+jobs_echo.EnvRepository("CronJob")+"_", &aws_v1_ecr.Repository{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+csrs.Env()+"_", &csrs.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+configmaps.Env()+"_", &configmaps.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+secrets.Env()+"_", &secrets.Config{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+clusterloader.Env()+"_", &clusterloader.Config{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+clusterloader.EnvTestOverride()+"_", &clusterloader.TestOverride{}))

	b.WriteByte('\n')
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+stress.Env()+"_", &stress.Config{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+stress.EnvRepository()+"_", &aws_v1_ecr.Repository{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+stress_in_cluster.Env()+"_", &stress_in_cluster.Config{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+stress_in_cluster.EnvK8sTesterStressRepository()+"_", &aws_v1_ecr.Repository{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+stress_in_cluster.EnvK8sTesterStressCLI()+"_", &stress_in_cluster.K8sTesterStressCLI{}))
	b.WriteByte('\n')
	b.WriteString(es.writeDoc(k8s_tester.ENV_PREFIX+stress_in_cluster.EnvK8sTesterStressCLIBusyboxRepository()+"_", &aws_v1_ecr.Repository{}))

	b.WriteByte('\n')
	b.WriteByte('\n')

	txt := b.String()

	return txt
}

type enableEnvVars struct {
}

var columns = []string{
	"environmental variable",
	"field type",
	"type",
	"go type",
}

func (es *enableEnvVars) writeDoc(pfx string, st interface{}) string {
	buf := bytes.NewBuffer(nil)
	tb := tablewriter.NewWriter(buf)
	tb.SetAutoWrapText(false)
	tb.SetAlignment(tablewriter.ALIGN_LEFT)
	tb.SetColWidth(1500)
	tb.SetCenterSeparator("*")
	tb.SetHeader(columns)

	ts := reflect.TypeOf(st)
	tp, vv := reflect.TypeOf(st).Elem(), reflect.ValueOf(st).Elem()
	for i := 0; i < tp.NumField(); i++ {
		jv := tp.Field(i).Tag.Get("json")
		if jv == "" || jv == "-" {
			continue
		}
		if vv.Field(i).Type().Kind() == reflect.Ptr {
			continue
		}

		ft := "settable via env var"
		if tp.Field(i).Tag.Get("read-only") == "true" {
			ft = "read-only"
		}

		jv = strings.Replace(jv, ",omitempty", "", -1)
		jv = strings.ToUpper(strings.Replace(jv, "-", "_", -1))
		env := pfx + jv

		tb.Append([]string{
			env,
			strings.ToUpper(ft),
			fmt.Sprintf("%s.%s", ts, tp.Field(i).Name),
			fmt.Sprintf("%s", vv.Field(i).Type()),
		})
	}

	tb.Render()
	return buf.String()
}
