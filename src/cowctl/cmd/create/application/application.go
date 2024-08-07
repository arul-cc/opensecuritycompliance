package application

import (
	"cowlibrary/applications"
	"cowlibrary/constants"
	cowlibutils "cowlibrary/utils"
	"cowlibrary/vo"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"cowctl/utils"
	"cowctl/utils/terminalutils"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.NoArgs,

		Use:   "application",
		Short: "Create a application",
		Long:  "Create application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd)
		},
	}

	cmd.Flags().StringP("file-name", "f", "", "Set your file name")
	cmd.Flags().String("config-path", "", "path for the configuration file.")
	return cmd
}

func runE(cmd *cobra.Command) error {

	additionalInfo, err := utils.GetAdditionalInfoFromCmd(cmd)
	if err != nil {
		return err
	}

	yamlFileName, yamlFilePath := ``, ``

	if cmd.Flags().HasFlags() {
		if flagName := cmd.Flags().Lookup("yaml-file"); flagName != nil {
			yamlFileName = flagName.Value.String()
		}
	}

	if cowlibutils.IsNotEmpty(yamlFileName) {
		yamlFilePath = cowlibutils.GetYamlFilesPathFromApplicationCatalog(additionalInfo, yamlFileName)
	}

	if cowlibutils.IsEmpty(yamlFileName) || cowlibutils.IsFileNotExist(yamlFilePath) {
		fileNameFromCmd, err := utils.GetValueAsFileNameFromCmdPrompt("Select a valid yaml file name", additionalInfo.PolicyCowConfig.PathConfiguration.ApplicationClassPath, []string{".yaml", ".yml"})
		if err != nil || cowlibutils.IsEmpty(fileNameFromCmd) {
			return errors.New("cannot get the application")
		}
		yamlFilePath = filepath.Join(additionalInfo.PolicyCowConfig.PathConfiguration.ApplicationClassPath, fileNameFromCmd)
		if cowlibutils.IsFileNotExist(yamlFilePath) {
			return fmt.Errorf("%s is not a valid application file path", yamlFilePath)
		}
	}

	applicationVO := &vo.UserDefinedApplicationVO{}

	yamlFileByts, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFileByts, &applicationVO)
	if err != nil {
		return err
	}

	additionalInfo.ApplictionScopeConfigVO = &vo.ApplictionScopeConfigVO{FileData: yamlFileByts}

	if len(applicationVO.Spec.LinkableApplicationClasses) > 0 {
		errorDetailVO := applications.ValidateLinkedApplications(applicationVO.Spec, additionalInfo, yamlFilePath)
		if errorDetailVO != nil {
			return fmt.Errorf(errorDetailVO.Issue)
		}
	}
	if applications.IsAppAlreadyPresent(applicationVO.Meta, additionalInfo) {
		isConfirmed, err := terminalutils.GetConfirmationFromCmdPrompt("application already presented in the directory. Are you sure you going to re-initialize ?")
		if !isConfirmed || err != nil {
			return err
		}
		applicationVO.IsVersionToBeOverride = true
	}

	supportedLanguage := constants.SupportedLanguageGo

	filePath := filepath.Join(additionalInfo.PolicyCowConfig.PathConfiguration.AppConnectionPath, supportedLanguage.String())

	if supportedLanguage == constants.SupportedLanguagePython {
		filePath = filepath.Join(filePath, "appconnections")
	}

	packagePath := filepath.Join(filePath, strings.ToLower(applicationVO.Meta.Name))

	if cowlibutils.IsFolderExist(packagePath) {
		isConfirmed, err := terminalutils.GetConfirmationFromCmdPrompt("An application class implementation already exists for this name (the version will be excluded). Are you sure you want to re-create it?")
		if !isConfirmed || err != nil {
			return err
		}
		applicationVO.IsVersionToBeOverride = true
	}

	application := applications.GetApplication(supportedLanguage)

	applicationVOClone := *applicationVO

	errorDetails := application.Create(applicationVO, additionalInfo)
	if len(errorDetails) > 0 {
		utils.DrawErrorTable(errorDetails)
		return errors.New(constants.ErroInvalidData)
	}

	applicationVOClone.IsVersionToBeOverride = true
	if supportedLanguage == constants.SupportedLanguageGo {

		errorDetails = (&applications.PythonApplicationHandler{Context: cmd.Context()}).Create(&applicationVOClone, additionalInfo)

	} else {
		errorDetails = (&applications.GoApplicationHandler{Context: cmd.Context()}).Create(&applicationVOClone, additionalInfo)
	}

	if len(errorDetails) > 0 {
		utils.DrawErrorTable(errorDetails)
		return errors.New(constants.ErroInvalidData)
	}

	emoji.Println("Application creation is now complete! You can view the application YAML file at ", filepath.Join(additionalInfo.PolicyCowConfig.PathConfiguration.DeclarativePath, "applications"), ", and the app connection codes are available inside ", additionalInfo.PolicyCowConfig.PathConfiguration.AppConnectionPath, ":smiling_face_with_sunglasses:")

	return nil

}
