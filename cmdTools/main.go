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
		// if pair.Key() == "/api/v1/groups" {
		// 	pathParams := utils.GetPathParam(pair.Key())
		// 	node := pair.Value()
		// 	if node.Post != nil {
		// 		err = buildCmdForHTTPMethod(node.Post.OperationId, node.Post.Description, pair.Key(), http.MethodPost, node.Post.Tags, pathParams)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	if node.Get != nil {
		// 		err = buildCmdForHTTPMethod(node.Get.OperationId, node.Get.Description, pair.Key(), http.MethodGet, node.Get.Tags, pathParams)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }
		// if pair.Key() == "/api/v1/groups/{groupId}" {
		// 	pathParams := utils.GetPathParam(pair.Key())
		// 	node := pair.Value()
		// 	if node.Get != nil {
		// 		err = buildCmdForHTTPMethod(node.Get.OperationId, node.Get.Description, pair.Key(), http.MethodGet, node.Get.Tags, pathParams)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	if node.Put != nil {
		// 		err = buildCmdForHTTPMethod(node.Put.OperationId, node.Put.Description, pair.Key(), http.MethodPut, node.Put.Tags, pathParams)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	if node.Delete != nil {
		// 		err = buildCmdForHTTPMethod(node.Delete.OperationId, node.Delete.Description, pair.Key(), http.MethodDelete, node.Delete.Tags, pathParams)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }
		// if pair.Key() == "/api/v1/apps" {
		// 	node := pair.Value()
		// 	if node.Post != nil {
		// 		fmt.Println("130", node.Post.RequestBody.Content.Value("application/json").Schema.Schema().OneOf)
		// 		fmt.Println("131", node.Post.RequestBody.Content.Value("application/json").Schema.Schema().Discriminator.PropertyName)
		// 		fmt.Println("132", node.Post.RequestBody.Content.Value("application/json").Schema.Schema().Properties)
		// 		for _, v := range node.Post.RequestBody.Content.Value("application/json").Schema.Schema().OneOf {
		// 			fmt.Println("133", v.GetReference())
		// 			fmt.Println("134", v.Schema().AllOf)
		// 			if v.GetReference() == "#/components/schemas/WsFederationApplication" {
		// 				for _, v := range v.Schema().AllOf {
		// 					// fmt.Println("136", v.Schema().Discriminator)
		// 					// fmt.Println("137", v.Schema().SchemaTypeRef)
		// 					// fmt.Println("138", v.Schema().ExclusiveMaximum)
		// 					// fmt.Println("139", v.Schema().ExclusiveMinimum)
		// 					// fmt.Println("140", v.Schema().Type)
		// 					// fmt.Println("141", v.Schema().AllOf)
		// 					// fmt.Println("142", v.Schema().OneOf)
		// 					// fmt.Println("143", v.Schema().AnyOf)
		// 					// fmt.Println("144", v.Schema().Discriminator)
		// 					// fmt.Println("145", v.Schema().Examples)
		// 					// fmt.Println("146", v.Schema().PrefixItems)
		// 					// fmt.Println("147", v.Schema().Contains)
		// 					// fmt.Println("148", v.Schema().MinContains)
		// 					// fmt.Println("149", v.Schema().MaxContains)
		// 					// fmt.Println("150", v.Schema().If)
		// 					// fmt.Println("151", v.Schema().Else)
		// 					// fmt.Println("152", v.Schema().Then)
		// 					// fmt.Println("153", v.Schema().DependentSchemas)
		// 					// fmt.Println("154", v.Schema().PatternProperties)
		// 					// fmt.Println("155", v.Schema().PropertyNames)
		// 					// fmt.Println("156", v.Schema().UnevaluatedItems)
		// 					// fmt.Println("157", v.Schema().UnevaluatedProperties)
		// 					// fmt.Println("158", v.Schema().Items)
		// 					// fmt.Println("159", v.Schema().Anchor)
		// 					// fmt.Println("160", v.Schema().Not)
		// 					fmt.Println("161", v.Schema().Properties)
		// 					c := orderedmap.Iterate(context.Background(), v.Schema().Properties)
		// 					for pair := range c {
		// 						fmt.Println("165", pair.Key())
		// 						// fmt.Println("166", pair.Value().Schema().Type[0])
		// 					}
		// 					// fmt.Println("162", v.Schema().Title)
		// 					// fmt.Println("163", v.Schema().MultipleOf)
		// 					// fmt.Println("164", v.Schema().Maximum)
		// 					// fmt.Println("165", v.Schema().Minimum)
		// 					// fmt.Println("166", v.Schema().MaxLength)
		// 					// fmt.Println("167", v.Schema().MinLength)
		// 					// fmt.Println("168", v.Schema().Pattern)
		// 					// fmt.Println("169", v.Schema().Format)
		// 					// fmt.Println("170", v.Schema().MaxItems)
		// 					// fmt.Println("171", v.Schema().MinItems)
		// 					// fmt.Println("172", v.Schema().UniqueItems)
		// 					// fmt.Println("173", v.Schema().MaxProperties)
		// 					// fmt.Println("174", v.Schema().MinProperties)
		// 					// fmt.Println("175", v.Schema().Required)
		// 					// fmt.Println("176", v.Schema().Enum)
		// 					// fmt.Println("177", v.Schema().AdditionalProperties)
		// 					// fmt.Println("178", v.Schema().Description)
		// 					// fmt.Println("179", v.Schema().Default)
		// 					// fmt.Println("180", v.Schema().Const)
		// 					// fmt.Println("181", v.Schema().Nullable)
		// 					// fmt.Println("182", v.Schema().ReadOnly)
		// 					// fmt.Println("183", v.Schema().WriteOnly)
		// 					// fmt.Println("184", v.Schema().XML)
		// 					// fmt.Println("185", v.Schema().ExternalDocs)
		// 					// fmt.Println("186", v.Schema().Example)
		// 					// fmt.Println("187", v.Schema().Deprecated)
		// 					// fmt.Println("188", v.Schema().Extensions)
		// 					// fmt.Println("189", v.Schema().ParentProxy)
		// 				}
		// 			}

		// 			// fmt.Println(v.GetReferenceOrigin())
		// 			// fmt.Println(v.GetSchemaKeyNode())
		// 		}
		// 	}
		// }
		pathParams := utils.GetPathParam(pair.Key())
		node := pair.Value()
		if node.Post != nil {
			err = buildCmdForHTTPMethod(node.Post.OperationId, node.Post.Description, pair.Key(), http.MethodPost, node.Post.Tags, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Get != nil {
			err = buildCmdForHTTPMethod(node.Get.OperationId, node.Get.Description, pair.Key(), http.MethodGet, node.Get.Tags, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Put != nil {
			err = buildCmdForHTTPMethod(node.Put.OperationId, node.Put.Description, pair.Key(), http.MethodPut, node.Put.Tags, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Delete != nil {
			err = buildCmdForHTTPMethod(node.Delete.OperationId, node.Delete.Description, pair.Key(), http.MethodDelete, node.Delete.Tags, pathParams)
			if err != nil {
				return err
			}
		}
		if node.Patch != nil {
			err = buildCmdForHTTPMethod(node.Patch.OperationId, node.Patch.Description, pair.Key(), http.MethodPatch, node.Patch.Tags, pathParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func buildCmdForHTTPMethod(methodName, description, endpoint, httpMethod string, tags, pathParams []string) error {
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
	if httpMethod != http.MethodGet && httpMethod != http.MethodDelete {
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
	if httpMethod != http.MethodGet && httpMethod != http.MethodDelete {
		templateData["data"] = true
	}
	err = utils.WriteFile(f, "cmdTools", "lowLevelCmd.tmpl", templateData)
	if err != nil {
		return err
	}
	return nil
}
