// Copyright 2018 The dexon-consensus-core Authors
// This file is part of the dexon-consensus-core library.
//
// The dexon-consensus-core library is free software: you can redistribute it
// and/or modify it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The dexon-consensus-core library is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dexon-consensus-core library. If not, see
// <http://www.gnu.org/licenses/>.

package simulation

import (
	"fmt"
	"sync"
	"time"

	"github.com/dexon-foundation/dexon-consensus-core/core/types"
	"github.com/dexon-foundation/dexon-consensus-core/simulation/config"
)

// simGovernance is a simulated governance contract implementing the
// core.Governance interface.
type simGovernance struct {
	lock                  sync.RWMutex
	notarySet             map[types.ValidatorID]struct{}
	expectedNumValidators int
	k                     int
	phiRatio              float32
	chainNum              uint32
	crs                   string
	dkgComplaint          map[uint64][]*types.DKGComplaint
	dkgMasterPublicKey    map[uint64][]*types.DKGMasterPublicKey
	lambda                time.Duration
}

// newSimGovernance returns a new simGovernance instance.
func newSimGovernance(
	numValidators int, consensusConfig config.Consensus) *simGovernance {
	return &simGovernance{
		notarySet:             make(map[types.ValidatorID]struct{}),
		expectedNumValidators: numValidators,
		k:                  consensusConfig.K,
		phiRatio:           consensusConfig.PhiRatio,
		chainNum:           consensusConfig.ChainNum,
		crs:                consensusConfig.GenesisCRS,
		dkgComplaint:       make(map[uint64][]*types.DKGComplaint),
		dkgMasterPublicKey: make(map[uint64][]*types.DKGMasterPublicKey),
		lambda:             time.Duration(consensusConfig.Lambda) * time.Millisecond,
	}
}

// GetNotarySet returns the current notary set.
func (g *simGovernance) GetNotarySet() map[types.ValidatorID]struct{} {

	g.lock.RLock()
	defer g.lock.RUnlock()

	// Return the cloned notarySet.
	ret := make(map[types.ValidatorID]struct{})
	for k := range g.notarySet {
		ret[k] = struct{}{}
	}
	return ret
}

// GetConfiguration returns the configuration at a given block height.
func (g *simGovernance) GetConfiguration(blockHeight uint64) *types.Config {
	return &types.Config{
		NumShards:  1,
		NumChains:  g.chainNum,
		GenesisCRS: g.crs,
		Lambda:     g.lambda,
		K:          g.k,
		PhiRatio:   g.phiRatio,
	}
}

// addValidator add a new validator into the simulated governance contract.
func (g *simGovernance) addValidator(vID types.ValidatorID) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, exists := g.notarySet[vID]; exists {
		return
	}
	if len(g.notarySet) == g.expectedNumValidators {
		panic(fmt.Errorf("attempt to add validator when ready"))
	}
	g.notarySet[vID] = struct{}{}
}

// AddDKGComplaint adds a DKGComplaint.
func (g *simGovernance) AddDKGComplaint(complaint *types.DKGComplaint) {
	g.dkgComplaint[complaint.Round] = append(
		g.dkgComplaint[complaint.Round], complaint)
}

// DKGComplaints returns the DKGComplaints of round.
func (g *simGovernance) DKGComplaints(round uint64) []*types.DKGComplaint {
	complaints, exist := g.dkgComplaint[round]
	if !exist {
		return []*types.DKGComplaint{}
	}
	return complaints
}

// AddDKGMasterPublicKey adds a DKGMasterPublicKey.
func (g *simGovernance) AddDKGMasterPublicKey(
	masterPublicKey *types.DKGMasterPublicKey) {
	g.dkgMasterPublicKey[masterPublicKey.Round] = append(
		g.dkgMasterPublicKey[masterPublicKey.Round], masterPublicKey)
}

// DKGMasterPublicKeys returns the DKGMasterPublicKeys of round.
func (g *simGovernance) DKGMasterPublicKeys(
	round uint64) []*types.DKGMasterPublicKey {
	masterPublicKeys, exist := g.dkgMasterPublicKey[round]
	if !exist {
		return []*types.DKGMasterPublicKey{}
	}
	return masterPublicKeys
}
