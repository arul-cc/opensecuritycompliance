// This file is autogenerated. Modify as per your task needs.
package main

type UserInputs struct {
	IncludeCriteria      string `yaml:"IncludeCriteria"`
	ExcludeCriteria      string `yaml:"ExcludeCriteria"`
	OpaConfigurationFile string `yaml:"OpaConfigurationFile"`
}

type Outputs struct {
	Source   string `yaml:"Source"`
	DataFile string
	LogFile  string
}
