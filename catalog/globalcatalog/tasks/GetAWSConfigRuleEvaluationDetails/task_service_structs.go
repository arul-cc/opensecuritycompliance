// This file is autogenerated. Modify as per your task needs.
package main

type UserInputs struct {
	Region []string `yaml:"Region" validate:"required"`
}

type Outputs struct {
	ComplianceStatus_                 string `description:"ComplianceStatus"`
	CompliancePCT_                    int    `description:"CompliancePCT"`
	AWSConfigRulesJSON                string `description:"AWS Config Rules Report JSON"`
	AWSConfigRuleEvaluationStatusJSON string `description:"AWS Config Rule Evaluation Status"`
	LogFile                           string `description:"Log file"`
	MetaDataFile                      string `description:"Meta file"`
}
