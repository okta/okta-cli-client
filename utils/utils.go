package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func ExtractIDs(jsonData []byte) ([]string, error) {
	var items []map[string]interface{}
	if err := json.Unmarshal(jsonData, &items); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var ids []string
	for _, item := range items {
		if id, ok := item["id"].(string); ok {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func WriteFile(file *os.File, filePath, tmpFile string, data map[string]interface{}) error {
	defer file.Close()
	w := bufio.NewWriter(file)
	funcs := template.FuncMap{
		"join":      strings.Join,
		"lower":     strings.ToLower,
		"hasPrefix": strings.HasPrefix,
	}
	cmdTmplFilePath := fmt.Sprintf("%v/%v", filePath, tmpFile)
	tmpl, err := template.New(tmpFile).Funcs(funcs).ParseFiles(cmdTmplFilePath)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

func GetTagList(c <-chan orderedmap.Pair[string, *v3high.PathItem]) map[string]bool {
	tagSet := make(map[string]bool)
	for pair := range c {
		node := pair.Value()
		if node.Get != nil {
			for _, tag := range node.Get.Tags {
				tagSet[tag] = true
			}
		}
		if node.Put != nil {
			for _, tag := range node.Put.Tags {
				tagSet[tag] = true
			}
		}
		if node.Post != nil {
			for _, tag := range node.Post.Tags {
				tagSet[tag] = true
			}
		}
		if node.Delete != nil {
			for _, tag := range node.Delete.Tags {
				tagSet[tag] = true
			}
		}
		if node.Patch != nil {
			for _, tag := range node.Patch.Tags {
				tagSet[tag] = true
			}
		}
	}
	return tagSet
}

func PrettyPrintObject(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", " ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func PrettyPrintByte(v []byte) (err error) {
	var jsonMap interface{}
	err = json.Unmarshal(v, &jsonMap)
	if err != nil {
		return err
	}
	return PrettyPrintObject(jsonMap)
}

func GetPathParam(url string) (pathParams []string) {
	re := regexp.MustCompile(`{(.*?)}`)
	res := re.FindAllStringSubmatch(url, -1)
	for _, r := range res {
		pathParams = append(pathParams, r[1])
	}
	return
}

func FirstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

func ExtractMap(payload string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(payload), &m)
	return m, err
}

func PrepareDataForRestore(jsonData []byte) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse backup data: %w", err)
	}

	removeMetadataFields(data)

	cleanedData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal processed data: %w", err)
	}

	return cleanedData, nil
}

func removeMetadataFields(data map[string]interface{}) {
	fieldsToRemove := []string{
		"_links", "_embedded", "id", "created", "lastUpdated",
		"lastLogin", "passwordChanged", "status",
	}

	for _, field := range fieldsToRemove {
		delete(data, field)
	}

	for _, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			removeMetadataFields(nestedMap)
		} else if nestedSlice, ok := value.([]interface{}); ok {
			for _, item := range nestedSlice {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					removeMetadataFields(nestedMap)
				}
			}
		}
	}
}
