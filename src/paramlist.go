package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	goopt "github.com/droundy/goopt"
	"os"
	"sort"
)

var service = createSSMService()

func main() {
	initGoOpt()

	for _, param := range describeParameters(goopt.Args) {
		fmt.Println(param)
	}
}

func createSSMService() *ssm.SSM {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	if err != nil {
		exitErrorf("Unable to describe parameters, %v", err)
	}

	// Create SSM service client
	return ssm.New(sess)
}

func initGoOpt() {
	goopt.Description = func() string {
		return "Get parameter names from AWS SSM with optional prefixes."
	}
	goopt.Version = version
	goopt.Summary = "paramlist [prefixes]"
	goopt.Parse(nil)
}

func describeParameters(prefixes []string) []string {
	paramNames := []string{}
	if len(prefixes) <= 0 {
		paramNames = getAllParamNames()
	} else {
		for _, prefix := range prefixes {
			for _, name := range getParamNames(prefix) {
				paramNames = appendIfMissing(paramNames, name)
			}
		}
	}
	sort.Strings(paramNames)
	return paramNames
}

func getParamNames(prefix string) []string {
	filters := []*ssm.ParametersFilter{&ssm.ParametersFilter{
		Key:    aws.String("Name"),
		Values: []*string{aws.String(prefix)},
	}}
	paramNames := []string{}
	err := service.DescribeParametersPages(&ssm.DescribeParametersInput{
		Filters: filters},
		func(page *ssm.DescribeParametersOutput, lastPage bool) bool {
			for _, parameter := range page.Parameters {
				paramNames = append(paramNames, aws.StringValue(parameter.Name))
			}
			return true
		})

	if err != nil {
		exitErrorf("Unable to describe parameters, %v", err)
	}

	return paramNames
}

func getAllParamNames() []string {
	return getParamNames(" ")
}

func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var version string = "😂👌💯🔥🔥😂"