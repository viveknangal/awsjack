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

package awsjack

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// This function retrieves "Security Groups"  related details for EC2 instances
func sgInfo(ecsvc *ec2.EC2, y *ec2.Instance) []string {
	sgPortList := []string{}

	sgList := []string{}
	for _, y := range y.SecurityGroups {
		sgList = append(sgList, *y.GroupId)
	}
	sgRef := ec2.DescribeSecurityGroupsInput{GroupIds: aws.StringSlice(sgList)}
	sgOutput, _ := ecsvc.DescribeSecurityGroups(&sgRef)

	for _, b := range sgOutput.SecurityGroups {
		for _, c := range b.IpPermissions {
			for _, d := range c.IpRanges {
				if *d.CidrIp == "0.0.0.0/0" {

					if c.FromPort == nil {
						sgPortList = append(sgPortList, "0-65535")
					} else if *c.FromPort == *c.ToPort {
						sgPortList = append(sgPortList, strconv.FormatInt(*c.FromPort, 10))
					} else {

						output := fmt.Sprintf("%v-%v", *c.FromPort, *c.ToPort)
						sgPortList = append(sgPortList, output)
					}
				}

			}

		}

	}

	return sgPortList
}
