package config

import (
	"runtime"

	"github.com/rocket-pool/smartnode/shared/types/config"
)

const (
	// Prater
	nimbusBnTagTest string = "statusim/nimbus-eth2:multiarch-v23.1.0"
	nimbusVcTagTest string = "statusim/nimbus-validator-client:multiarch-v23.1.0"

	// Mainnet
	nimbusBnTagProd string = "statusim/nimbus-eth2:multiarch-v23.1.0"
	nimbusVcTagProd string = "statusim/nimbus-validator-client:multiarch-v23.1.0"

	defaultNimbusMaxPeersArm uint16 = 100
	defaultNimbusMaxPeersAmd uint16 = 160
)

// Configuration for Nimbus
type NimbusConfig struct {
	Title string `yaml:"-"`

	// The max number of P2P peers to connect to
	MaxPeers config.Parameter `yaml:"maxPeers,omitempty"`

	// Common parameters that Nimbus doesn't support and should be hidden
	UnsupportedCommonParams []string `yaml:"-"`

	// The Docker Hub tag for the BN
	BnContainerTag config.Parameter `yaml:"bnContainerTag,omitempty"`

	// The Docker Hub tag for the VC
	VcContainerTag config.Parameter `yaml:"vcContainerTag,omitempty"`

	// Custom command line flags for the BN
	AdditionalBnFlags config.Parameter `yaml:"additionalBnFlags,omitempty"`

	// Custom command line flags for the VC
	AdditionalVcFlags config.Parameter `yaml:"additionalVcFlags,omitempty"`
}

// Generates a new Nimbus configuration
func NewNimbusConfig(cfg *RocketPoolConfig) *NimbusConfig {
	return &NimbusConfig{
		Title: "Nimbus Settings",

		MaxPeers: config.Parameter{
			ID:                   "maxPeers",
			Name:                 "Max Peers",
			Description:          "The maximum number of peers your client should try to maintain. You can try lowering this if you have a low-resource system or a constrained network.",
			Type:                 config.ParameterType_Uint16,
			Default:              map[config.Network]interface{}{config.Network_All: getNimbusDefaultPeers()},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth2},
			EnvironmentVariables: []string{"BN_MAX_PEERS"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		BnContainerTag: config.Parameter{
			ID:          "bnContainerTag",
			Name:        "Beacon Node Container Tag",
			Description: "The tag name of the Nimbus Beacon Node container you want to use on Docker Hub.",
			Type:        config.ParameterType_String,
			Default: map[config.Network]interface{}{
				config.Network_Mainnet: nimbusBnTagProd,
				config.Network_Prater:  nimbusBnTagTest,
				config.Network_Devnet:  nimbusBnTagTest,
			},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth2},
			EnvironmentVariables: []string{"BN_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		VcContainerTag: config.Parameter{
			ID:          "containerTag",
			Name:        "Validator Client Container Tag",
			Description: "The tag name of the Nimbus Validator Client container you want to use on Docker Hub.",
			Type:        config.ParameterType_String,
			Default: map[config.Network]interface{}{
				config.Network_Mainnet: nimbusVcTagProd,
				config.Network_Prater:  nimbusVcTagTest,
				config.Network_Devnet:  nimbusVcTagTest,
			},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Validator},
			EnvironmentVariables: []string{"VC_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		AdditionalBnFlags: config.Parameter{
			ID:                   "additionalBnFlags",
			Name:                 "Additional Beacon Client Flags",
			Description:          "Additional custom command line flags you want to pass Nimbus's Beacon Client, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth2},
			EnvironmentVariables: []string{"BN_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},

		AdditionalVcFlags: config.Parameter{
			ID:                   "additionalVcFlags",
			Name:                 "Additional Validator Client Flags",
			Description:          "Additional custom command line flags you want to pass Nimbus's Validator Client, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Validator},
			EnvironmentVariables: []string{"VC_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the parameters for this config
func (cfg *NimbusConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.MaxPeers,
		&cfg.BnContainerTag,
		&cfg.VcContainerTag,
		&cfg.AdditionalBnFlags,
		&cfg.AdditionalVcFlags,
	}
}

// Get the common params that this client doesn't support
func (cfg *NimbusConfig) GetUnsupportedCommonParams() []string {
	return cfg.UnsupportedCommonParams
}

// Get the Docker container name of the validator client
func (cfg *NimbusConfig) GetValidatorImage() string {
	return cfg.VcContainerTag.Value.(string)
}

// Get the name of the client
func (cfg *NimbusConfig) GetName() string {
	return "Nimbus"
}

// The the title for the config
func (cfg *NimbusConfig) GetConfigTitle() string {
	return cfg.Title
}

func getNimbusDefaultPeers() uint16 {
	if runtime.GOARCH == "arm64" {
		return defaultNimbusMaxPeersArm
	}

	return defaultNimbusMaxPeersAmd
}
