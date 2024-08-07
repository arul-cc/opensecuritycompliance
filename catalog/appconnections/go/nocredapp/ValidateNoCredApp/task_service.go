// This file is autogenerated. Modify as per your task needs.

package main

// ValidateNoCredApp :
func (inst *TaskInstance) ValidateNoCredApp(inputs *UserInputs, outputs *Outputs) (err error) {
	outputs.IsValidated, err = inputs.Validate()
	if err != nil {
		outputs.ValidationMessage = err.Error()
	} else {
		outputs.ValidationMessage = "Credentials validated successfully"
	}

	return nil
}