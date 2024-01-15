package keeper

import (
	"context"
	"fmt"

	"github.com/settlus/chain/x/settlement/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = SettlementKeeper{}

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	*SettlementKeeper
}

// Params implements the Query/Params gRPC method
func (k SettlementKeeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// UTXR implements the Query/UTXR gRPC method
func (k SettlementKeeper) UTXR(c context.Context, req *types.QueryUTXRRRequest) (*types.QueryUTXRResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	utxr := k.GetUTXRByRequestId(
		ctx,
		req.TenantId,
		req.RequestId,
	)
	if utxr == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("utxr with [request_id: %s] not found", req.RequestId))
	}

	return &types.QueryUTXRResponse{Utxr: *utxr}, nil
}

// UTXRs implements the Query/UTXRs gRPC method
func (k SettlementKeeper) UTXRs(c context.Context, req *types.QueryUTXRsRequest) (*types.QueryUTXRsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var utxrs []types.UTXR
	ctx := sdk.UnwrapSDKContext(c)
	if !k.CheckTenantExist(ctx, req.TenantId) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("tenant with [tenant_id: %d] not found", req.TenantId))
	}

	utxrStore := k.GetUTXRStore(ctx, req.TenantId)

	pageRes, err := query.Paginate(utxrStore, req.Pagination, func(_ []byte, value []byte) error {
		var utxr types.UTXR
		k.cdc.MustUnmarshal(value, &utxr)

		utxrs = append(utxrs, utxr)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUTXRsResponse{Utxrs: utxrs, Pagination: pageRes}, nil
}

// Tenant implements the Query/Tenant gRPC method
func (k SettlementKeeper) Tenant(c context.Context, req *types.QueryTenantRequest) (*types.QueryTenantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	tenant := k.GetTenant(ctx, req.TenantId)
	if tenant == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("tenant with [tenant_id: %d] not found", req.TenantId))
	}

	buildTenantWithTreasury := k.buildTenantWithTreasury(ctx, tenant)
	return &types.QueryTenantResponse{Tenant: buildTenantWithTreasury}, nil
}

// Tenants implements the Query/Tenants gRPC method
func (k SettlementKeeper) Tenants(c context.Context, req *types.QueryTenantsRequest) (*types.QueryTenantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var result []types.TenantWithTreasury
	ctx := sdk.UnwrapSDKContext(c)

	tenantStore := k.GetTenantStore(ctx)

	pageRes, err := query.Paginate(tenantStore, req.Pagination, func(_ []byte, value []byte) error {
		var tenant types.Tenant
		k.cdc.MustUnmarshal(value, &tenant)
		tenantWithTreasury := k.buildTenantWithTreasury(ctx, &tenant)
		result = append(result, tenantWithTreasury)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTenantsResponse{Tenants: result, Pagination: pageRes}, nil
}

func (k SettlementKeeper) buildTenantWithTreasury(ctx sdk.Context, tenant *types.Tenant) types.TenantWithTreasury {
	treasuryAccount, treasuryBalance := k.GetTenantTreasury(ctx, tenant.Id)
	tenantWithTreasury := types.TenantWithTreasury{
		Tenant: tenant,
		Treasury: &types.Treasury{
			Address: treasuryAccount.String(),
		},
	}

	if tenant.PayoutMethod == types.PayoutMethod_Native {
		balance := sdk.NewCoin(tenant.Denom, treasuryBalance.AmountOf(tenant.Denom))
		tenantWithTreasury.Treasury.Balance = &balance
	}

	return tenantWithTreasury
}
