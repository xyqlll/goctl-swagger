package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

func Do(filename string, host string, basePath string, schemes string, in *plugin.Plugin) error {
	swagger, err := applyGenerate(in, host, basePath, schemes)
	if err != nil {
		fmt.Println(err)
	}
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")

	if err := enc.Encode(swagger); err != nil {
		fmt.Println(err)
	}

	output := in.Dir + "/" + filename

	data, err := GetOpenApi3Data(formatted.Bytes())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, data, 0666)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func GetOpenApi3Data(formatted []byte) ([]byte, error) {
	url := "https://converter.swagger.io/api/convert"
	method := "POST"

	payload := strings.NewReader(string(formatted))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "converter.swagger.io")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))

	return body, nil
}
