package operation

import "encoding/json"

func ParseOperations(rawOperations string) []Operation {
	var operations []Operation

	var err = json.Unmarshal([]byte(rawOperations), &operations)
	if err != nil {
		panic(err)
	}

	return operations
}
