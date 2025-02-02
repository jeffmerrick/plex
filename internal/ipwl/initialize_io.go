package ipwl

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/web3"
)

var (
	inputs           string
	scatteringMethod string
)

func InitializeIo(toolPath string, scatteringMethod string, inputVectors map[string][]string) ([]IO, error) {
	// Open the file and load its content
	tool, toolInfo, err := ReadToolConfig(toolPath)
	if err != nil {
		return nil, err
	}

	// Validate input keys
	err = validateInputKeys(inputVectors, tool.Inputs)
	if err != nil {
		return nil, err
	}

	// Handle scattering methods and create the ist
	var inputsList [][]string
	switch scatteringMethod {
	case "dotProduct":
		inputsList, err = dotProductScattering(inputVectors)
		if err != nil {
			return nil, err
		}
	case "crossProduct":
		inputsList, err = crossProductScattering(inputVectors)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid scattering method: %s", scatteringMethod)
	}

	var userId string

	if web3.IsValidEthereumAddress(os.Getenv("RECIPIENT_WALLET")) {
		userId = os.Getenv("RECIPIENT_WALLET")
	} else {
		fmt.Print("Invalid wallet address detected. Using empty string for user ID.\n")
		userId = ""
	}

	var ioList []IO

	for _, inputs := range inputsList {
		io, err := createSingleIo(inputs, tool, toolInfo, userId, inputVectors)
		if err != nil {
			return nil, err
		}
		ioList = append(ioList, io)
	}

	return ioList, nil
}

func validateInputKeys(inputVectors map[string][]string, toolInputs map[string]ToolInput) error {
	for inputKey := range inputVectors {
		if _, exists := toolInputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", toolInputs)
			return fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
		}
	}
	return nil
}

func dotProductScattering(inputVectors map[string][]string) ([][]string, error) {
	var vectorLength int
	for _, v := range inputVectors {
		if vectorLength == 0 {
			vectorLength = len(v)
			continue
		}
		if len(v) != vectorLength {
			return nil, fmt.Errorf("all input arguments must have the same length for dotProduct scattering method")
		}
	}

	var inputsList [][]string
	keys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i < vectorLength; i++ {
		tmp := []string{}
		for _, k := range keys {
			tmp = append(tmp, inputVectors[k][i])
		}
		inputsList = append(inputsList, tmp)
	}

	return inputsList, nil
}

func crossProductScattering(inputVectors map[string][]string) ([][]string, error) {
	cartesian := func(arrs ...[]string) [][]string {
		result := [][]string{{}}
		for _, arr := range arrs {
			var temp [][]string
			for _, res := range result {
				for _, str := range arr {
					product := append([]string{}, res...)
					product = append(product, str)
					temp = append(temp, product)
				}
			}
			result = temp
		}
		return result
	}

	keys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	arrays := make([][]string, len(inputVectors))
	for i, k := range keys {
		arrays[i] = inputVectors[k]
	}
	inputsList := cartesian(arrays...)

	return inputsList, nil
}

func createSingleIo(inputs []string, tool Tool, toolInfo ToolInfo, userId string, inputVectors map[string][]string) (IO, error) {
	io := IO{
		Tool:    toolInfo,
		Inputs:  make(map[string]FileInput),
		Outputs: make(map[string]Output),
		State:   "created",
		ErrMsg:  "",
		UserID:  userId,
	}

	inputKeys := make([]string, 0, len(inputVectors))
	for k := range inputVectors {
		inputKeys = append(inputKeys, k)
	}

	// Sort the inputKeys slice to ensure a consistent order
	sort.Strings(inputKeys)

	for i, inputValue := range inputs {
		inputKey := inputKeys[i]

		if strings.Count(inputValue, "/") == 1 {
			parts := strings.Split(inputValue, "/")
			cid := parts[0]
			fileName := parts[1]
			if !ipfs.IsValidCID(cid) {
				return io, fmt.Errorf("invalid CID: %s", cid)
			}
			io.Inputs[inputKey] = FileInput{
				Class:    tool.Inputs[inputKey].Type,
				FilePath: fileName,
				IPFS:     cid,
			}
		} else {
			cid, err := ipfs.WrapAndPinFile(inputValue)
			if err != nil {
				return io, err
			}
			io.Inputs[inputKey] = FileInput{
				Class:    tool.Inputs[inputKey].Type,
				FilePath: filepath.Base(inputValue),
				IPFS:     cid,
			}
		}
	}

	for outputKey, outputValue := range tool.Outputs {
		io.Outputs[outputKey] = FileOutput{
			Class:    outputValue.Type,
			FilePath: "",
			IPFS:     "",
		}
	}

	return io, nil
}
