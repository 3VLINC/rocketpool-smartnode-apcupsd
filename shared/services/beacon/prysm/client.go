package prysm

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "regexp"
    "strconv"
    "strings"

    "github.com/rocket-pool/smartnode/shared/services/beacon"
)


// Beacon endpoints
const REQUEST_ETH2_CONFIG_PATH string = "/eth/v1alpha1/beacon/config"
const REQUEST_BEACON_HEAD_PATH string = "/eth/v1alpha1/beacon/chainhead"
const REQUEST_VALIDATOR_PATH string = "/eth/v1alpha1/validator"


// Beacon response types
type Eth2ConfigResponse struct {
    Config struct {
        GenesisForkVersion string       `json:"GenesisForkVersion"`
        BLSWithdrawalPrefixByte string  `json:"BLSWithdrawalPrefixByte"`
        DomainBeaconProposer string     `json:"DomainBeaconProposer"`
        DomainBeaconAttester string     `json:"DomainBeaconAttester"`
        DomainRandao string             `json:"DomainRandao"`
        DomainDeposit string            `json:"DomainDeposit"`
        DomainVoluntaryExit string      `json:"DomainVoluntaryExit"`
        SlotsPerEpoch string            `json:"SlotsPerEpoch"`
    } `json:"config"`
}
type BeaconHeadResponse struct {
    HeadEpoch string                    `json:"headEpoch"`
    FinalizedEpoch string               `json:"finalizedEpoch"`
    JustifiedEpoch string               `json:"justifiedEpoch"`
}
type ValidatorResponse struct {
    PublicKey string                    `json:"publicKey"`
    WithdrawalCredentials string        `json:"withdrawalCredentials"`
    EffectiveBalance string             `json:"effectiveBalance"`
    Slashed bool                        `json:"slashed"`
    ActivationEligibilityEpoch string   `json:"activationEligibilityEpoch"`
    ActivationEpoch string              `json:"activationEpoch"`
    ExitEpoch string                    `json:"exitEpoch"`
    WithdrawableEpoch string            `json:"withdrawableEpoch"`
}


// Client
type Client struct {
    providerUrl string
}


/**
 * Create client
 */
func NewClient(providerUrl string) *Client {
    return &Client{
        providerUrl: providerUrl,
    }
}


/**
 * Get the eth2 config
 */
func (c *Client) GetEth2Config() (*beacon.Eth2Config, error) {

    // Get config
    var config Eth2ConfigResponse
    if responseBody, err := c.getRequest(REQUEST_ETH2_CONFIG_PATH); err != nil {
        return nil, errors.New("Error retrieving eth2 config: " + err.Error())
    } else if err := json.Unmarshal(responseBody, &config); err != nil {
        return nil, errors.New("Error unpacking eth2 config: " + err.Error())
    }

    // Create response
    response := &beacon.Eth2Config{}

    // Decode data and update
    if genesisForkVersion, err := deserializeBytes(config.Config.GenesisForkVersion); err != nil {
        return nil, errors.New("Error decoding genesis fork version: " + err.Error())
    } else {
        response.GenesisForkVersion = genesisForkVersion
    }
    if blsWithdrawalPrefixByteInt, err := strconv.Atoi(config.Config.BLSWithdrawalPrefixByte); err != nil {
        return nil, errors.New("Error decoding BLS withdrawal prefix byte: " + err.Error())
    } else {
        response.BLSWithdrawalPrefixByte = byte(blsWithdrawalPrefixByteInt)
    }
    if domainBeaconProposer, err := deserializeBytes(config.Config.DomainBeaconProposer); err != nil {
        return nil, errors.New("Error decoding beacon proposer domain: " + err.Error())
    } else {
        response.DomainBeaconProposer = domainBeaconProposer
    }
    if domainBeaconAttester, err := deserializeBytes(config.Config.DomainBeaconAttester); err != nil {
        return nil, errors.New("Error decoding beacon attester domain: " + err.Error())
    } else {
        response.DomainBeaconAttester = domainBeaconAttester
    }
    if domainRandao, err := deserializeBytes(config.Config.DomainRandao); err != nil {
        return nil, errors.New("Error decoding randao domain: " + err.Error())
    } else {
        response.DomainRandao = domainRandao
    }
    if domainDeposit, err := deserializeBytes(config.Config.DomainDeposit); err != nil {
        return nil, errors.New("Error decoding deposit domain: " + err.Error())
    } else {
        response.DomainDeposit = domainDeposit
    }
    if domainVoluntaryExit, err := deserializeBytes(config.Config.DomainVoluntaryExit); err != nil {
        return nil, errors.New("Error decoding voluntary exit domain: " + err.Error())
    } else {
        response.DomainVoluntaryExit = domainVoluntaryExit
    }
    if slotsPerEpoch, err := strconv.Atoi(config.Config.SlotsPerEpoch); err != nil {
        return nil, errors.New("Error decoding slots per epoch: " + err.Error())
    } else {
        response.SlotsPerEpoch = uint64(slotsPerEpoch)
    }

    // Return
    return response, nil

}


/**
 * Get the beacon head
 */
func (c *Client) GetBeaconHead() (*beacon.BeaconHead, error) {

    // Get beacon head
    var head BeaconHeadResponse
    if responseBody, err := c.getRequest(REQUEST_BEACON_HEAD_PATH); err != nil {
        return nil, errors.New("Error retrieving beacon head: " + err.Error())
    } else if err := json.Unmarshal(responseBody, &head); err != nil {
        return nil, errors.New("Error unpacking beacon head: " + err.Error())
    }

    // Create response
    response := &beacon.BeaconHead{}

    // Decode data and update
    if headEpoch, err := strconv.Atoi(head.HeadEpoch); err != nil {
        return nil, errors.New("Error decoding head epoch: " + err.Error())
    } else {
        response.Epoch = uint64(headEpoch)
    }
    if finalizedEpoch, err := strconv.Atoi(head.FinalizedEpoch); err != nil {
        return nil, errors.New("Error decoding finalized epoch: " + err.Error())
    } else {
        response.FinalizedEpoch = uint64(finalizedEpoch)
    }
    if justifiedEpoch, err := strconv.Atoi(head.JustifiedEpoch); err != nil {
        return nil, errors.New("Error decoding justified epoch: " + err.Error())
    } else {
        response.JustifiedEpoch = uint64(justifiedEpoch)
    }

    // Return
    return response, nil

}


/**
 * Get a validator's status
 */
func (c *Client) GetValidatorStatus(pubkey string) (*beacon.ValidatorStatus, error) {
    return &beacon.ValidatorStatus{}, nil
}


/**
 * Make GET request to beacon server
 */
func (c *Client) getRequest(requestPath string) ([]byte, error) {

    // Send request
    response, err := http.Get(c.providerUrl + requestPath)
    if err != nil { return nil, err }
    defer response.Body.Close()

    // Get response
    body, err := ioutil.ReadAll(response.Body)
    if err != nil { return nil, err }

    // Return
    return body, nil

}


// Deserialize a byte array
func deserializeBytes(value string) ([]byte, error) {

    // Check format
    if !regexp.MustCompile("^\\[(\\d+( \\d+)*)?\\]$").MatchString(value) {
        return nil, errors.New("Invalid byte array format")
    }

    // Get byte strings
    byteStrings := strings.Split(value[1:len(value)-1], " ")

    // Get and return bytes
    bytes := []byte{}
    for _, byteString := range byteStrings {
        if byteInt, err := strconv.Atoi(byteString); err != nil {
            return nil, errors.New("Invalid byte")
        } else {
            bytes = append(bytes, byte(byteInt))
        }
    }
    return bytes, nil

}

