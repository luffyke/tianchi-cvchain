package main

import (
    "fmt"
    "encoding/json"

    "github.com/pkg/errors"
    "github.com/hyperledger/fabric/bccsp"
    "github.com/hyperledger/fabric/bccsp/factory"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
    "github.com/hyperledger/fabric/protos/peer"
)

const DECKEY = "DECKEY"
const ENCKEY = "ENCKEY"
const IV = "IV"

type CvChain struct {
    bccspInst bccsp.BCCSP
}

type CV struct {
    Company   string `json:"company"`
    Position  string `json:"position"`
}

func (t *CvChain) Init(stub shim.ChaincodeStubInterface) peer.Response {
    return shim.Success(nil)
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
        tMap, err := stub.GetTransient()
        if err != nil {
            return shim.Error(fmt.Sprintf("Could not retrieve transient, err %s", err))
        }
        // result, err = enc(stub, args, tMap[ENCKEY], tMap[IV])
        return t.Enc(stub, args[0:], tMap[ENCKEY], tMap[IV])
    } else if fn == "decRecord" {
        tMap, err := stub.GetTransient()
        if err != nil {
            return shim.Error(fmt.Sprintf("Could not retrieve transient, err %s", err))
        }
        // result, err = dec(stub, args, tMap[DECKEY])
        return t.Dec(stub, args[0:], tMap[DECKEY])
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

func (t *CvChain) Enc(stub shim.ChaincodeStubInterface, args []string, encKey, IV []byte) peer.Response {
    // create the encrypter entity - we give it an ID, the bccsp instance, the key and (optionally) the IV
    ent, err := entities.NewAES256EncrypterEntity("ID", t.bccspInst, encKey, IV)
    if err != nil {
        //return "", fmt.Errorf("entities.NewAES256EncrypterEntity failed, err %s", err)
        return shim.Error(fmt.Sprintf("entities.NewAES256EncrypterEntity failed, err %s", err))
    }
    /*
    if len(args) != 2 {
        return shim.Error("Expected 2 parameters to function Encrypter")
    }
    */
    /*
	key := args[0]
    cleartextValue := []byte(args[1])
    */
    // key
    var key = args[0] + "-" + args[1]
    // value
    var cv = CV{Company: args[2], Position: args[3]}
    cvAsBytes, _ := json.Marshal(cv)

    // here, we encrypt cleartextValue and assign it to key
    err = encryptAndPutState(stub, ent, key, cvAsBytes)
    if err != nil {
        // return "", fmt.Errorf("encryptAndPutState failed: %s", err)
        return shim.Error(fmt.Sprintf("signEncryptAndPutState failed, err %+v", err))
	}
    // return string(cvAsBytes), nil
    return shim.Success(nil)
}

func (t *CvChain) Dec(stub shim.ChaincodeStubInterface, args []string, decKey []byte) peer.Response {
    // create the encrypter entity - we give it an ID, the bccsp instance, the key and (optionally) the IV
    ent, err := entities.NewAES256EncrypterEntity("ID", t.bccspInst, decKey, []byte(""))
    if err != nil {
        return shim.Error(fmt.Sprintf("entities.NewAES256EncrypterEntity failed, err %s", err))
        //return "", fmt.Errorf("entities.NewAES256EncrypterEntity failed, err %s", err)
    }
    /*
    if len(args) != 1 {
        return shim.Error("Expected 1 parameters to function Decrypter")
    }
    */
    // key
    var key = args[0] + "-" + args[1]
    // here we decrypt the state associated to key
    value, err := getStateAndDecrypt(stub, ent, key)
    if err != nil {
        return shim.Error(fmt.Sprintf("getStateAndDecrypt failed, err %+v", err))
        //return "", fmt.Errorf("getStateAndDecrypt failed: %s", err)
    }
    cv := CV{}
    json.Unmarshal(value, &cv)
    // here we return the decrypted value as a result
    // return cv.Company, nil
    return shim.Success([]byte(cv.Company))
}

// getStateAndDecrypt retrieves the value associated to key,
// decrypts it with the supplied entity and returns the result
// of the decryption
func getStateAndDecrypt(stub shim.ChaincodeStubInterface, ent entities.Encrypter, key string) ([]byte, error) {
    // at first we retrieve the ciphertext from the ledger
    ciphertext, err := stub.GetState(key)
    if err != nil {
        return nil, err
    }
    // GetState will return a nil slice if the key does not exist.
    // Note that the chaincode logic may want to distinguish between
    // nil slice (key doesn't exist in state db) and empty slice
    // (key found in state db but value is empty). We do not
    // distinguish the case here
    if len(ciphertext) == 0 {
        return nil, errors.New("no ciphertext to decrypt")
    }
    return ent.Decrypt(ciphertext)
}

// encryptAndPutState encrypts the supplied value using the
// supplied entity and puts it to the ledger associated to
// the supplied KVS key
func encryptAndPutState(stub shim.ChaincodeStubInterface, ent entities.Encrypter, key string, value []byte) error {
    // at first we use the supplied entity to encrypt the value
    ciphertext, err := ent.Encrypt(value)
    if err != nil {
        return err
    }
    return stub.PutState(key, ciphertext)
}

// main function starts up the chaincode in the container during instantiate
func main() {
    /*
    if err := shim.Start(new(CvChain)); err != nil {
        fmt.Printf("Error starting CvChain chaincode: %s", err)
    }
    */
    factory.InitFactories(nil)
    err := shim.Start(&CvChain{factory.GetDefault()})
    if err != nil {
        fmt.Printf("Error starting CvChain chaincode: %s", err)
    }
}