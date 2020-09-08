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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//"diskInfoType" holds the data retrieved from AWS EBS API
type diskInfoType struct {
	date      string
	diskCount int
}

// This function used for fetching EBS volume details
func diskInfo(y *ec2.Instance) diskInfoType {

	deviceMap := map[string]int64{}

	rootdevice := *y.RootDeviceName
	for _, f := range y.BlockDeviceMappings {
		deviceMap[aws.StringValue(f.DeviceName)] = f.Ebs.AttachTime.Unix()
	}
	dd := time.Unix(int64(deviceMap[rootdevice]), 0)
	splitdate := strings.Split(dd.Format(time.RFC850), " ")[1]
	diskAttached := len(deviceMap)
	dItype := diskInfoType{date: splitdate, diskCount: diskAttached}
	return dItype

}
