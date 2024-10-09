package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/okta/okta-cli-client/utils"

	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	oasv3Spec, err := os.ReadFile("template.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}

	document, err := libopenapi.NewDocument(oasv3Spec)
	if err != nil {
		panic(fmt.Sprintf("cannot create new document: %e", err))
	}

	docModel, errors := document.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = generateCmd(ctx, docModel)
	if err != nil {
		fmt.Println(err.Error())
	}
}

const (
	packageName = "okta"
)

func generateCmd(ctx context.Context, docModel *libopenapi.DocumentModel[v3high.Document]) (err error) {
	c := orderedmap.Iterate(ctx, docModel.Model.Paths.PathItems)
	listFileName := utils.GetTagList(c)
	err = createFileWithDefaultTemplate(listFileName)
	if err != nil {
		return err
	}
	c = orderedmap.Iterate(ctx, docModel.Model.Paths.PathItems)
	err = buildCmdFile(c)
	if err != nil {
		return err
	}
	return nil
}

func createFileWithDefaultTemplate(listFileName map[string]bool) error {
	for fileName := range listFileName {
		filePath := fmt.Sprintf("%v/%vCmd.go", packageName, fileName)
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		data := map[string]interface{}{
			"packageName":   packageName,
			"name":          fileName,
			"nameLowerCase": utils.FirstToLower(fileName),
		}
		err = utils.WriteFile(f, "cmdTools", "highLevelCmd.tmpl", data)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildCmdFile(c <-chan orderedmap.Pair[string, *v3high.PathItem]) error {
	var err error
	for pair := range c {
		pathParams := utils.GetPathParam(pair.Key())
		node := pair.Value()
		if node.Post != nil {
			err = buildCmdForHTTPMethod(node.Post, pair.Key(), http.MethodPost, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Get != nil {
			err = buildCmdForHTTPMethod(node.Get, pair.Key(), http.MethodGet, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Put != nil {
			err = buildCmdForHTTPMethod(node.Put, pair.Key(), http.MethodPut, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Delete != nil {
			err = buildCmdForHTTPMethod(node.Delete, pair.Key(), http.MethodDelete, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Patch != nil {
			err = buildCmdForHTTPMethod(node.Patch, pair.Key(), http.MethodPatch, pathParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func buildCmdForHTTPMethod(ops *v3high.Operation, endpoint, httpMethod string, pathParams []string) error {
	methodName := ops.OperationId
	description := ops.Description
	tags := ops.Tags
	var fileName string
	if len(tags) != 1 {
		return fmt.Errorf("get multiple tags for end point %v method %v", endpoint, httpMethod)
	} else {
		fileName = tags[0]
	}
	f, err := os.OpenFile(fmt.Sprintf("%v/%vCmd.go", packageName, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	subCommand := strings.ReplaceAll(methodName, fileName, "")

	sanitizedOperationID := cases.Title(language.English, cases.NoLower).String(methodName)
	sanitizedPathParams := make([]string, 0)
	for _, pathParam := range pathParams {
		sanitizedPathParams = append(sanitizedPathParams, fmt.Sprintf("%v%v", sanitizedOperationID, pathParam))
	}
	requiredFlags := make([]string, 0)
	requiredFlags = append(requiredFlags, pathParams...)
	if checkRequestBodyExist(ops) {
		requiredFlags = append(requiredFlags, "data")
	}

	templateData := map[string]interface{}{
		"name":          fileName,
		"operationId":   sanitizedOperationID,
		"description":   description,
		"pathParams":    sanitizedPathParams,
		"requiredFlags": requiredFlags,
		"subCommand":    subCommand,
	}
	if checkRequestBodyExist(ops) {
		templateData["data"] = true
	}

	err = utils.WriteFile(f, "cmdTools", "lowLevelCmd.tmpl", templateData)
	if err != nil {
		return err
	}
	return nil
}

func checkRequestBodyExist(ops *v3high.Operation) bool {
	return ops.RequestBody != nil
}
