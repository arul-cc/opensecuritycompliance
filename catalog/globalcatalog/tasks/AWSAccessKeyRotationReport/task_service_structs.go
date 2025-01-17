// This file is autogenerated. Modify as per your task needs.
package main

type UserInputs struct {
	MaxAccessKeyAge     int    `yaml:"MaxAccessKeyAge" validate:"required"`
	AWSCredentialReport string `yaml:"AWSCredentialReport" validate:"required"`
}

type Outputs struct {
	AccessKeyRotationReport string
	LogFile                 string
	MetaFile                string
}
