// Copyright © 2021 Weald Technology Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgresql_test

import (
	"context"
	"os"
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/chaind/services/chaindb"
	"github.com/wealdtech/chaind/services/chaindb/postgresql"
)

func TestSetDeposits(t *testing.T) {
	ctx := context.Background()
	s, err := postgresql.New(ctx,
		postgresql.WithLogLevel(zerolog.Disabled),
		postgresql.WithConnectionURL(os.Getenv("CHAINDB_DATABASE_URL")),
	)
	require.NoError(t, err)

	// Fetch a block so we can set the deposits' block root.
	blocks, err := s.BlocksBySlot(ctx, 0)
	require.NoError(t, err)
	require.Len(t, blocks, 1)

	deposit := &chaindb.Deposit{
		InclusionSlot:      0,
		InclusionBlockRoot: blocks[0].Root,
		InclusionIndex:     0,
		ValidatorPubKey: spec.BLSPubKey{
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		},
		WithdrawalCredentials: []byte{0x0c, 0x0d, 0x0e, 0x0f},
		Amount:                32000000000,
	}
	otherDeposit := &chaindb.Deposit{
		InclusionSlot:      0,
		InclusionBlockRoot: blocks[0].Root,
		InclusionIndex:     1,
		ValidatorPubKey: spec.BLSPubKey{
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		},
		WithdrawalCredentials: []byte{0x2c, 0x2d, 0x2e, 0x2f},
		Amount:                32000000000,
	}

	// Try to set outside of a transaction; should fail.
	require.EqualError(t, s.SetDeposit(ctx, deposit), postgresql.ErrNoTransaction.Error())

	ctx, cancel, err := s.BeginTx(ctx)
	require.NoError(t, err)
	defer cancel()

	// Set the deposits.
	require.NoError(t, s.SetDeposit(ctx, deposit))
	require.NoError(t, s.SetDeposit(ctx, otherDeposit))

	// Attempt to set the same deposit again; should succeed.
	require.NoError(t, s.SetDeposit(ctx, deposit))
}

func TestDeposits(t *testing.T) {
	ctx := context.Background()
	s, err := postgresql.New(ctx,
		postgresql.WithLogLevel(zerolog.Disabled),
		postgresql.WithConnectionURL(os.Getenv("CHAINDB_DATABASE_URL")),
	)
	require.NoError(t, err)

	// Fetch a block so we can set the deposits' block root.
	blocks, err := s.BlocksBySlot(ctx, 0)
	require.NoError(t, err)
	require.Len(t, blocks, 1)

	deposit := &chaindb.Deposit{
		InclusionSlot:      0,
		InclusionBlockRoot: blocks[0].Root,
		InclusionIndex:     0,
		ValidatorPubKey: spec.BLSPubKey{
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		},
		WithdrawalCredentials: []byte{0x0c, 0x0d, 0x0e, 0x0f},
		Amount:                32000000000,
	}
	otherDeposit := &chaindb.Deposit{
		InclusionSlot:      0,
		InclusionBlockRoot: blocks[0].Root,
		InclusionIndex:     1,
		ValidatorPubKey: spec.BLSPubKey{
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		},
		WithdrawalCredentials: []byte{0x2c, 0x2d, 0x2e, 0x2f},
		Amount:                32000000000,
	}

	ctx, cancel, err := s.BeginTx(ctx)
	require.NoError(t, err)
	defer cancel()

	// Set the deposits.
	require.NoError(t, s.SetDeposit(ctx, deposit))
	require.NoError(t, s.SetDeposit(ctx, otherDeposit))

	// Fetch the deposits individually.
	deposits, err := s.DepositsByPublicKey(ctx, []spec.BLSPubKey{
		{
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		},
	})
	require.NoError(t, err)
	require.Len(t, deposits, 1)

	deposits, err = s.DepositsByPublicKey(ctx, []spec.BLSPubKey{
		{
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		},
	})
	require.NoError(t, err)
	require.Len(t, deposits, 1)

	// Fetch the deposits together.
	deposits, err = s.DepositsByPublicKey(ctx, []spec.BLSPubKey{
		{
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
			0x80, 0x81, 0x82, 0x83, 0x84, 0x84, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		},
		{
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
			0x90, 0x91, 0x92, 0x93, 0x94, 0x94, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		},
	})
	require.NoError(t, err)
	require.Len(t, deposits, 2)
}