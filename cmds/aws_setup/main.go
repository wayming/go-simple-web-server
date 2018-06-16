package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const cmdCreateVPC = "CreateVPC"
const cmdDeleteVPC = "DeleteVPC"
const cmdListVPC = "ListVPC"
const cmdCreateDefaultVPC = "CreateDefaultVPC"

func printAWSError(e error) {
	if e != nil {
		if aerr, ok := e.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(e.Error())
		}
	}
}

func getVPC(awsSession *session.Session) *ec2.DescribeVpcsOutput {
	var filter ec2.Filter
	input := &ec2.DescribeVpcsInput{Filters: []*ec2.Filter{&filter}}

	result, err := ec2.New(awsSession).DescribeVpcs(input)
	if err != nil {
		printAWSError(err)
		return nil
	}

	return result
}

func getSubnets(awsSession *session.Session, vpc string) *ec2.DescribeSubnetsOutput {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpc)},
			},
		},
	}

	result, err := ec2.New(awsSession).DescribeSubnets(input)
	if err != nil {
		printAWSError(err)
		return nil
	}

	return result
}

func getInternetGateway(awsSession *session.Session, vpc string) *ec2.DescribeInternetGatewaysOutput {
	input := &ec2.DescribeInternetGatewaysInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("attachment.vpc-id"),
				Values: []*string{
					aws.String(vpc),
				},
			},
		},
	}

	result, err := ec2.New(awsSession).DescribeInternetGateways(input)
	if err != nil {
		printAWSError(err)
		return nil
	}

	return result
}

func createDefaultVPC(awsSession *session.Session) {
	input := &ec2.CreateDefaultVpcInput{}
	result, err := ec2.New(awsSession).CreateDefaultVpc(input)
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println("VPC created")
	fmt.Println(result)
}

func createVPC(awsSession *session.Session) {
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	}
	result, err := ec2.New(awsSession).CreateVpc(input)
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println("VPC created")
	fmt.Println(result)
}

func deleteVPC(awsSession *session.Session, vpcName string) {
	input := &ec2.DeleteVpcInput{
		VpcId: aws.String(vpcName),
	}

	internetGatewaysOutput := getInternetGateway(awsSession, vpcName)
	for _, internetGateway := range internetGatewaysOutput.InternetGateways {
		fmt.Println("Delete dependent internet gateway " + *internetGateway.InternetGatewayId + " for VPC " + vpcName)
		detachInternetGateway(awsSession, vpcName, *internetGateway.InternetGatewayId)
		deleteInternetGateway(awsSession, *internetGateway.InternetGatewayId)
	}

	subnetsOutput := getSubnets(awsSession, vpcName)
	for _, subnet := range subnetsOutput.Subnets {
		fmt.Println("Delete dependent subnet " + *subnet.SubnetId + " for VPC " + vpcName)
		deleteSubnet(awsSession, *subnet.SubnetId)
	}

	result, err := ec2.New(awsSession).DeleteVpc(input)
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println(result)
}

func detachInternetGateway(awsSession *session.Session, vpcName string, intenetGatewayName string) {
	input := &ec2.DetachInternetGatewayInput{
		VpcId:             aws.String(vpcName),
		InternetGatewayId: aws.String(intenetGatewayName),
	}

	result, err := ec2.New(awsSession).DetachInternetGateway(input)
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println(result)
}

func deleteInternetGateway(awsSession *session.Session, intenetGatewayName string) {
	input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String(intenetGatewayName),
	}

	result, err := ec2.New(awsSession).DeleteInternetGateway(input)
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println(result)
}

func deleteSubnet(awsSession *session.Session, subnetName string) {
	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetName),
	}

	result, err := ec2.New(awsSession).DeleteSubnet(input)
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println(result)
}

func listVPC(awsSession *session.Session) {
	vpcs := getVPC(awsSession)
	fmt.Println(vpcs)

	for _, vpc := range vpcs.Vpcs {
		listSubnets(awsSession, *vpc.VpcId)
		listIntenetGateway(awsSession, *vpc.VpcId)
	}

}

func listSubnets(awsSession *session.Session, vpc string) {
	subnets := getSubnets(awsSession, vpc)
	fmt.Println(subnets)
}

func listIntenetGateway(awsSession *session.Session, vpc string) {
	internetGateway := getInternetGateway(awsSession, vpc)
	fmt.Println(internetGateway)
}

func main() {
	command := flag.String(
		"exec", cmdListVPC,
		"ListVPC, CreateDefaultVPC, CreateVPC or DeleteVPC")
	vpcName := flag.String("vpc", "", "name of the VPC")
	flag.Parse()

	if *command == cmdDeleteVPC && *vpcName == "" {
		fmt.Println("Please specify the name of vpc")
		return
	}

	var session, err = session.NewSession(&aws.Config{Region: aws.String("ap-southeast-2")})
	if err != nil {
		printAWSError(err)
		return
	}

	fmt.Println(command)
	switch *command {
	case cmdCreateVPC:
		createVPC(session)
	case cmdCreateDefaultVPC:
		createDefaultVPC(session)
	case cmdDeleteVPC:
		deleteVPC(session, *vpcName)
	case cmdListVPC:
		listVPC(session)
	default:
		fmt.Println("Unknown command")
	}

}
