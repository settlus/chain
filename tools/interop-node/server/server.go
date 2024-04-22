package server

import (
	"context"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/settlus/chain/tools/interop-node/client"
	cfg "github.com/settlus/chain/tools/interop-node/config"
	"github.com/settlus/chain/tools/interop-node/feeder"
	"github.com/settlus/chain/tools/interop-node/signer"
	"github.com/settlus/chain/tools/interop-node/subscriber"
	"github.com/settlus/chain/tools/interop-node/types"
)

const (
	IterationInterval      = 200 * time.Millisecond
	GracefulShutdownPeriod = 1500 * time.Millisecond
)

type Server struct {
	ctx    context.Context
	config *cfg.Config
	logger log.Logger

	sc *client.SettlusClient

	subscribers map[string]subscriber.Subscriber
	feeder      *feeder.Feeder
}

// NewServer creates a new interop server
func NewServer(
	config *cfg.Config,
	ctx context.Context,
	logger log.Logger,
) (*Server, error) {
	logger = logger.With("server", "interop-node")

	s := signer.NewSigner(ctx, config)
	address, err := types.GetAddressFromPubKey(s.PubKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get address from pubkey: %w", err)
	}
	config.Feeder.Address = address

	sc, err := client.NewSettlusClient(config, ctx, s, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlus client: %w", err)
	}

	subscribers, err := subscriber.InitSubscribers(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to init chain clients: %w", err)
	}
	subscribersMap := make(map[string]subscriber.Subscriber)
	for _, ss := range subscribers {
		subscribersMap[ss.Id()] = ss
	}

	feeder, err := feeder.NewFeeder(config, sc, subscribers, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create block feeder: %w", err)
	}

	return &Server{
		ctx:         ctx,
		logger:      logger,
		config:      config,
		sc:          sc,
		subscribers: subscribersMap,
		feeder:      feeder,
	}, nil
}

// Start starts the oracle feeder server
func (s *Server) Start() {
	s.logger.Info("Starting oracle feeder server")
	for _, sub := range s.subscribers {
		sub.Start(s.ctx)
	}

	s.logger.Info("Starting oracle feeder server loop")
	go s.startIteration()
}

func (s *Server) startIteration() {
	ticker := time.NewTicker(IterationInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.handleIteration()
		case <-s.ctx.Done():
			s.logger.Info("Server context canceled, stopping...")
			return
		}
	}
}

// handleIteration handles a single iteration of the oracle feeder server
func (s *Server) handleIteration() {
	latestHeight, err := s.sc.GetLatestHeight(s.ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to get latest height: %v", err))
		return
	}
	nextHeight := latestHeight + 1

	f := s.feeder
	switch {
	case f.WantAbstain(nextHeight):
		if err := f.HandleAbstain(s.ctx); err != nil {
			s.logger.Error(fmt.Sprintf("failed to handle abstain: %v", err))
		}
	case f.IsVotingPeriod(nextHeight):
		if err := f.HandleVote(s.ctx); err != nil {
			s.logger.Error(fmt.Sprintf("failed to handle vote: %v", err))
		}
	case f.IsPreVotingPeriod(nextHeight):
		if err := f.HandlePrevote(s.ctx); err != nil {
			s.logger.Error(fmt.Sprintf("failed to handle prevote: %v", err))
		}
	default:
		f.FetchNewRoundInfo(s.ctx)
	}
}

// Close gracefully closes connections
func (s *Server) Close() {
	s.logger.Info("Shutting down oracle feeder server...")

	time.Sleep(GracefulShutdownPeriod)

	s.sc.Close()
	for _, ss := range s.subscribers {
		ss.Stop()
	}
}
