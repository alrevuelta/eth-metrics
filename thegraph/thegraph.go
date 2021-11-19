package thegraph

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	//log "github.com/sirupsen/logrus"
)

// TODO:
// - Allow to query up to a given block
// - Implement some caching features, perhaps by lastpage token/id

// See attestantio graphs
var ENDPOINTS = map[string]string{
	"mainnet": "https://api.thegraph.com/subgraphs/name/attestantio/eth2deposits",
	"prater":  "https://api.thegraph.com/subgraphs/name/attestantio/eth2deposits-pater",
	"pyrmont": "https://api.thegraph.com/subgraphs/name/attestantio/eth2deposits-pyrmont",
}

// Max results per page that thegraph allows. More than it
// will result in an error
const MAX_RESULTS_PAGE = uint64(1000)

// See https://thegraph.com/docs/developer/graphql-api
var QUERY_SCHEMA = `{"query":
    "{deposits(first: __first__, where: {id_gt: \"__lastid__\", withdrawalCredentials_in : [__withcredlist__] }) { validatorPubKey, id }}",
    "variables":null
}`

type Thegraph struct {
	network  string
	endpoint string
	withCred []string
}

// Response structure for the previous query
type Deposit struct {
	ValidatorPubKey string `json:"validatorPubKey"`
	Id              string `json:"id"`
}

type Deposits struct {
	Deposits []Deposit `json:"deposits"`
}

type Data struct {
	Data Deposits `json:"data"`
}

func NewThegraph(network string, withCredList []string) (*Thegraph, error) {

	if len(withCredList) == 0 {
		return nil, errors.New("at least one withdrawal credential must be provided")
	}

	// Verify that the withdrawal credentials are valid
	for _, cred := range withCredList {
		_, err := ValidateAndDecodeWithdrawalCredentials(cred)
		if err != nil {
			return nil, errors.Wrap(err, "Withdrawal credential is not valid: "+cred)
		}
	}

	// Verify that the network is supported
	graphEndpoint, exists := ENDPOINTS[network]
	if !exists {
		return nil, errors.New("Network not supported: " + network)
	}

	return &Thegraph{
		network:  network,
		withCred: withCredList,
		endpoint: graphEndpoint,
	}, nil
}

// Get the deposited keys in the eth2 deposit contract by withdrawal credentials
func (a *Thegraph) GetDepositedKeys() ([][]byte, error) {
	depositedKeys := make([][]byte, 0)

	done := false

	id := "0x0000000000000000"
	for !done {
		pageKeys, startId, err := getKeysPage(id, a.endpoint, a.withCred)
		if err != nil {
			return nil, errors.Wrap(err, "error getting page")
		}
		id = startId
		if uint64(len(pageKeys)) < MAX_RESULTS_PAGE {
			done = true
		}
		depositedKeys = append(depositedKeys, pageKeys...)
	}
	return depositedKeys, nil
}

// TheGraph is limited to 1000 results per page, so if we have more keys
// we need to run multiple calls. Page starts at 0
func getKeysPage(id string, graphEndpoint string, withCredentials []string) ([][]byte, string, error) {
	var quotedWithCredentials []string
	for _, cred := range withCredentials {
		// Each with credential is surrounded by escaped quotes \"withcred\"
		temp := strings.Replace(strconv.Quote(cred), "\"", `\"`, -1)
		quotedWithCredentials = append(quotedWithCredentials, temp)
	}

	query := strings.Replace(QUERY_SCHEMA, "__withcredlist__", strings.Join(quotedWithCredentials, ","), -1)
	query = strings.Replace(query, "__first__", strconv.FormatUint(MAX_RESULTS_PAGE, 10), -1)
	query = strings.Replace(query, "__lastid__", id, -1)

	var jsonStr = []byte(query)
	req, err := http.NewRequest("POST", graphEndpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, "", errors.Wrap(err, "could not create request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 10
	resp, err := client.Do(req)

	if err != nil {
		return nil, "", errors.Wrap(err, "could not send request")
	}

	defer resp.Body.Close()

	response := &Data{}

	if resp.Status != "200 OK" {
		return nil, "", errors.New("the http response was different than 200, " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, response); err != nil {
		return nil, "", errors.Wrap(err, "could not unmarshal the body of the response")
	}

	// Detect for errors in the response
	if strings.Contains(string(body), "errors") {
		return nil, "", errors.New("there was an error in the TheGraph call: " + string(body))
	}

	pageKeys := make([][]byte, 0)
	for _, s := range response.Data.Deposits {
		decKey, err := hexutil.Decode(s.ValidatorPubKey)
		if err != nil {
			return nil, "", errors.Wrap(err, "could not decode public validator key")
		}
		pageKeys = append(pageKeys, decKey)
	}

	return pageKeys, response.Data.Deposits[len(response.Data.Deposits)-1].Id, nil
}

// Makes sure the withdrawal credentials comply with:
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/validator.md#withdrawal-credentials
func ValidateAndDecodeWithdrawalCredentials(withCred string) ([]byte, error) {
	if !strings.HasPrefix(withCred, "00") { // BLS_WITHDRAWAL_PREFIX
		if !strings.HasPrefix(withCred, "01") { // ETH1_ADDRESS_WITHDRAWAL_PREFIX
			// Prefix does not match
			return nil, errors.New("withdrawal credentials prefix does not match the spec")
		} else {
			if !strings.HasPrefix(withCred, "010000000000000000000000") {
				// Eth1 address is not left padded
				return nil, errors.New("eth1 withdrawal credentials are not left padded as the spec")
			}
		}
	}

	validDecodedWithCred, err := hex.DecodeString(withCred)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode withdrawal credentials")
	}
	return validDecodedWithCred, nil
}
