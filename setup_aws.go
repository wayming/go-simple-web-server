package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
)

func main(){
    svc := ec2.New(session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-2"),
    }))
    input := &ec2.CreateVpcInput{
        CidrBlock: aws.String("10.0.0.0/16"),
    }
    fmt.Println("VPC created")

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
