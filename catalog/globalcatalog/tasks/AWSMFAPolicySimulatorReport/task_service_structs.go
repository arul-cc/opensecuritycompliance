// This file is autogenerated. Modify as per your task needs.
package main

type UserInputs struct {
	Users                       string `yaml:"Users" validate:"required"`
	UserStatus                  string `yaml:"UserStatus" validate:"required,oneof=include exclude"`
	Roles                       string `yaml:"Roles" validate:"required"`
	RoleStatus                  string `yaml:"RoleStatus" validate:"required,oneof=include exclude"`
	Groups                      string `yaml:"Groups" validate:"required"`
	GroupStatus                 string `yaml:"GroupStatus" validate:"required,oneof=include exclude"`
	MFARecommendationFile       string `yaml:"MFARecommendationFile" validate:"required"`
	AccountAuthorizationDetails string `yaml:"AccountAuthorizationDetails" validate:"required"`
}

type Outputs struct {
	MFAPolicySimulatorReport string
	LogFile                  string
}