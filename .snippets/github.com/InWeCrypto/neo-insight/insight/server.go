package insight

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"github.com/inwecrypto/neo-insight/claim"
	"github.com/inwecrypto/neodb"
	"github.com/inwecrypto/neogo/rpc"
	"github.com/julienschmidt/httprouter"
	"github.com/ybbus/jsonrpc"
)

var logger slf4go.Logger

// OpenLogger .
func OpenLogger() {
	logger = slf4go.Get("neo-insight")
}

type handler func(params []interface{}) (interface{}, *JSONRPCError)

type syncAddress struct {
	Address string
	Times   int
}

func (address *syncAddress) String() string {
	return fmt.Sprintf("%s (%d)", address.Address, address.Times)
}

// Server insight api jsonrpc 2.0 server
type Server struct {
	mutex        sync.Mutex
	cnf          *config.Config
	router       *httprouter.Router
	remote       *url.URL
	dispatch     map[string]handler
	engine       *xorm.Engine
	redisclient  *redis.Client
	syncChan     chan *syncAddress
	syncFlag     map[string]string
	syncTimes    int
	syncDuration time.Duration
}

type loggerHandler struct {
	handler http.Handler
}

func (l *loggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.DebugF("http route: %s %s", r.Method, r.URL)
	l.handler.ServeHTTP(w, r)
}

// NewServer create new server
func NewServer(cnf *config.Config) (*Server, error) {

	remote, err := url.Parse(cnf.GetString("insight.neo", "http://xxxxxx:10332"))

	if err != nil {
		return nil, err
	}

	username := cnf.GetString("insight.neodb.username", "xxx")
	password := cnf.GetString("insight.neodb.password", "xxx")
	port := cnf.GetString("insight.neodb.port", "6543")
	host := cnf.GetString("insight.neodb.host", "localhost")
	scheme := cnf.GetString("insight.neodb.schema", "postgres")

	engine, err := xorm.NewEngine(
		"postgres",
		fmt.Sprintf(
			"user=%v password=%v host=%v dbname=%v port=%v sslmode=disable",
			username, password, host, scheme, port,
		),
	)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cnf.GetString("insight.redis.address", "localhost:6379"),
		Password: cnf.GetString("insight.redis.password", "xxxxxx"), // no password set
		DB:       int(cnf.GetInt64("insight.redis.db", 1)),          // use default DB
	})

	server := &Server{
		cnf:          cnf,
		router:       httprouter.New(),
		remote:       remote,
		dispatch:     make(map[string]handler),
		engine:       engine,
		redisclient:  client,
		syncChan:     make(chan *syncAddress, cnf.GetInt64("insight.sync_chan_length", 1024)),
		syncFlag:     make(map[string]string),
		syncTimes:    int(cnf.GetInt64("insight.sync_times", 20)),
		syncDuration: time.Second * cnf.GetDuration("insight.sync_duration", 4),
	}

	return server, nil
}

func openDB(cnf *config.Config) (*sql.DB, error) {
	driver := cnf.GetString("insight.database.driver", "xxxx")
	username := cnf.GetString("insight.database.username", "xxx")
	password := cnf.GetString("insight.database.password", "xxx")
	port := cnf.GetString("insight.database.port", "6543")
	host := cnf.GetString("insight.database.host", "localhost")
	schema := cnf.GetString("insight.database.schema", "postgres")
	maxconn := cnf.GetInt64("insight.database.maxconn", 10)

	db, err := sql.Open(driver, fmt.Sprintf("user=%v password=%v host=%v dbname=%v port=%v sslmode=disable", username, password, host, schema, port))

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(maxconn))

	return db, nil
}

func makeJSONRPCError(w http.ResponseWriter, id uint, code int, message string, data interface{}) {
	response := &jsonrpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &jsonrpc.RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	jsonresponse, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "server internal error", http.StatusInternalServerError)
		logger.ErrorF("marshal response error :%s", err)
		return
	}

	w.WriteHeader(200)

	if _, err := w.Write(jsonresponse); err != nil {
		logger.ErrorF("write response error :%s", err)
	}
}

func makeRPCRequest(r *http.Request) (*jsonrpc.RPCRequest, error) {
	request := jsonrpc.RPCRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	err := decoder.Decode(&request)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func makeJSONRPCResponse(w http.ResponseWriter, id uint, data interface{}) {
	response := &jsonrpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  data,
	}

	jsonresponse, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "server internal error", http.StatusInternalServerError)
		logger.ErrorF("marshal response error :%s", err)
		return
	}

	w.WriteHeader(200)

	if _, err := w.Write(jsonresponse); err != nil {
		logger.ErrorF("write response error :%s", err)
	}
}

// Run insight server
func (server *Server) Run() {

	server.router.POST(server.cnf.GetString("insight.extend", "/extend"), server.DipsatchJSONRPC)

	server.router.POST(server.cnf.GetString("insight.proxy", "/"), server.ReverseProxy)

	server.dispatch["balance"] = server.getBalance
	server.dispatch["claim"] = server.getClaim

	go server.syncCached()

	logger.Fatal(http.ListenAndServe(
		server.cnf.GetString("insight.listen", ":10332"),
		&loggerHandler{
			handler: server.router,
		},
	))
}

func (server *Server) syncCached() {

	for address := range server.syncChan {

		logger.DebugF("sync address claimed utxos %s", address)

		startTime := time.Now()

		unclaimed, err := server.doGetClaim(address.Address)

		claimTimes := time.Now().Sub(startTime)

		if err != nil {
			logger.ErrorF("sync claim for address %s err, %s", address, err)
			continue
		}

		logger.DebugF("[doGetClaim] claim %s spent times %s", address.Address, claimTimes)

		data, err := json.Marshal(unclaimed)

		if err != nil {
			logger.ErrorF("sync claim for address %s err, %s", address, err)
			continue
		}

		err = server.redisclient.Set(address.Address, data, time.Hour*24).Err()

		if err != nil {
			logger.ErrorF("cached claim for address %s err, %s", address, err)
			continue
		}

		logger.DebugF(" sync address claimed utxos %s -- success", address)

		address.Times--

		if address.Times > 0 {

			requeue := address

			syncDuration := server.syncDuration

			if claimTimes > 20*time.Second {
				syncDuration = time.Minute * 10
			}

			if syncDuration < server.syncDuration {
				syncDuration = server.syncDuration
			}

			time.AfterFunc(syncDuration, func() {
				logger.DebugF("requeue sync address %s", requeue)

				server.syncChan <- requeue

				logger.DebugF("requeue sync address %s -- success", requeue)
			})

		} else {
			logger.DebugF("delete sync address %s", address)
			server.removeAddress(address.Address)
			logger.DebugF("delete sync address %s -- success", address)
		}
	}

}

func (server *Server) removeAddress(address string) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	delete(server.syncFlag, address)

}

// ReverseProxy reverse proxy handler
func (server *Server) ReverseProxy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	reverseProxy := httputil.NewSingleHostReverseProxy(server.remote)

	reverseProxy.ServeHTTP(w, r)
}

// DipsatchJSONRPC .
func (server *Server) DipsatchJSONRPC(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	logger.DebugF("call extend api :%s", r.RemoteAddr)

	request, err := makeRPCRequest(r)

	if err != nil {
		makeJSONRPCError(w, 0, JSONRPCParserError, "parse error", nil)
		return
	}

	if method, ok := server.dispatch[request.Method]; ok {
		result, err := method(request.Params.([]interface{}))

		if err != nil {
			makeJSONRPCError(w, request.ID, err.ID, err.Message, result)
		} else {
			makeJSONRPCResponse(w, request.ID, result)
		}

	} else {
		makeJSONRPCError(w, request.ID, JSONRPCMethodNotFound, fmt.Sprintf("method %s not found", request.Method), nil)
	}
}

func (server *Server) getBalance(params []interface{}) (interface{}, *JSONRPCError) {
	if len(params) < 2 {
		return nil, errorf(JSONRPCInvalidParams, "expect address and asset parameters")
	}

	address, ok := params[0].(string)

	if !ok {
		return nil, errorf(JSONRPCInvalidParams, "address parameter must be string")
	}

	asset, ok := params[1].(string)

	if !ok {
		return nil, errorf(JSONRPCInvalidParams, "asset parameter must be string")
	}

	utxos, err := server.unspent(address, asset)

	if err != nil {
		return nil, errorf(JSONRPCInnerError, "get %s balance %s err:\n\t%s", address, asset, err)
	}

	return utxos, nil
}

func (server *Server) unspent(address string, asset string) ([]*rpc.UTXO, error) {

	tutxos := make([]*neodb.UTXO, 0)

	err := server.
		engine.
		Where(`address = ? and asset = ? and spent_block = -1`, address, asset).
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

func (server *Server) getClaim(params []interface{}) (interface{}, *JSONRPCError) {

	if len(params) < 1 {
		return nil, errorf(JSONRPCInvalidParams, "expect address and asset parameters")
	}

	address, ok := params[0].(string)

	if !ok {
		return nil, errorf(JSONRPCInvalidParams, "address parameter must be string")
	}

	unclaimed, ok := server.getCachedClaim(address)

	if !ok {
		unclaimed = &rpc.Unclaimed{
			Unavailable: "0",
			Available:   "0",
		}
	}

	return unclaimed, nil
}

func (server *Server) markAddress(address string) bool {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	flag := false

	if _, ok := server.syncFlag[address]; !ok {
		server.syncFlag[address] = address

		flag = true

		logger.DebugF("queued claim task: %s", address)
	}

	return flag
}

func (server *Server) getCachedClaim(address string) (unclaimed *rpc.Unclaimed, ok bool) {

	logger.DebugF("get claim: %s", address)

	flag := server.markAddress(address)

	if flag {
		server.syncChan <- &syncAddress{
			Address: address,
			Times:   server.syncTimes,
		}
	}

	val, err := server.redisclient.Get(address).Result()

	if err == redis.Nil {
		return nil, false
	}

	if err != nil {
		logger.DebugF("get cached claim for address %s err , %s", address, err)
		return nil, false
	}

	err = json.Unmarshal([]byte(val), &unclaimed)

	if err != nil {
		logger.DebugF("get cachde claim for address %s err , %s", address, err)
		return nil, false
	}

	ok = true

	logger.DebugF("get cached claim for address %s -- success", address)

	return
}

// Asserts .
const (
	GasAssert = "0x602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
	NEOAssert = "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
)

func (server *Server) unclaimed(address string) ([]*rpc.UTXO, error) {
	tutxos := make([]*neodb.UTXO, 0)

	err := server.
		engine.
		Where(`address = ? and asset = ? and claimed = FALSE`, address, NEOAssert).
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

func (server *Server) getBlocksFee(start int64, end int64) (float64, int64, error) {

	var rows []map[string]string
	var err error

	if end == -1 {
		rows, err = server.engine.QueryString(`select sum(sys_fee), max(block) from neo_block where block >= ?`, start)
	} else {
		rows, err = server.engine.QueryString(`select sum(sys_fee), max(block) from neo_block where block >= ? and block < ?`, start, end)
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

func (server *Server) getBlocks(start int64, end int64) ([]*neodb.Block, error) {

	blocks := make([]*neodb.Block, 0)

	if end == -1 {
		err := server.engine.Where(`block >= ?`, start).Find(&blocks)
		if err != nil {
			return nil, err
		}
	} else {
		if err := server.engine.Where(`block >= ? and block <= ?`, start, end).Cols("block", "sys_fee").Find(&blocks); err != nil {
			return nil, err
		}
	}

	return blocks, nil
}

func (server *Server) doGetClaim(address string) (*rpc.Unclaimed, error) {

	logger.DebugF("[doGetClaim]start get claim :%s", address)

	utxos, err := server.unclaimed(address)

	if err != nil {
		return nil, fmt.Errorf("[doGetClaim]get %s get unclaimed utxo err:\n\t%s", address, err)
	}

	logger.DebugF("[doGetClaim]get address %s unclaimed utxo -- success", address)

	start, end := claim.CalcBlockRange(utxos)

	logger.Debug("[doGetClaim]start get blocks", start, end)
	blocks, err := server.getBlocks(start, end)
	logger.Debug("[doGetClaim]end get blocks", len(blocks), blocks[len(blocks)-1].Block)

	if err != nil {
		return nil, nil
	}

	logger.DebugF("[doGetClaim] calc address %s unclaimed gas", address)

	unavailable, available, err := claim.CalcUnclaimedGas(utxos, blocks)

	if err != nil {
		return nil, fmt.Errorf("[doGetClaim]get address %s unclaimed gas fee err:\n\t%s", address, err)
	}

	logger.DebugF("[doGetClaim] calc address %s unclaimed gas -- success", address)

	claims := make([]*rpc.UTXO, 0)

	for _, utxo := range utxos {
		if utxo.SpentBlock != -1 {
			claims = append(claims, utxo)
		}
	}

	unclaimed := &rpc.Unclaimed{
		Available:   fmt.Sprintf("%.8f", round(available, 8)),
		Unavailable: fmt.Sprintf("%.8f", round(unavailable, 8)),
		Claims:      claims,
	}

	// jsondata, _ := json.Marshal(unclaimed)

	logger.DebugF("[doGetClaim]finish get claim: %s available: %.8f unavailable: %.8f", address, round(available, 8), round(unavailable, 8))

	return unclaimed, nil
}

func round(f float64, n int) float64 {
	data := fmt.Sprintf("%.9f", f)

	data = data[0 : len(data)-1]

	r, _ := strconv.ParseFloat(data, 8)

	return r
}
