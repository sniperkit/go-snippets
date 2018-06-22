package claim

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/dynamicgo/config"
	"github.com/go-xorm/xorm"
	"github.com/inwecrypto/neodb"
	"github.com/inwecrypto/neogo/rpc"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// var log = slf4go.Get("test")

func getUnClaimedGas2(start, end int64) float64 {
	generated := float64(0)
	for i := start; i < end; i++ {
		tmp := generateGas(i + 1)

		if tmp == 0 {
			break
		}

		generated += tmp
	}

	return generated
}

func getUnClaimedGas3(start, end int64) float64 {
	generated := float64(0)
	for i := start; i < end; i++ {
		tmp := generateGas(i)

		if tmp == 0 {
			break
		}

		generated += tmp
	}

	return generated
}

func TestGenerateCas(t *testing.T) {
	val := float64(2)
	// generated := getUnClaimedGas2(1992850, 2001227)

	// generated2 := getUnClaimedGas3(1992850, 2001227)

	// require.Equal(t, generated, generated2)

	// for i := int64(1992850); i < int64(2001227); i++ {
	// 	println(i, generateGas(i))
	// }

	// println(generateGas(1992850), generateGas(2001227))

	gas := round((val * (getUnClaimedGas(1992850, 2001227) + 2) / totalNEO), 8)

	fmt.Printf("%v\n", gas)
}

var engine *xorm.Engine

func init() {
	conf, err := config.NewFromFile("../../conf/claim.json")

	if err != nil {
		panic(err)
	}

	engine, err = createEngine(conf, "claim")

	if err != nil {
		panic(err)
	}
}

func getBlocksFee(start int64, end int64) (float64, int64, error) {

	var rows []map[string]string
	var err error

	if end == -1 {
		rows, err = engine.QueryString(`select sum(sys_fee), max(block) from neo_block where block >= ?`, start)
	} else {
		rows, err = engine.QueryString(`select sum(sys_fee), max(block) from neo_block where block >= ? and block < ?`, start, end)
	}

	if err != nil {
		return 0, end, err
	}

	if len(rows) == 0 {
		return 0, end, nil
	}

	sum, err := strconv.ParseFloat(rows[0]["sum"], 32)

	if err != nil {
		return 0, end, nil
	}

	max, err := strconv.ParseFloat(rows[0]["max"], 32)

	if err != nil {
		return 0, end, nil
	}

	return sum, int64(max), nil
}

func createEngine(conf *config.Config, name string) (*xorm.Engine, error) {
	driver := conf.GetString(fmt.Sprintf("%s.driver", name), "postgres")
	datasource := conf.GetString(fmt.Sprintf("%s.datasource", name), "")

	return xorm.NewEngine(driver, datasource)
}

func getBlocks(start int64, end int64) ([]*neodb.Block, error) {

	blocks := make([]*neodb.Block, 0)

	if end == -1 {
		err := engine.Where(`block >= ?`, start).Find(&blocks)
		if err != nil {
			return nil, err
		}
	} else {
		if err := engine.Where(`block >= ? and block <= ?`, start, end).Cols("block", "sys_fee").Find(&blocks); err != nil {
			return nil, err
		}
	}

	return blocks, nil
}

func TestVNext(t *testing.T) {
	log.Debug("start fetch unspent")
	utxos, err := unspent("AanTXadhgdHzGbmy5ZBPXxR4iPHMivzVPZ", NEOAssert)
	log.Debug("end fetch unspent")
	require.NoError(t, err)

	// println(printResult(utxos))

	start, end := calcBlockRange(utxos)

	println(start, end, len(utxos))
	log.Debug("start get blocks", start, end)
	blocks, err := getBlocks(start, end)
	log.Debug("end get blocks", len(blocks), blocks[len(blocks)-1].Block)

	require.NoError(t, err)

	u, v, err := calcUnclaimedGas(utxos, blocks)

	require.NoError(t, err)

	println(fmt.Sprintf("%.08f %.08f %.08f", u, v, u+v))

	println(utxos[0].Gas)

}

func printResult(val interface{}) string {
	data, _ := json.MarshalIndent(val, "", "\t")

	return string(data)
}

const (
	GasAssert = "0x602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
	NEOAssert = "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
)

func unspent(address string, asset string) ([]*rpc.UTXO, error) {

	tutxos := make([]*neodb.UTXO, 0)

	err :=
		engine.
			Where(`address = ? and asset = ? and claimed = FALSE`, address, asset).
			Find(&tutxos)

	if err != nil {
		return nil, err
	}

	utxos := make([]*rpc.UTXO, 0)

	for _, t := range tutxos {
		utxos = append(utxos, &rpc.UTXO{
			TransactionID: t.TX,
			Vout: rpc.Vout{
				Address: t.Address,
				Asset:   t.Asset,
				N:       t.N,
				Value:   t.Value,
			},
			CreateTime: t.CreateTime.Format(time.RFC3339Nano),
			Block:      t.CreateBlock,
			SpentBlock: t.SpentBlock,
		})
	}

	return utxos, nil
}
