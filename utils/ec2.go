// MIT License

// Copyright (c) 2020 Vivek Aggarwal

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This "awsjack" fetch details realted to Instances & IAM
package awsjack

import (
	"fmt"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go/service/pricing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//"instanceDetail" holds the data  of individual instance retrieved from AWS EC2 API
type instanceDetail struct {
	Count int
	//RestartDate  string
	Price        float64
	LaunchDate   string
	DiskCount    int
	Sgroup       []string
	Name         string
	PrivateIP    string
	InstanceType string
	Az           string
	Vpc          string
	PublicIP     string
	Key          string
	ImageID      string
	InstanceID   string
	State        string
}

//"instanceInfoOutput" holds the data  of all EC2 instances retrieved from AWS EC2 API
type instanceInfoOutput struct {
	InstanceDetailList []instanceDetail
	Mesg               string
	Regionlist         []string
}

// This function call "instanceInfo" which retrieves instance details
func Ec2Details(region string) instanceInfoOutput {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	ec2svc := ec2.New(sess)

	dd := instanceInfo(ec2svc, region)

	return dd
}

// This function used for fetching Details Pricing inventory for a given  Instance  type
func getProductInput(region string, instanceType string, usagetype string) *pricing.GetProductsOutput {
	regionMap := map[string]string{
		"us-east-2":      "US East (Ohio)",
		"us-east-1":      "US East (N. Virginia)",
		"us-west-1":      "US West (N. California)",
		"us-west-2":      "US West (Oregon)",
		"af-south-1":     "Africa (Cape Town)",
		"ap-east-1":      "Asia Pacific (Hong Kong)",
		"ap-south-1":     "Asia Pacific (Mumbai)",
		"ap-northeast-3": "Asia Pacific (Osaka-Local)",
		"ap-northeast-2": "Asia Pacific (Seoul)",
		"ap-southeast-1": "Asia Pacific (Singapore)",
		"ap-southeast-2": "Asia Pacific (Sydney)",
		"ap-northeast-1": "Asia Pacific (Tokyo)",
		"ca-central-1":   "Canada (Central)",
		"cn-north-1":     "China (Beijing)",
		"cn-northwest-1": "China (Ningxia)",
		"eu-central-1":   "EU (Frankfurt)",
		"eu-west-1":      "EU (Ireland)",
		"eu-west-2":      "EU (London)",
		"eu-south-1":     "EU (Milan)",
		"eu-west-3":      "EU (Paris)",
		"eu-north-1":     "EU (Stockholm)",
		"me-south-1":     "Middle East (Bahrain)",
		"sa-east-1":      "South America (Sao Paulo)",
		"us-gov-east-1":  "AWS GovCloud (US-East)",
		"us-gov-west-1":  "AWS GovCloud (US)"}

	regionName := regionMap[region]

	input := pricing.GetProductsInput{ServiceCode: aws.String("AmazonEC2"),
		Filters: []*pricing.Filter{
			{
				Field: aws.String("instanceType"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(instanceType),
			},
			{
				Field: aws.String("operatingSystem"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Linux"),
			},
			{
				Field: aws.String("location"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(regionName),
			},
			{
				Field: aws.String("operation"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("RunInstances"),
			},
			{
				Field: aws.String("usagetype"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String(usagetype),
			},
		},
	}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	priceSvc := pricing.New(sess)
	priceoutput, _ := priceSvc.GetProducts(&input)
	return priceoutput
}

// This function used for fetching Pricing detail for a EC2 Instance  type
func getprice(region string, instanceType string) float64 {

	var pricePerDay float64

	usageTypeMap := map[string]string{"us-east-2": "USE2", "us-east-1": "USE1", "us-west-1": "USW1", "us-west-2": "USW2", "af-south-1": "AFS1", "ap-east-1": "APE1", "ap-south-1": "APS3", "ap-northeast-3": "APN3", "ap-northeast-2": "APN2", "ap-southeast-1": "APS1", "ap-southeast-2": "APS2", "ap-northeast-1": "APN1", "ca-central-1": "CAN1", "cn-north-1": "CNN1", "cn-northwest-1": "CNW1", "eu-central-1": "EUC1", "eu-west-1": "EUW1", "eu-west-2": "EUW2", "eu-south-1": "EUS1", "eu-west-3": "EUS3", "eu-north-1": "EUN1", "me-south-1": "MES1", "sa-east-1": "SAE1", "us-gov-east-1": "AWS GovCloud (US-East)", "us-gov-west-1": "AWS GovCloud (US)"}

	usageTypeVal := usageTypeMap[region] + "-BoxUsage:" + instanceType
	priceoutput := getProductInput(region, instanceType, usageTypeVal)

	pricelist := priceoutput.PriceList

	if len(pricelist) == 0 && usageTypeVal != ("BoxUsage:"+instanceType) {
		usageTypeVal := "BoxUsage:" + instanceType
		priceoutput := getProductInput(region, instanceType, usageTypeVal)
		pricelist = priceoutput.PriceList
	}
	if len(pricelist) > 0 {
		dd := pricelist[0]
		for _, value1 := range ((dd["terms"].(map[string]interface{}))["OnDemand"]).(map[string]interface{}) {
			interim1 := value1.(map[string]interface{})["priceDimensions"]
			for _, value2 := range interim1.(map[string]interface{}) {
				cost := (value2.(map[string]interface{})["pricePerUnit"]).(map[string]interface{})["USD"]

				priceFloat, _ := strconv.ParseFloat(fmt.Sprint(cost), 64)
				pricePerHour := (math.Round((priceFloat)*100) / 100)

				pricePerDay = (math.Round((pricePerHour*24)*100) / 100)

			}

		}
	}
	return pricePerDay
}

// This function retrieves instance details for all instances in a given region
func instanceInfo(ec2svc *ec2.EC2, region string) instanceInfoOutput {
	instancesOutput, _ := ec2svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	count := 0
	instancePriceSum := 0.0
	insDetailList := []instanceDetail{}
	insDetail := instanceDetail{}

	insInfoOutput := instanceInfoOutput{}
	priceList := map[string]float64{}
	for _, j := range instancesOutput.Reservations {

		for _, y := range j.Instances {
			count = count + 1
			pip := ""
			key := ""
			privip := ""
			if y.PrivateIpAddress == nil {
				privip = "NA"
			} else {
				privip = *y.PrivateIpAddress
			}

			if y.PublicIpAddress == nil {
				pip = "NA"
			} else {
				pip = *y.PublicIpAddress
			}

			if y.KeyName == nil {
				fmt.Println("lenght=")
			} else {
				key = *y.KeyName
			}
			nameTag := "NA"
			if len(y.Tags) > 0 {
				for _, g := range y.Tags {

					if aws.StringValue(g.Key) == "Name" {
						nameTag = aws.StringValue(g.Value)
					}
				}
			}
			sgOutput := sgInfo(ec2svc, y)
			diskOutput := diskInfo(y)
			instanceType := aws.StringValue(y.InstanceType)
			val, ok := priceList[instanceType]
			var instancePrice float64
			if ok {
				instancePrice = val
			} else {

				instancePrice = getprice(region, instanceType)
				priceList[instanceType] = instancePrice

			}

			insDetail.Count = count
			insDetail.Price = instancePrice
			insDetail.LaunchDate = diskOutput.date
			insDetail.DiskCount = diskOutput.diskCount
			insDetail.Sgroup = sgOutput
			insDetail.Name = nameTag
			insDetail.PrivateIP = privip
			insDetail.InstanceType = aws.StringValue(y.InstanceType)
			insDetail.Az = aws.StringValue(y.Placement.AvailabilityZone)
			insDetail.State = aws.StringValue(y.State.Name)
			insDetail.Vpc = aws.StringValue(y.VpcId)
			insDetail.PublicIP = pip
			insDetail.Key = key
			insDetail.ImageID = aws.StringValue(y.ImageId)
			insDetail.InstanceID = aws.StringValue(y.InstanceId)
			insDetailList = append(insDetailList, insDetail)
			instancePriceSum = instancePriceSum + instancePrice
			fmt.Println(count, ")",
				instancePrice, diskOutput.date, diskOutput.diskCount, sgOutput, nameTag,
				privip, *y.InstanceType, *y.Placement.AvailabilityZone, *y.State.Name,
				*y.VpcId, pip, key, *y.ImageId, *y.InstanceId)
		}

	}
	insInfoOutput.InstanceDetailList = insDetailList
	if len(insInfoOutput.InstanceDetailList) < 1 {
		insInfoOutput.Mesg = "No Instances found for region=\"" + region + "\""

	} else {
		insInfoOutput.Mesg = "Instances list for Region=\"" + region + "\" & their per/day Instance running cost=\"" + fmt.Sprintf("%.2f", instancePriceSum) + "\""
	}
	return insInfoOutput
}
