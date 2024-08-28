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

func WriteFile(file *os.File, filePath, tmpFile string, data map[string]interface{}) error {
	defer file.Close()
	w := bufio.NewWriter(file)
	funcs := template.FuncMap{"join": strings.Join}
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
