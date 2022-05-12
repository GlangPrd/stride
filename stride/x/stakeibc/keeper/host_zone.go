package keeper

import (
	"encoding/binary"

	"github.com/Stride-labs/stride/x/stakeibc/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetHostZoneCount get the total number of hostZone
func (k Keeper) GetHostZoneCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.HostZoneCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetHostZoneCount set the total number of hostZone
func (k Keeper) SetHostZoneCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.HostZoneCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendHostZone appends a hostZone in the store with a new id and update the count
func (k Keeper) AppendHostZone(
	ctx sdk.Context,
	hostZone types.HostZone,
) uint64 {
	// Create the hostZone
	count := k.GetHostZoneCount(ctx)

	// Set the ID of the appended value
	hostZone.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostZoneKey))
	appendedValue := k.cdc.MustMarshal(&hostZone)
	store.Set(GetHostZoneIDBytes(hostZone.Id), appendedValue)

	// Update hostZone count
	k.SetHostZoneCount(ctx, count+1)

	return count
}

// SetHostZone set a specific hostZone in the store
func (k Keeper) SetHostZone(ctx sdk.Context, hostZone types.HostZone) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostZoneKey))
	b := k.cdc.MustMarshal(&hostZone)
	store.Set(GetHostZoneIDBytes(hostZone.Id), b)
}

// GetHostZone returns a hostZone from its id
func (k Keeper) GetHostZone(ctx sdk.Context, id uint64) (val types.HostZone, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostZoneKey))
	b := store.Get(GetHostZoneIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveHostZone removes a hostZone from the store
func (k Keeper) RemoveHostZone(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostZoneKey))
	store.Delete(GetHostZoneIDBytes(id))
}

// GetAllHostZone returns all hostZone
func (k Keeper) GetAllHostZone(ctx sdk.Context) (list []types.HostZone) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostZoneKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.HostZone
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetHostZoneIDBytes returns the byte representation of the ID
func GetHostZoneIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetHostZoneIDFromBytes returns ID in uint64 format from a byte array
func GetHostZoneIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}