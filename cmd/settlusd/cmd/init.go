package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cometbft/cometbft/libs/cli"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/types"

	"github.com/cosmos/go-bip39"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

const (
	FlagPersistentPeers      = "persistent-peers"
	FlagEnableTelemetry      = "enable-telemetry"
	FlagStateSyncRpcServers  = "state-sync.rpc-servers"
	FlagStateSyncTrustHeight = "state-sync.trust-height"
	FlagStateSyncTrustHash   = "state-sync.trust-hash"
	FlagPrivValidatorListen  = "priv-validator-listen"
)

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out))); err != nil {
		return err
	}

	return nil
}

func initializeStateSync(config *cfg.Config, cmd *cobra.Command) error {
	stateSyncRpcServers, _ := cmd.Flags().GetString(FlagStateSyncRpcServers)
	if stateSyncRpcServers == "" {
		return nil
	}

	config.StateSync.RPCServers = strings.Split(stateSyncRpcServers, ",")

	stateSyncTrustHeight, _ := cmd.Flags().GetInt64(FlagStateSyncTrustHeight)
	if stateSyncTrustHeight <= 0 {
		return fmt.Errorf("invalid trust height: %d", stateSyncTrustHeight)
	}

	config.StateSync.TrustHeight = stateSyncTrustHeight

	stateSyncTrustHash, _ := cmd.Flags().GetString(FlagStateSyncTrustHash)
	if stateSyncTrustHash == "" {
		return fmt.Errorf("trust hash cannot be empty")
	}

	config.StateSync.TrustHash = stateSyncTrustHash

	config.StateSync.Enable = true
	return nil
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init MONIKER",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			// Set peers in and out to an 8:1 ratio to prevent choking
			config.P2P.MaxNumInboundPeers = 240
			config.P2P.MaxNumOutboundPeers = 30
			config.P2P.ListenAddress = "tcp://0.0.0.0:26656"

			persistentPeers, _ := cmd.Flags().GetString(FlagPersistentPeers)
			if persistentPeers != "" {
				config.P2P.PersistentPeers = persistentPeers
			}

			config.RPC.ListenAddress = "tcp://0.0.0.0:26657"

			config.Mempool.Size = 10000

			if err := initializeStateSync(config, cmd); err != nil {
				return err
			}
			config.StateSync.TrustPeriod = 336 * time.Hour

			privValidatorListen, _ := cmd.Flags().GetBool(FlagPrivValidatorListen)
			if privValidatorListen {
				config.PrivValidatorListenAddr = "tcp://127.0.0.1:26658"
			}

			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("settlus_1001-%v", tmrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string

			recoverKey, _ := cmd.Flags().GetBool(genutilcli.FlagRecover)
			if recoverKey {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			nodeID, _, err := initializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(genutilcli.FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				fmt.Println("genesis.json file already exists: ", genFile)
				return nil
			}

			appState, err := json.MarshalIndent(mbm.DefaultGenesis(cdc), "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return errors.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			genDoc.ConsensusParams = types.DefaultConsensusParams()
			genDoc.ConsensusParams.Validator.PubKeyTypes = []string{types.ABCIPubKeyTypeSecp256k1}

			if err := genutil.ExportGenesisFile(genDoc, genFile); err != nil {
				return errors.Wrap(err, "Failed to export genesis file")
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(genutilcli.FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(genutilcli.FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagPersistentPeers, "", "persistent peer addresses separated by comma")
	cmd.Flags().Bool(FlagEnableTelemetry, false, "enable telemetry server")
	cmd.Flags().String(FlagStateSyncRpcServers, "", "state sync rpc servers separated by comma")
	cmd.Flags().Int64(FlagStateSyncTrustHeight, 0, "state sync trust height")
	cmd.Flags().String(FlagStateSyncTrustHash, "", "state sync trust hash")
	cmd.Flags().Bool(FlagPrivValidatorListen, false, "enable priv validator listen")

	return cmd
}

func initializeNodeValidatorFilesFromMnemonic(config *cfg.Config, mnemonic string) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, fmt.Errorf("invalid mnemonic")
	}
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, err
	}

	nodeID = string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0o777); err != nil {
		return "", nil, err
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0o777); err != nil {
		return "", nil, err
	}

	var filePV *privval.FilePV
	if len(mnemonic) == 0 {
		if tmos.FileExists(pvKeyFile) {
			filePV = privval.LoadFilePV(pvKeyFile, pvStateFile)
		} else {
			filePV = privval.NewFilePV(secp256k1.GenPrivKey(), pvKeyFile, pvStateFile)
			filePV.Save()
		}
	} else {
		privKey := secp256k1.GenPrivKeySecp256k1([]byte(mnemonic))
		filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
		filePV.Save()
	}

	tmValPubKey, err := filePV.GetPubKey()
	if err != nil {
		return "", nil, err
	}

	valPubKey, err = cryptocodec.FromTmPubKeyInterface(tmValPubKey)
	if err != nil {
		return "", nil, err
	}

	return nodeID, valPubKey, nil
}
