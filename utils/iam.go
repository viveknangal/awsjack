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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

//"iamDetail" holds the data of individual User retrieved from AWS IAM API
type iamDetail struct {
	Count                   int
	Username                string
	PasswordLastUsed        string
	CreateDate              string
	MfaConig                string
	UserAttchedGroupsOutput string
	UserAttchedPolicyOutput string
	UserInlinePolicyOutput  string
}

//"iamOutput" holds  data  of all Users retrieved from AWS IAM API
type iamOutput struct {
	IamDetailList []iamDetail
	Mesg          string
	Regionlist    []string
}

// This function used for fetching IAM User details
func IamDetails() iamOutput {
	//sess, _ := session.NewSession(&aws.Config{})
	sess, _ := session.NewSession(&aws.Config{})
	fmt.Println("IAMMMMMMMMMM inGIT")
	iamsvc := iam.New(sess)
	iamOutput1, _ := iamsvc.ListUsers(&iam.ListUsersInput{})
	count := 0
	iDetailList := []iamDetail{}
	iDetail := iamDetail{}
	iOutput := iamOutput{}

	for _, value := range iamOutput1.Users {
		count++

		mfaOutput, _ := iamsvc.ListMFADevices(&iam.ListMFADevicesInput{UserName: value.UserName})

		mfaConig := "Enabled"
		if len(mfaOutput.MFADevices) == 0 {

			mfaConig = "Disabled"
		}

		PasswordLastUsed := ""
		if value.PasswordLastUsed == nil {
			PasswordLastUsed = "Password Disabled"
		} else {
			PasswordLastUsed = strings.Split((value.PasswordLastUsed).Format(time.RFC850), " ")[1]
		}
		InlinePolicyOutput, _ := iamsvc.ListUserPolicies(&iam.ListUserPoliciesInput{UserName: value.UserName})
		UserPolicies, _ := iamsvc.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{UserName: value.UserName})
		UserGroups, _ := iamsvc.ListGroupsForUser(&iam.ListGroupsForUserInput{UserName: value.UserName})

		userAttchedGroupsOutput := []string{}
		for _, b := range UserGroups.Groups {
			userAttchedGroupsOutput = append(userAttchedGroupsOutput, aws.StringValue(b.GroupName))
		}
		if len(userAttchedGroupsOutput) == 0 {
			userAttchedGroupsOutput = append(userAttchedGroupsOutput, "None")
		}
		userAttchedPolicyOutput := []string{}
		for _, j := range UserPolicies.AttachedPolicies {
			userAttchedPolicyOutput = append(userAttchedPolicyOutput, aws.StringValue(j.PolicyName))
		}
		if len(userAttchedPolicyOutput) == 0 {
			userAttchedPolicyOutput = append(userAttchedPolicyOutput, "None")
		}
		userInlinePolicyOutput := aws.StringValueSlice(InlinePolicyOutput.PolicyNames)
		if len(userInlinePolicyOutput) == 0 {
			userInlinePolicyOutput = append(userInlinePolicyOutput, "None")
		}

		fmt.Printf(">>>>> %v>>>>>%35v | %10v |%10v |%v| %v |%v | %v \n", count, aws.StringValue(value.UserName),
			PasswordLastUsed, strings.Split((value.CreateDate).Format(time.RFC850), " ")[1],
			mfaConig, userAttchedGroupsOutput, userAttchedPolicyOutput, userInlinePolicyOutput)

		iDetail.Count = count
		iDetail.Username = aws.StringValue(value.UserName)
		iDetail.PasswordLastUsed = PasswordLastUsed
		iDetail.CreateDate = strings.Split((value.CreateDate).Format(time.RFC850), " ")[1]
		iDetail.MfaConig = mfaConig
		iDetail.UserAttchedGroupsOutput = strings.Join(userAttchedGroupsOutput, ",")
		iDetail.UserAttchedPolicyOutput = strings.Join(userAttchedPolicyOutput, ",")
		iDetail.UserInlinePolicyOutput = strings.Join(userInlinePolicyOutput, ",")
		iDetailList = append(iDetailList, iDetail)
	}
	iOutput.IamDetailList = iDetailList
	return iOutput
}
