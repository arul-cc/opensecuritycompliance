package constants

const TaskMain = `// This file is autogenerated. Please do not modify
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"github.com/google/uuid"

	"cowlibrary/constants"
)

func handlePanic() {
	r := recover()
	if r != nil {
		os.WriteFile("logs.txt", debug.Stack(), os.ModePerm)
		os.WriteFile("task_output.json", []byte(` + "`" + `{"error":"` + "`" + `+fmt.Sprintf("%v, %s", r, "Please review the stack trace in the logs.txt file within the task.")+` + "`" + `"}` + "`" + `), os.ModePerm)
	}
}

func main() {
	defer handlePanic()
	inst := new(TaskInstance)

	taskInput := &TaskInputs{}
	taskOutput := &TaskOutputs{Outputs: &Outputs{}}
	errorOutput := make(map[string]string)
	
	inputObj := &TaskInputs{}
	if _, err := os.Stat("inputs.yaml"); err == nil {
		byts, err := os.ReadFile("inputs.yaml")
		if err == nil {
			byts = []byte(os.ExpandEnv(string(byts)))
			err = yaml.Unmarshal(byts, inputObj)
			if err != nil {
				taskInputObj := &TaskInputsV2{}
				err = yaml.Unmarshal(byts, taskInputObj)
				if err == nil {
					inputObj.SystemInputs = taskInputObj.SystemInputs
					inputObj.UserInputs = taskInputObj.UserInputs
					inputObj.FromDate_, _ = time.Parse("2006-01-02", taskInputObj.FromDate_)
					inputObj.ToDate_, _ = time.Parse("2006-01-02", taskInputObj.ToDate_)
				}
			}
		}
		taskInputByts, err := json.Marshal(inputObj)
		if err != nil {
			return
		}
		taskInputFilePath := "task_input.json"
		taskInputByts = []byte(os.ExpandEnv(string(taskInputByts)))
		err = os.WriteFile(taskInputFilePath, taskInputByts, os.ModePerm)
		if err != nil {
			errorOutput["error"] = err.Error()
			return
		}
	}

	taskInputFile := ""
	flag.StringVar(&taskInputFile, "f", "task_input.json", "task input file path")
	flag.Parse()

	taskOutputFile := ""
	i := strings.LastIndex(taskInputFile, "/")
	if i != -1 {
		taskOutputFile = taskInputFile[:i+1]
	}
	taskOutputFile += "task_output.json"

	err := readFromFile(taskInputFile, taskInput)
	if err != nil {
		errorOutput["error"] = err.Error()
		writeToFile(taskOutputFile, errorOutput)
		return
	}

	// INFO : Added for unit testing
	if len(taskInput.SystemObjects) == 0 {
		defaultSystemObjectByts := []byte(os.ExpandEnv(string(constants.SystemObjects)))
		json.Unmarshal(defaultSystemObjectByts, &taskInput.SystemObjects)
		metaDataTemplate := &MetaDataTemplate{}
		metaDataTemplate.ControlID = uuid.New().String()
		metaDataTemplate.PlanExecutionGUID = uuid.New().String()
		metaDataTemplate.RuleGUID = uuid.New().String()
		metaDataTemplate.RuleTaskGUID = uuid.New().String()
		taskInput.MetaData = metaDataTemplate
	}
	inst.TaskInputs = taskInput
	err = inst.{{TaskName}}(taskInput.UserInputs, taskOutput.Outputs)
	if err != nil {
		errorOutput["error"] = err.Error()
		writeToFile(taskOutputFile, errorOutput)
		return
		// panic(err.Error())
	}
	writeToFile(taskOutputFile, taskOutput)
}

func readFromFile(fileName string, dest interface{}) error {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, dest)
	if err != nil {
		return err
	}
	return nil
}

func writeToFile(fileName string, data interface{}) error {
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, payload, 0644)
	if err != nil {
		return err
	}
	return nil
}

type TaskInstance struct {
	*TaskInputs
}

`
const TaskServer = (`// This file is autogenerated. Please do not modify

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
)

const servicename = "{{TaskName}}"

func getVarname(err error) string {
	errorString := err.Error()
	index := strings.Index(errorString, "Inputs.")
	if index > 0 {
		errorString = errorString[index+len("Inputs."):]
		return errorString[:strings.Index(errorString, " ")]
	}
	return ""
}

func validate(taskInput *TaskInputs) (UserInputs, error) {
	input := UserInputs{}
	refSourceVal := reflect.ValueOf(&input)
	errorString := ""
	for {
		payload, _ := json.Marshal(taskInput.UserInputs)
		err := json.Unmarshal(payload, &input)
		if err == nil {
			break
		}
		varName := getVarname(err)
		errorString = errorString + varName + " is not type of " + refSourceVal.Elem().FieldByName(varName).Type().String() + ", "
		delete(taskInput.UserInputs, varName)
	}
	if len(errorString) > 0 {
		return input, fmt.Errorf("[ " + errorString[:len(errorString)-2] + " ]")
	}
	return input, nil
}

func (inst *TaskInstance) Validate(taskInput *TaskInputs, taskOutput *TaskOutputs) error {
	_, err := validate(taskInput)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (inst *TaskInstance) Execute(taskInput *TaskInputs, taskOutput *TaskOutputs) (err error) {
	newInst := new(TaskInstance)
	defer func() {
		if r := recover(); r != nil {
			errorString := fmt.Sprintf("Error while executing task:%v", r)
			err = fmt.Errorf(errorString)
			fileName := "error.log"
			logFile, fileOpenErr := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if fileOpenErr == nil {
				defer logFile.Close()
				logFile.WriteString(errorString + "\n")
				logFile.Write(debug.Stack())
				logFile.WriteString("\n\n")
			}
		}
	}()

	output := Outputs{}
	input, err := validate(taskInput)
	if err != nil {
		fmt.Println(err)
		return err
	}
	payload, _ := json.Marshal(taskInput.UserInputs)
	err = json.Unmarshal(payload, &input)
	fmt.Println(err)
	if err != nil {
		return err
	}
	newInst.SystemInputs = &taskInput.SystemInputs
	err = newInst.{{TaskName}}(&input, &output)
	fmt.Println(err)
	if err != nil {
		return err
	}
	payload, err = json.Marshal(&output)
	fmt.Println(err)
	if err != nil {
		return err
	}
	err = json.Unmarshal(payload, &(taskOutput.Outputs))
	return nil
}


`)

var TaskServerStructWithPackage = `// This file is autogenerated. Please do not modify
package main
import (
	"time"
	{{import()}}
)

` + TaskServerStructs

var TaskServerStructs = `

const (
	appObject    = "app"
	serverObject = "server"
)

const (
	_ = iota
	schemaInputs
	schemaFacts
	schemaOutputs
)

type TaskInputs struct {
	SystemInputs ` + "`" + `yaml:",inline"` + "`" + `
	UserInputs   *UserInputs ` + "`" + `yaml:"userInputs"` + "`" + `
	FromDate_     time.Time ` + "`" + `yaml:"fromDate,omitempty"` + "`" + `
    ToDate_       time.Time ` + "`" + `yaml:"toDate,omitempty"` + "`" + `
}

type TaskInputsV2 struct {
	SystemInputs ` + "`" + `yaml:",inline"` + "`" + `
	UserInputs   *UserInputs ` + "`" + `yaml:"userInputs"` + "`" + `
	FromDate_     string ` + "`" + `yaml:"fromDate,omitempty"` + "`" + `
    ToDate_       string ` + "`" + `yaml:"toDate,omitempty"` + "`" + `
}

type SystemInputs struct {
	UserObject    *ObjectTemplate ` + "`" + `yaml:"userObject"` + "`" + `
	SystemObjects []*ObjectTemplate ` + "`" + `yaml:"systemObjects"` + "`" + `
	MetaData      *MetaDataTemplate` + "`" + `yaml:"-"` + "`" + `
}

type ObjectTemplate struct {
	App         *AppAbstract  ` + "`" + `yaml:"app,omitempty"` + "`" + `
	Server      *ServerAbstract ` + "`" + `yaml:"server,omitempty"` + "`" + `
	Credentials []*Credential   ` + "`" + `yaml:"credentials,omitempty"` + "`" + `
}

type MetaDataTemplate struct {
	RuleGUID          string
	RuleTaskGUID      string
	ControlID         string
	PlanExecutionGUID string
}

type TaskOutputs struct {
	Outputs *Outputs
}

type AppAbstract struct {
	*AppBase  ` + "`" + `yaml:",inline"` + "`" + `
	ID          string               ` + "`" + `json:"id,omitempty"` + "`" + `
	AppSequence int                 ` + "`" + `json:"appSequence,omitempty"` + "`" + `
	AppTags     map[string][]string ` + "`" + `json:"appTags,omitempty"` + "`" + `
	ActionType  string              ` + "`" + `json:"actionType,omitempty"` + "`" + `
	AppObjects  map[string]interface{}
	Servers     []*ServerAbstract ` + "`" + `json:"servers,omitempty"` + "`" + `
	UserDefinedCredentials {{UserDefinedCredentials}} ` + "`" + `json:"userDefinedCredentials" yaml:"userDefinedCredentials"` + "`" + `
	LinkedApplications  map[string][]*AppAbstract ` + "`" + `json:"linkedApplications,omitempty"` + "`" + `
}

type AppBase struct {
	ApplicationName string                 ` + "`" + `json:"appName,omitempty" yaml:"name" binding:"required" validate:"required"` + "`" + `
	ApplicationGUID string                 ` + "`" + `json:"applicationguid,omitempty" yaml:"applicationguid,omitempty"` + "`" + `
	AppGroupGUID    string                 ` + "`" + `json:"appgroupguid,omitempty" yaml:"appgroupguid,omitempty"` + "`" + `
	ApplicationURL  string                 ` + "`" + `json:"AppURL,omitempty" yaml:"appURL,omitempty"` + "`" + `
	ApplicationPort string                 ` + "`" + `json:"Port,omitempty" yaml:"appPort,omitempty"` + "`" + `
	OtherInfo       map[string]interface{} ` + "`" + `yaml:"otherinfo,omitempty"` + "`" + `
}

type ServerBase struct {
	ServerGUID      string
	ServerName      string ` + "`" + `json:"servername,omitempty"` + "`" + `
	ApplicationGUID string ` + "`" + `json:"appid,omitempty"` + "`" + `
	ServerType      string ` + "`" + `json:"servertype,omitempty"` + "`" + `
	ServerURL       string ` + "`" + `json:"serverurl,omitempty"` + "`" + `
	ServerHostName  string ` + "`" + `json:"serverhostname,omitempty"` + "`" + `
}

type ServerAbstract struct {
	ServerBase  ` + "`" + `yaml:",inline"` + "`" + `
	ID            string              ` + "`" + `json:"id,omitempty"` + "`" + `
	ServerTags    map[string][]string ` + "`" + `json:"servertags,omitempty"` + "`" + `
	ServerBootSeq int                 ` + "`" + `json:"serverbootseq,omitempty"` + "`" + `
	ActionType    string              ` + "`" + `json:"actiontype,omitempty"` + "`" + `
	OSInfo        struct {
		OSDistribution string ` + "`" + `json:"osdistribution,omitempty"` + "`" + `
		OSKernelLevel  string ` + "`" + `json:"oskernellevel,omitempty"` + "`" + `
		OSPatchLevel   string ` + "`" + `json:"ospatchlevel,omitempty"` + "`" + `
	} ` + "`" + `json:"osinfo,omitempty"` + "`" + `
	IPv4Addresses map[string]string ` + "`" + `json:"ipv4addresses,omitempty"` + "`" + `
	Volumes       map[string]string ` + "`" + `json:"volumes,omitempty"` + "`" + `
	OtherInfo     struct {
		CPU      int ` + "`" + `json:"cpu,omitempty"` + "`" + `
		GBMemory int ` + "`" + `json:"memory_gb,omitempty"` + "`" + `
	} ` + "`" + `json:"otherinfo,omitempty"` + "`" + `
	ClusterInfo struct {
		ClusterName    string            ` + "`" + `json:"clustername,omitempty"` + "`" + `
		ClusterMembers []*ServerAbstract ` + "`" + `json:"clustermembers,omitempty"` + "`" + `
	} ` + "`" + `json:"clusterinfo,omitempty"` + "`" + `
	Servers []*ServerAbstract ` + "`" + `json:"servers,omitempty"` + "`" + `
}

// Credential : Holds Customer Credentials
type Credential struct {
	CredentialBase  ` + "`" + `yaml:",inline"` + "`" + `
	ID            string                 ` + "`" + `json:"id,omitempty" yaml:"id,omitempty"` + "`" + `
	PasswordHash  []byte                 ` + "`" + `json:"passwordhash,omitempty" yaml:"passwordhash,omitempty"` + "`" + `
	Password      string                 ` + "`" + `json:"passwordstring,omitempty" yaml:"password,omitempty"` + "`" + `
	LoginURL      string                 ` + "`" + `json:"loginurl,omitempty" yaml:"loginURL,omitempty" binding:"required,url" validate:"required,url"` + "`" + `
	SSHPrivateKey []byte                 ` + "`" + `json:"sshprivatekey,omitempty" yaml:"sshprivatekey,omitempty"` + "`" + `
	CredTags      map[string][]string    ` + "`" + `json:"credtags,omitempty" yaml:"tags,omitempty"` + "`" + `
	OtherCredInfo map[string]interface{} ` + "`" + `json:"othercredinfomap,omitempty" yaml:"otherCredentials" binding:"required" validate:"required"` + "`" + `
}


type CredentialBase struct {
    CredGUID   string ` + "`" + `json:"credguid,omitempty" yaml:"credguid,omitempty"` + "`" + `
    CredType   string ` + "`" + `json:"credtype,omitempty" yaml:"credType,omitempty"` + "`" + `
    SourceGUID string ` + "`" + `json:"sourceguid,omitempty" yaml:"sourceguid,omitempty"` + "`" + `
    SourceType string ` + "`" + `json:"sourcetype,omitempty" yaml:"sourcetype,omitempty"` + "`" + `
    UserID     string ` + "`" + `json:"userID,omitempty" yaml:"userid,omitempty"` + "`" + `
}


`

const TaskServiceStructs = `// This file is autogenerated. Modify as per your task needs.
package main

type UserInputs struct {
	BucketName string
}

type Outputs struct {
	ComplianceStatus_ string
	CompliancePCT_    int
	LogFile           string
}
`

const TaskServiceStructs_V2 = `package main

type UserInputs struct {
	{{replace_input_fields}}
}

type Outputs struct {
	ComplianceStatus_          string 
	CompliancePCT_             int   
	{{replace_output_fields}}
}
`

const TaskService = `
package main

import (
	cowStorage "appconnections/minio"
	"cowlibrary/vo"
	"encoding/json"
	"fmt"
	"time"
)

// {{TaskName}} :
func (inst *TaskInstance) {{TaskName}}(inputs *UserInputs, outputs *Outputs) (err error) {

	compliancePCT, complianceStatus := 0, "NON_COMPLIANT"
	output := map[string]interface{}{
		"Outputs": map[string]string{
			"Status": "Task executed successfully",
		},
	}
	defer func() {
		err = func() error {

			systemInputs := vo.SystemInputs{}
			systemInputsByteData, err := json.Marshal(inst.SystemInputs)
			if err != nil {
				return err
			}
			err = json.Unmarshal(systemInputsByteData, &systemInputs)
			if err != nil {
				return err
			}

			outputs.LogFile, err = cowStorage.UploadJSONFile(fmt.Sprintf("%v-%v%v", "log-", time.Now().Unix(), ".json"), output, systemInputs)
			if err != nil {
				return err
			}

			outputs.CompliancePCT_ = compliancePCT
			outputs.ComplianceStatus_ = complianceStatus

			return nil
		}()
	}()
	
	{{ValidationMethod}}
	// TODO : PLACEHOLDER FOR RULE BUSINESS LOGIC

	return nil
}



`

const TaskService_V2 = `package main

import (
	{{replace_with_imports}}
    
	
	/*
	"bytes"
	
	"encoding/csv"
	"os"
	"io"
	"fmt"
	"encoding/json"
	"log"
	
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/minio/minio-go" */
)



func (inst *TaskInstance) {{replace_method_name}}(inputs *UserInputs, outputs *Outputs) (err error) {
    

	{{replace_final_code}}

	if err !=nil{
		return err
	}
	outputs.CompliancePCT_ = 100
	outputs.ComplianceStatus_ = "Compliant"


    return nil

}




{{replace_methods}}


`

const FileStore = `// This file is autogenerated. Please do not modify
package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dmnlk/stringUtils"
	"github.com/minio/minio-go"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/operations"
)

func (inst *TaskInstance) registerMinio(inputs *UserInputs) (*minio.Client, string, error) {
	var endpoint, accessKey, bucketName, secretKey string
	var err error

	for _, systemInput := range inst.SystemObjects {
		switch {
		case systemInput.App != nil: // Application object
			if systemInput.App.ApplicationName == "minio" {
				for _, cred := range systemInput.Credentials {
					if cred.LoginURL != "" && cred.CredTags["servicename"] != nil && cred.CredTags["servicetype"] != nil {
						endpoint = cred.LoginURL
						accessKey = cred.OtherCredInfo["MINIO_ACCESS_KEY"].(string)
						secretKey = cred.OtherCredInfo["MINIO_SECRET_KEY"].(string)
						bucketName, _ = cred.OtherCredInfo["BucketName"].(string)
						break
					}
				}

			}
		case systemInput.Server != nil: // Server Object
		}
	}
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKey, secretKey, false)
	if err != nil {
		log.Printf("%v", err)
		return nil, "", err
	}

	log.Printf("%#v", minioClient) // minioClient is now setup

	if stringUtils.IsEmpty(bucketName) {
		bucketName = inputs.BucketName
	}

	exists, err := minioClient.BucketExists(bucketName)

	if !exists && err == nil {
		// Bucket does not exist. Create one
		err = minioClient.MakeBucket(bucketName, "")
		if err != nil {
			log.Printf("%v", err)
			return nil, "", err
		}
	}
	if err == nil && exists {
		log.Printf("We already own %s", bucketName)
	}

	return minioClient, endpoint, err
}

func (inst *TaskInstance) createAndUploadFile(data interface{}, fileName, bucketName, endpoint string,
	minioClient *minio.Client) (string, string, error) {
	folderName := inst.getFolderName()
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", "", err
	}
	tempFileName := fmt.Sprintf("%v%v", strings.Replace(folderName, "/", "-", -1), fileName)
	err = os.WriteFile(tempFileName, payload, 0644)
	if err != nil {
		return "", "", err
	}
	if err = inst.uploadFile(minioClient, bucketName, folderName+fileName, tempFileName); err != nil {
		return "", "", err
	}
	fileName = "http://" + endpoint + "/" + bucketName + "/" + folderName + fileName
	defer os.Remove(tempFileName)
	return fileName, folderName, nil
}

func (inst *TaskInstance) getFolderName() string {
	return inst.getHash(inst.MetaData.PlanExecutionGUID, inst.MetaData.ControlID, inst.MetaData.RuleGUID, inst.MetaData.RuleTaskGUID)
}

func (inst *TaskInstance) getHash(values ...string) string {
	hash := ""
	for _, value := range values {
		h := sha1.New()
		h.Write([]byte(value))
		hash = hash + fmt.Sprintf("%x/", h.Sum(nil))
	}
	log.Println("Hash of input:", hash, "length:", len(hash))
	return hash
}

func (inst *TaskInstance) uploadFile(minioClient *minio.Client, bucketName string, objectName string, fileName string) (err error) {
	contentType := "application/json"

	// Upload the file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, fileName, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	return nil
}

func (inst *TaskInstance) downloadFile(minioClient *minio.Client, bucketName string, hash string, fileName string) (err error) {
	length := len(hash)
	if length > 0 && hash[length-1] != '/' {
		hash += "/"
	}
	objectName := hash + fileName

	// Download the file with FGetObject
	err = minioClient.FGetObject(bucketName, objectName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	log.Printf("Successfully downloaded %s\n", objectName)

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	fmt.Printf("File Name = %s /t Size = %d", fileInfo.Name(), fileInfo.Size())
	return nil
}

func (inst *TaskInstance) downloadFileV2(storageName string, minioClient *minio.Client,
	bucketName, objectName, hash, suffixFileName string) (string, error) {
	folderName := strings.Replace(inst.getFolderName(), "/", "-", -1)
	dataFileName := fmt.Sprintf("%v%v-%v-%v", folderName, time.Now().Unix(), rand.Int31(), suffixFileName)
	if strings.HasPrefix(objectName, "networkfile://") {
		objectName = strings.Replace(objectName, "networkfile://", "", 1)
		err := inst.downloadFileWithRclone(storageName, objectName, dataFileName)
		if err != nil {
			return "", err
		}
	} else {
		if len(hash) > 0 {
			if hash[len(hash)-1] != '/' {
				hash += "/"
			}
			objectName = hash + objectName
		}
		_, err := inst.downloadFileFromMinio(minioClient, bucketName, objectName, dataFileName)
		if err != nil {
			return "", err
		}
	}
	return dataFileName, nil
}

func (inst *TaskInstance) downloadFileAndUnmarshalV2(storageName string, minioClient *minio.Client,
	bucketName, objectName, hash string, object interface{}) error {
	dataFileName, err := inst.downloadFileV2(storageName, minioClient, bucketName, objectName, hash, "data.json")
	if err != nil {
		return err
	}
	defer os.Remove(dataFileName)
	return inst.unmarshalFile(dataFileName, object)
}

func (inst *TaskInstance) downloadFileWithRclone(storageName, objectName, dataFileName string) error {

	fsrc, err := fs.NewFs(storageName+":")
	if err != nil {
		return err
	}
	fdst, err := fs.NewFs("./")
	if err != nil {
		return err
	}

	return operations.CopyFile(context.Background(), fdst, fsrc, dataFileName, objectName)
}

func (inst *TaskInstance) downloadFileFromMinio(minioClient *minio.Client, bucketName, objectName, fileName string) (outFileName string, err error) {

	// Download the file with FGetObject
	err = minioClient.FGetObject(bucketName, objectName, fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}

	return fileInfo.Name(), nil
}

func (inst *TaskInstance) unmarshalFile(dataFile string, object interface{}) (err error) {
	dataPayload, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataPayload, object)
}

func (inst *TaskInstance) downloadFileAndUnmarshal(minioClient *minio.Client, bucketName, objectName, folderName string, object interface{}) (err error) {
	dataFile := fmt.Sprintf("%v%v-%v%v", strings.Replace(folderName, "/", "-", -1), "data", time.Now().Unix(), ".json")
	_, err = inst.downloadFileFromMinio(minioClient, bucketName, objectName, dataFile)
	if err != nil {
		return err
	}
	defer os.Remove(dataFile)
	return inst.unmarshalFile(dataFile, object)
}

func readFileHelperWithExtension(fileName, extension string, target interface{}) {
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	if !strings.HasSuffix(fileName, extension) {
		fileName += extension
	}
	readFileHelper(fileName, target, 4, 0)
}

func readFileHelper(fileName string, target interface{}, nestedLevelLimit int, count int) {
	// if !strings.HasSuffix(fileName, ".json") {
	// 	fileName += ".json"
	// }
	if nestedLevelLimit < count {
		return
	}
	count++
	if !strings.Contains(fileName, string(os.PathSeparator)) {
		fileName = "files" + string(os.PathSeparator) + fileName
	} else {
		fileName = ".." + string(os.PathSeparator) + fileName
	}

	fs, err := os.Stat(fileName)
	if err != nil || fs.IsDir() {
		readFileHelper(fileName, target, nestedLevelLimit, count)
		return
	}

	byts, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("error while reading " + fileName)
		return
	}
	json.Unmarshal(byts, target)
}



`

const MetaDataValueJSON = `{
	"RuleGUID": "",
	"RuleTaskGUID": "",
	"ControlID": "",
	"PlanExecutionGUID": ""
}`

const SystemObjectYAML = `SystemObjects:
- App:
  appName: minio
  appurl: continube.com
  labels:
  app:
	- minio
  appURL: cowstorage:9000
  appPort:
  UserDefinedCredentials:
	MINIO_ACCESS_KEY: '$MINIO_ROOT_USER'
	MINIO_SECRET_KEY: '$MINIO_ROOT_PASSWORD'`

const SystemObjectValueJSON = `[
	{
		"App": {
			"ApplicationName": "minio",
			"ApplicationGUID": "",
			"AppGroupGUID": "",
			"ApplicationURL": "minio.domain.com",
			"ID": "4",
			"AppTags": {
				"app": [
					"minio"
				]
			}
		},
		"Credentials": [
			{
				"CredGUID": "",
				"SourceGUID": "",
				"SourceType": "app",
				"ID": "1",
				"Password": "",
				"LoginURL": "cowstorage:9000",
				"CredTags": {
					"servicename": [
						"minio"
					],
					"servicetype": [
						"storage"
					]
				},
				"OtherCredInfo": {
					"MINIO_ACCESS_KEY": "",
					"MINIO_SECRET_KEY": ""
				}
			}
		]
	}
]`

const TaskInputValueJSON = `{
    "BucketName": "demo"
}`

const UserObjectAppValueJSON = `{
	"App": {
		"AppTags": {
			"app": []
		}
	},
	"Credentials": [
		{
			"OtherCredInfo": {}
		}
	]
}`

const UserObjectServerValueJSON = `{
	"Server": {
		"servername": "",
		"servertags": {
			"app": []
		},
		"osinfo": {},
		"otherinfo": {},
		"clusterinfo": {}
	},
	"Credentials": [
		{
			"OtherCredInfo": {}
		}
	]
}`

const AutoGeneratedFilePrefix = "autogenerated_"

const TaskInputYAML = `
userObject:
  name: ""
  appURL: ""
  appPort: 0
  userDefinedCredentials:
userInputs:
  BucketName: demo
fromDate:
toDate:
`
const MetaYAML = `authors:
- username@gmail.com
domain: continube
createdDate: {{CreatedDate}}
name: {{TaskName}}
displayName: {{TaskName}}
version: '1.0'
description: {{TaskName}}
shaToken: ''
showInCatalog: true
icon: fa-solid fa-database
type: {{Type}}
tags:
- generic
applicationType: generic
userObjectJSONInBase64: ''
systemObjectJSONInBase64: ''
inputs:
- name: BucketName # A unique identifier of the task input
  description: minio bucket name for the process # A concise description of the task input
  type: string # data type of the input. Available Types STRING, INT, FLOAT, FILE
  allowedValues: [] # Optional. Specifies allowed values for the input. Use a comma-separated list for multiple values.
  defaultValue: demo # optional. You can specify the default value (either a string or a number) at this point, for now, it supports a single value.
  showField: true  # boolean: true | false
  isRequired: true # boolean: true | false
outputs:
- name: ComplianceStatus_
  description: compliance status of the task
  type: string
- name: CompliancePCT_
  description: compliance percentage of the task
  type: int
`
