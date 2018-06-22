package claim

import (
	"github.com/dynamicgo/slf4go"
	"github.com/inwecrypto/neodb"
	"github.com/inwecrypto/neogo/rpc"
	"github.com/inwecrypto/neogo/tx"
)

// CalcBlockRange .
func CalcBlockRange(unclaimed []*rpc.UTXO) (start int64, end int64) {
	return calcBlockRange(unclaimed)
}

func calcBlockRange(unclaimed []*rpc.UTXO) (start int64, end int64) {

	hasUnspent := false

	minBlock := int64(-1)
	maxBlock := int64(0)

	for _, utxo := range unclaimed {
		if utxo.SpentBlock == -1 {
			hasUnspent = true
		}

		if minBlock == -1 {
			minBlock = utxo.Block
		}

		if minBlock > utxo.Block {
			minBlock = utxo.Block
		}

		if maxBlock < utxo.SpentBlock {
			maxBlock = utxo.SpentBlock
		}
	}

	if hasUnspent {
		end = -1
	} else {
		end = maxBlock
	}

	start = minBlock

	return
}

func blocksFee(blocks []*neodb.Block, start, end int64) (float64, int64) {
	if len(blocks) == 0 {
		return 0, 0
	}

	offset := start - blocks[0].Block

	if end == -1 {
		end = blocks[len(blocks)-1].Block
	}

	fee := float64(0)

	endoffset := end - blocks[0].Block

	if endoffset > int64(len(blocks)) {
		endoffset = int64(len(blocks))
	}

	for i := offset; i < endoffset; i++ {
		fee += blocks[i].SysFee
	}

	return fee, end
}

var log = slf4go.Get("test")

// CalcUnclaimedGas .
func CalcUnclaimedGas(unclaimed []*rpc.UTXO, blocks []*neodb.Block) (unavailable, available float64, err error) {
	return calcUnclaimedGas(unclaimed, blocks)
}

func calcUnclaimedGas(unclaimed []*rpc.UTXO, blocks []*neodb.Block) (unavailable, available float64, err error) {
	unavailableFixed8 := tx.MakeFixed8(0)
	availableFixed8 := tx.MakeFixed8(0)

	for _, utxo := range unclaimed {

		sysfee, end := blocksFee(blocks, utxo.Block, utxo.SpentBlock)

		start := utxo.Block

		gas := sysfee + getUnClaimedGas(start, end)

		val, err := utxo.Value()

		if err != nil {
			return 0, 0, err
		}

		gasFixed8 := tx.MakeFixed8(val * gas / totalNEO)

		// log.DebugF("gas %08f %d %d %f %s", gas, start, end, val, gasFixed8.String())

		if utxo.SpentBlock != -1 {
			availableFixed8 += gasFixed8
		} else {
			unavailableFixed8 += gasFixed8
		}

		utxo.Gas = gasFixed8.String()
	}

	available = availableFixed8.Float64()
	unavailable = unavailableFixed8.Float64()

	return
}
