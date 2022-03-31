package metrics

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
)

// all necessary fuctions to read and parse the list of public-keys of Eth2 validator
func ReadCustomValidatorsFile(valFile string) (vals [][]byte, err error) {
	valList := NewValidatorList()

	// Read the file and parse the public keys
	jFile, err := os.Open(valFile)
	if err != nil {
		return vals, err
	}
	defer jFile.Close()

	jBytes, err := ioutil.ReadAll(jFile)
	if err != nil {
		return vals, err
	}

	json.Unmarshal(jBytes, &valList)

	// get the correct format for each of the
	keys, err := valList.Keys()
	if err != nil {
		return keys, err
	}
	return keys, nil
}

type ValidatorList []string

func NewValidatorList() ValidatorList {
	return make([]string, 0)
}

func (v *ValidatorList) Keys() ([][]byte, error) {
	valKeys := make([][]byte, 0)
	i := 0
	for idx, key := range *v {
		i = idx
		valKey, err := hexutil.Decode(key)
		if err != nil {
			return valKeys, err
		}
		valKeys = append(valKeys, valKey)
	}
	log.Infof("the list of custom validators had %d validator keys", i+1)
	return valKeys, nil
}
