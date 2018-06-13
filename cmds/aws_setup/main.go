package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const cmdCreateVPC = "CreateVPC"
const cmdDeleteVPC = "DeleteVPC"
const cmdListVPC = "ListVPC"

func createVPC() {
	var session, err = session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")})
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	}
	fmt.Println("VPC created")

	svc := ec2.New(session)
	result, err := svc.CreateVpc(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func deleteVPC() {
	svc := ec2.New(session.New())
	fmt.Println(svc)

}

func listVPC() {
	svc := ec2.New(session.New())
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("isDefault"),
				Values: []*string{aws.String("true")},
			},
		},
	}

	result, err := svc.DescribeVpcs(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
func main() {

	var command = os.Args[1]
	fmt.Println(command)
	switch command {
	case cmdCreateVPC:
		fmt.Print("Not supported yet")
	case cmdDeleteVPC:
		fmt.Print("Not supported yet")
	case cmdListVPC:
		listVPC()
	default:
		fmt.Println("Unknown command")
	}

}
