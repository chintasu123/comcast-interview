package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type InputData struct {
	Number1 struct {
		N string `json:"N"`
	} `json:"number_1"`
	String1 struct {
		S string `json:"S"`
	} `json:"string_1"`
	String2 struct {
		S string `json:"S"`
	} `json:"string_2"`
	Map1 struct {
		M struct {
			Bool1 struct {
				BOOL string `json:"BOOL"`
			} `json:"bool_1"`
			Null1 struct {
				NULL string `json:"NULL "`
			} `json:"null_1"`
			List1 struct {
				L []struct {
					S    string `json:"S,omitempty"`
					N    string `json:"N,omitempty"`
					BOOL string `json:"BOOL,omitempty"`
					NULL string `json:"NULL,omitempty"`
				} `json:"L"`
			} `json:"list_1"`
		} `json:"M"`
	} `json:"map_1"`
	List2 struct {
		L string `json:"L"`
	} `json:"list_2"`
	List3 struct {
		L []string `json:"L"`
	} `json:"list_3"`
	Field7 struct {
		S string `json:"S"`
	} `json:""`
}

type OutputData struct {
	Number1 float64 `json:"number_1"`
	String1 string  `json:"string_1"`
	String2 int64   `json:"string_2"`
	Map1    struct {
		List1 []interface{} `json:"list_1"`
		Null1 interface{}   `json:"null_1"`
	} `json:"map_1"`
}

func sanitizeKey(key string) string {
	return strings.TrimSpace(key)
}

func transform(inputData InputData) OutputData {
	var outputData OutputData

	// Transform Number1
	outputData.Number1, _ = strconv.ParseFloat(inputData.Number1.N, 64)

	// Transform String1
	outputData.String1 = strings.TrimSpace(inputData.String1.S)

	// Transform String2
	string2Time, err := time.Parse(time.RFC3339, inputData.String2.S)
	if err == nil {
		outputData.String2 = string2Time.Unix()
	}

	// Transform Map1
	outputData.Map1.List1 = []interface{}{}
	for _, item := range inputData.Map1.M.List1.L {
		if item.N != "" {
			key := sanitizeKey(item.N)
			atoi, err := strconv.Atoi(key)
			if err != nil {
				continue
			}
			outputData.Map1.List1 = append(outputData.Map1.List1, atoi)
		} else if item.BOOL == "f" {
			key := sanitizeKey(item.BOOL)
			for _, value := range []string{"1", "t", "T", "TRUE", "true", "True"} {
				if key == value {
					outputData.Map1.List1 = append(outputData.Map1.List1, true)
					break
				}
			}
			for _, value := range []string{"0", "f", "F", "FALSE", "false", "False"} {
				if key == value {
					outputData.Map1.List1 = append(outputData.Map1.List1, false)
					break
				}
			}
		} else if item.NULL != "" {
			key := sanitizeKey(item.BOOL)
			for _, value := range []string{"1", "t", "T", "TRUE", "true", "True"} {
				if key == value {
					outputData.Map1.List1 = append(outputData.Map1.List1, nil)
					break
				}
			}
		} else if item.S != "" {
			key := sanitizeKey(item.S)
			parse, err := time.Parse(time.RFC3339, key)
			if err != nil {
				continue
			}
			outputData.Map1.List1 = append(outputData.Map1.List1, parse.Unix())
		}
	}
	if inputData.Map1.M.Null1.NULL != "" {
		outputData.Map1.Null1 = nil
	}

	return outputData
}

func main() {
	inputJSON := `{
	  "number_1": {
	    "N": "1.50"
	  },
	  "string_1": {
	    "S": "784498 "
	  },
	  "string_2": {
	    "S": "2014-07-16T20:55:46Z"
	  },
	  "map_1": {
	    "M": {
	      "bool_1": {
	        "BOOL": "truthy"
	      },
	      "null_1": {
	        "NULL ": "true"
	      },
	      "list_1": {
	        "L": [
	          {
	            "S": ""
	          },
	          {
	            "N": "011"
	          },
	          {
	            "N": "5215s"
	          },
	          {
	            "BOOL": "f"
	          },
	          {
	            "NULL": "0"
	          }
	        ]
	      }
	    }
	  },
	  "list_2": {
	    "L": "noop"
	  },
	  "list_3": {
	    "L": [
	      "noop"
	    ]
	  },
	  "": {
	    "S": "noop"
	  }
	}`

	var inputData InputData
	if err := json.Unmarshal([]byte(inputJSON), &inputData); err != nil {
		fmt.Println("Error:", err)
		return
	}

	outputData := []OutputData{transform(inputData)}

	outputJSON, err := json.Marshal(outputData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(outputJSON))
}
