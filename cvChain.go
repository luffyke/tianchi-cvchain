package main

import (
    "fmt"
    "encoding/json"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

type CvChain struct {
}

type CV struct {
    Company   string `json:"company"`
	Position  string `json:"position"`
}

func (t *CvChain) Init(stub shim.ChaincodeStubInterface) peer.Response {

}

func (t *CvChain) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
    fn, args := stub.GetFunctionAndParameters()
    var result string
    var err error
    if fn == "addRecord" {
        result, err = add(stub, args)
    } else if fn == "getRecord" {
        result, err = get(stub, args)
    } else if fn == "encRecord" {
        result, err = enc(stub, args)
    } else if fn == "decRecord" {
        result, err = dec(stub, args)
    }
    if err != nil {
        return shim.Error(err.Error())
    }
    // Return the result as success payload
    return shim.Success([]byte(result))
}

func add(stub shim.ChaincodeStubInterface, args []string) (string, error) {
    /*
    if len(args) != 2 {
        return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
    }
    */
    // key
    var key = args[0] + "-" + args[1]
    // value
    var cv = CV{Company: args[2], Position: args[3]}
    cvAsBytes, _ := json.Marshal(cv)
    err := stub.PutState(key, cvAsBytes)
    if err != nil {
        return "", fmt.Errorf("Failed to set asset: %s", key)
    }
    return string(cvAsBytes), nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
    /*
    if len(args) != 1 {
        return "", fmt.Errorf("Incorrect arguments. Expecting a key")
    }
    */
    // key
    var key = args[0] + "-" + args[1]
    value, err := stub.GetState(key)
    if err != nil {
        return "", fmt.Errorf("Failed to get asset: %s with error: %s", key, err)
    }
    if value == nil {
        return "", fmt.Errorf("Asset not found: %s", key)
    }
    cv := CV{}
	json.Unmarshal(value, &cv)
    return cv.Company, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
    if err := shim.Start(new(CvChain)); err != nil {
        fmt.Printf("Error starting CvChain chaincode: %s", err)
    }
}