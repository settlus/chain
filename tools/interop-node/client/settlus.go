package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
	cosmoscodectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cometlog "github.com/tendermint/tendermint/libs/log"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	oracletypes "github.com/settlus/chain/x/oracle/types"

	"github.com/settlus/chain/app"
	"github.com/settlus/chain/evmos/encoding"
	"github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/signer"
)

const (
	MaxConnections   = 10
	TxRetryCount     = 3
	TxRetryDelay     = time.Millisecond * 500
	HeightRetryCount = 3
	HeightRetryDelay = time.Millisecond * 50
)

var (
	HTTPProtocols = regexp.MustCompile("https?://")
)

type SettlusClient struct {
	client     *httpclient.HTTP
	grpcClient *grpc.ClientConn

	txConfig cosmosclient.TxConfig
	chainId  string
	gasLimit uint64
	fees     sdk.Coins
	feePayer string

	signer        signer.Signer
	accountnumber uint64
	sequence      uint64

	logger cometlog.Logger
}

// NewSettlusClient creates a new SettlusClient instance
func NewSettlusClient(config *config.Config, ctx context.Context, s signer.Signer, logger cometlog.Logger) (*SettlusClient, error) {
	rpcClient, gRpcClient, err := getSettlusRpcs(config.Settlus.RpcUrl, config.Settlus.GrpcUrl, config.Settlus.Insecure)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlus rpc clients: %w", err)
	}

	interfaceRegistry := getInterfaceRegistry()
	txConfig := getTxConfig(interfaceRegistry)

	account, err := getAccount(ctx, gRpcClient, config.Feeder.Address, interfaceRegistry)
	if err != nil {
		return nil, err
	}

	fees, err := getFees(config.Settlus.Fees)
	if err != nil {
		return nil, err
	}

	return &SettlusClient{
		client:        rpcClient,
		grpcClient:    gRpcClient,
		txConfig:      txConfig,
		chainId:       config.Settlus.ChainId,
		gasLimit:      config.Settlus.GasLimit,
		fees:          fees,
		feePayer:      config.Feeder.FeePayer,
		signer:        s,
		accountnumber: account.GetAccountNumber(),
		sequence:      account.GetSequence(),
		logger:        logger,
	}, nil
}

func getSettlusRpcs(rpcUrl, grpcUrl string, insecure bool) (*httpclient.HTTP, *grpc.ClientConn, error) {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(rpcUrl)
	if err != nil {
		return nil, nil, err
	}

	// Tweak the transport
	httpTransport, ok := (httpClient.Transport).(*http.Transport)
	if !ok {
		return nil, nil, fmt.Errorf("invalid HTTP Transport: %T", httpTransport)
	}
	httpTransport.MaxConnsPerHost = MaxConnections

	rpcClient, err := httpclient.NewWithClient(rpcUrl, "/websocket", httpClient)
	if err != nil {
		return nil, nil, err
	}

	if err = rpcClient.Start(); err != nil {
		return nil, nil, err
	}

	// Get grpc client
	gc, err := CreateGrpcConnection(insecure, grpcUrl)
	if err != nil {
		return nil, nil, err
	}

	return rpcClient, gc, nil
}

// getFees parses the fees from the given config
func getFees(fee config.Fee) (sdk.Coins, error) {
	amount, ok := sdk.NewIntFromString(fee.Amount)
	if !ok {
		return nil, fmt.Errorf("invalid amount %s", fee)
	}

	return sdk.NewCoins(sdk.NewCoin(fee.Denom, amount)), nil
}

// getInterfaceRegistry creates a new interface registry
func getInterfaceRegistry() cosmoscodectypes.InterfaceRegistry {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	return encodingConfig.InterfaceRegistry
}

// getTxConfig creates a new transaction config
func getTxConfig(ir cosmoscodectypes.InterfaceRegistry) cosmosclient.TxConfig {
	marshaler := cosmoscodec.NewProtoCodec(ir)
	return authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)
}

// getAccount gets the account from the given address
func getAccount(ctx context.Context, gRpcClient *grpc.ClientConn, feederAddress string, ir cosmoscodectypes.InterfaceRegistry) (authtypes.AccountI, error) {
	qc := authtypes.NewQueryClient(gRpcClient)
	res, err := qc.Account(ctx, &authtypes.QueryAccountRequest{Address: feederAddress})
	if err != nil {
		return nil, err
	}

	var account authtypes.AccountI
	if err := ir.UnpackAny(res.Account, &account); err != nil {
		return nil, err
	}
	return account, nil
}

// Close closes the SettlusClient
func (sc *SettlusClient) Close() {
	err := sc.client.Stop()
	if err != nil {
		return
	}
}

// CreateGrpcConnection creates a new gRPC client connection from the given configuration
func CreateGrpcConnection(isInsecure bool, grpcAddress string) (*grpc.ClientConn, error) {
	var grpcOpts []grpc.DialOption
	if isInsecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})))
	}

	address := HTTPProtocols.ReplaceAllString(grpcAddress, "")
	return grpc.Dial(address, grpcOpts...)
}

// LatestHeight get the latest height from the RPC client.
func (sc *SettlusClient) latestHeight(ctx context.Context) (int64, error) {
	status, err := sc.client.Status(ctx)
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// GetLatestHeight tries getting the latest height from the RPC client.
func (sc *SettlusClient) GetLatestHeight(ctx context.Context) (int64, error) {
	for retryCount := 0; retryCount < HeightRetryCount; retryCount++ {
		latestBlockHeight, err := sc.latestHeight(ctx)
		if err == nil {
			sc.logger.Debug("GetLatestHeight", "blocknumber", latestBlockHeight)
			return latestBlockHeight, nil
		}

		time.Sleep(HeightRetryDelay)
	}

	return -1, nil
}

// buildTx builds a cosmos transaction
func (sc *SettlusClient) buildTx(msg sdk.Msg) ([]byte, error) {
	txBuilder := sc.txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, fmt.Errorf("failed to set message: %w", err)
	}
	txBuilder.SetGasLimit(sc.gasLimit)
	txBuilder.SetFeeAmount(sc.fees)

	if sc.feePayer != "" {
		addr, err := sdk.AccAddressFromBech32(sc.feePayer)
		if err != nil {
			return nil, fmt.Errorf("failed to parse fee payer address: %w", err)
		}
		txBuilder.SetFeePayer(addr)
	}

	pubKey := sc.signer.PubKey()

	sigData := &signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     sigData,
		Sequence: sc.sequence,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, fmt.Errorf("failed to set signatures: %w", err)
	}

	signerData := authsigning.SignerData{
		ChainID:       sc.chainId,
		AccountNumber: sc.accountnumber,
		Sequence:      sc.sequence,
		PubKey:        pubKey,
		Address:       sdk.AccAddress(pubKey.Address()).String(),
	}

	signBytes, err := sc.txConfig.SignModeHandler().GetSignBytes(signing.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())
	if err != nil {
		return nil, fmt.Errorf("failed to get sign bytes: %w", err)
	}

	signature, err := sc.signer.Sign(signBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	sigData.Signature = signature

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, fmt.Errorf("failed to set signatures: %w", err)
	}

	return sc.txConfig.TxEncoder()(txBuilder.GetTx())
}

// sendTx sends a transaction to the Settlus node
func (sc *SettlusClient) sendTx(ctx context.Context, tx []byte) (*coretypes.ResultBroadcastTxCommit, error) {
	res, err := sc.client.BroadcastTxCommit(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast tx: %w", err)
	}

	if res.CheckTx.Code != 0 {
		if res.CheckTx.Code == errors.ErrWrongSequence.ABCICode() {
			return res, fmt.Errorf("failed to broadcast tx, sequence number mismatch, %s", res.CheckTx.Log)
		}
		return nil, fmt.Errorf("failed to broadcast tx, check tx failed, code: %d, log: %s", res.CheckTx.Code, res.CheckTx.Log)
	}

	if res.DeliverTx.Code != 0 {
		if res.DeliverTx.Code == errors.ErrWrongSequence.ABCICode() {
			return res, fmt.Errorf("failed to broadcast tx, sequence number mismatch, %s", res.DeliverTx.Log)
		}
		return nil, fmt.Errorf("failed to broadcast tx, deliver tx failed, code: %d, log: %s", res.DeliverTx.Code, res.DeliverTx.Log)
	}

	// If the tx was successful, increment the sequence
	sc.sequence++

	sc.logger.Debug("tx sent successfully: ", res.Hash.String())

	return nil, nil
}

// BuildAndSendTxWithRetry builds and sends a transaction to the Settlus node
func (sc *SettlusClient) BuildAndSendTxWithRetry(ctx context.Context, msg sdk.Msg) error {
	for attempt := 0; attempt < TxRetryCount; attempt++ {
		tx, err := sc.buildTx(msg)
		if err != nil {
			return fmt.Errorf("failed to build tx: %v", err)
		}

		res, err := sc.sendTx(ctx, tx)
		if err == nil {
			return nil
		}

		if res != nil && res.CheckTx.Code == errors.ErrWrongSequence.ABCICode() {
			if err := sc.UpdateSequenceFromError(res.CheckTx.Log); err != nil {
				return fmt.Errorf("failed to update sequence number: %w", err)
			}
			// retry immediately if sequence number mismatch
			continue
		}

		if attempt < TxRetryCount-1 {
			sc.logger.Debug(fmt.Sprintf("failed to send tx, %v,retrying in %s", err, TxRetryDelay))
			time.Sleep(TxRetryDelay)
		}
	}

	return fmt.Errorf("failed to send tx after %d attempts", TxRetryCount)
}

// UpdateSequenceFromError updates the sequence number from the error log
func (sc *SettlusClient) UpdateSequenceFromError(log string) error {
	chunk := strings.Split(log, "expected ")[1]
	sequence := strings.Split(chunk, ",")[0]
	sequenceInt, err := strconv.ParseUint(sequence, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to broadcast tx because of wrong sequence, and failed to parse sequence number from the error log, %w", err)
	}
	sc.sequence = sequenceInt
	sc.logger.Debug(fmt.Sprintf("sequence number mismatch, updated sequence number to %d", sequenceInt))
	return nil
}

func (sc *SettlusClient) FetchNewRoundInfo(ctx context.Context) (*oracletypes.RoundInfo, error) {
	qc := oracletypes.NewQueryClient(sc.grpcClient)
	res, err := qc.CurrentRoundInfo(ctx, &oracletypes.QueryCurrentRoundInfoRequest{})
	if err != nil {
		return nil, err
	}

	return res.RoundInfo, nil
}
