package tribserver

import (
	//	"errors"
	"encoding/json"
	"fmt"
	"github.com/cmu440/tribbler/libstore"
	"github.com/cmu440/tribbler/rpc/tribrpc"
	"github.com/cmu440/tribbler/util"
	"net"
	"net/http"
	"net/rpc"
	"sort"
	"time"
)

type tribServer struct {
	lib libstore.Libstore
}

// NewTribServer creates, starts and returns a new TribServer. masterServerHostPort
// is the master storage server's host:port and port is this port number on which
// the TribServer should listen. A non-nil error should be returned if the TribServer
// could not be started.
//
// For hints on how to properly setup RPC, see the rpc/tribrpc package.
func NewTribServer(masterServerHostPort, myHostPort string) (TribServer, error) {
	lib, err := libstore.NewLibstore(masterServerHostPort, myHostPort, libstore.Normal)
	tribServer := &tribServer{lib: lib}
	// Get the port number
	_, port, err := net.SplitHostPort(myHostPort)
	// Create the server socket that will listen for incoming RPCs.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}
	// Wrap the tribServer before registering it for RPC.
	err = rpc.RegisterName("TribServer", tribrpc.Wrap(tribServer))
	if err != nil {
		return nil, err
	}
	// Setup the HTTP handler that will server incoming RPCs and
	// serve requests in a background goroutine.
	rpc.HandleHTTP()
	go http.Serve(listener, nil)
	return tribServer, nil
}

func (ts *tribServer) CreateUser(args *tribrpc.CreateUserArgs, reply *tribrpc.CreateUserReply) error {
	userID := args.UserID
	userIDKey := util.FormatUserKey(userID)
	result, err := ts.lib.Get(userIDKey)
	if result != "" {
		reply.Status = tribrpc.Exists
		return nil
	}
	err = ts.lib.Put(userIDKey, userID)
	if err != nil {
		return err
	}
	reply.Status = tribrpc.OK
	return nil
}

func (ts *tribServer) AddSubscription(args *tribrpc.SubscriptionArgs, reply *tribrpc.SubscriptionReply) error {
	userID := args.UserID
	targetUserID := args.TargetUserID
	//First check if user and target user exist
	userIDKey := util.FormatUserKey(userID)
	targetUserIDKey := util.FormatUserKey(targetUserID)
	result, err := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	result, err = ts.lib.Get(targetUserIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchTargetUser
		return nil
	}
	userSubKey := util.FormatSubListKey(userID)
	err = ts.lib.AppendToList(userSubKey, targetUserID)
	if err != nil {
		reply.Status = tribrpc.Exists
		return nil
	}
	reply.Status = tribrpc.OK
	return nil
}

func (ts *tribServer) RemoveSubscription(args *tribrpc.SubscriptionArgs, reply *tribrpc.SubscriptionReply) error {
	userID := args.UserID
	targetUserID := args.TargetUserID
	//First check if user and target user exist
	userIDKey := util.FormatUserKey(userID)
	targetUserIDKey := util.FormatUserKey(targetUserID)
	result, err := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	result, err = ts.lib.Get(targetUserIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchTargetUser
		return nil
	}
	userSubKey := util.FormatSubListKey(userID)
	err = ts.lib.RemoveFromList(userSubKey, targetUserID)
	if err != nil {
		reply.Status = tribrpc.NoSuchTargetUser
		return nil
	}
	reply.Status = tribrpc.OK
	return nil
}

func (ts *tribServer) GetSubscriptions(args *tribrpc.GetSubscriptionsArgs, reply *tribrpc.GetSubscriptionsReply) error {
	userID := args.UserID
	//First check if the user exists
	userIDKey := util.FormatUserKey(userID)
	result, err := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	userSubKey := util.FormatSubListKey(userID)
	list, err := ts.lib.GetList(userSubKey)
	// TODO: Check what this error message is, and accordingly return stuff
	// for example, if err == "KeyNotFound", return blank list.
	// OR, as an alternative to the above approach, we could return nil error but with a blank list
	if err != nil {
		return err
	}
	reply.Status = tribrpc.OK
	reply.UserIDs = list
	return nil
}

func (ts *tribServer) PostTribble(args *tribrpc.PostTribbleArgs, reply *tribrpc.PostTribbleReply) error {
	//timestamp the tribble
	postTime := time.Now()
	contents := args.Contents
	userID := args.UserID
	//First check if the user exists
	userIDKey := util.FormatUserKey(userID)
	result, err := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	tribListKey := util.FormatTribListKey(userID)
	postKey := util.FormatPostKey(userID, postTime.UnixNano())
	tribble := tribrpc.Tribble{
		UserID:   userID,
		Posted:   postTime,
		Contents: contents}
	marshaledTribble, err := json.Marshal(tribble)
	if err != nil {
		fmt.Println("PostTribble error:", err)
		return err
	}
	//update tribble list and user's post list on storage server
	err = ts.lib.Put(postKey, string(marshaledTribble))
	if err != nil {
		fmt.Println("PostTribble error:", err)
		return err
	}
	err = ts.lib.AppendToList(tribListKey, postKey)
	if err != nil {
		fmt.Println("PostTribble error:", err)
		return err
	}
	reply.PostKey = postKey
	reply.Status = tribrpc.OK
	return nil
}

func (ts *tribServer) DeleteTribble(args *tribrpc.DeleteTribbleArgs, reply *tribrpc.DeleteTribbleReply) error {
	postKey := args.PostKey
	userID := args.UserID
	//First check if the user exists
	userIDKey := util.FormatUserKey(userID)
	result, err := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	result, err = ts.lib.Get(postKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchPost
		return nil
	}
	//update tribble list and user's post list on storage server
	tribListKey := util.FormatTribListKey(userID)
	err = ts.lib.RemoveFromList(tribListKey, postKey)
	if err != nil {
		return err
	}
	err = ts.lib.Delete(postKey)
	if err != nil {
		return err
	}
	reply.Status = tribrpc.OK
	return nil
}

func (ts *tribServer) GetTribbles(args *tribrpc.GetTribblesArgs, reply *tribrpc.GetTribblesReply) error {
	userID := args.UserID
	//First check if the user exists
	userIDKey := util.FormatUserKey(userID)
	result, _ := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	//Then get the tribble list of this user
	tribListKey := util.FormatTribListKey(userID)
	list, _ := ts.lib.GetList(tribListKey)
	//Next get all posts, up to 100
	tribbles := make([]tribrpc.Tribble, 0)
	count := 0
	for i := len(list) - 1; i >= 0; i-- {
		postKey := list[i]
		marshaledTribble, err := ts.lib.Get(postKey)
		var tribble tribrpc.Tribble
		json.Unmarshal([]byte(marshaledTribble), &tribble)
		if err == nil {
			tribbles = append(tribbles, tribble)
			count += 1
		}
		if count >= 100 {
			break
		}
	}
	//fmt.Println("GetTribbles returning", len(tribbles), "results")
	reply.Tribbles = tribbles
	reply.Status = tribrpc.OK
	return nil
}

func (ts *tribServer) GetTribblesBySubscription(args *tribrpc.GetTribblesArgs, reply *tribrpc.GetTribblesReply) error {
	userID := args.UserID
	//First check if the user exists
	userIDKey := util.FormatUserKey(userID)
	result, _ := ts.lib.Get(userIDKey)
	if result == "" {
		reply.Status = tribrpc.NoSuchUser
		return nil
	}
	//Then get the subscription list of this user
	subListKey := util.FormatSubListKey(userID)
	subList, _ := ts.lib.GetList(subListKey)
	tribbles := make([]tribrpc.Tribble, 0)
	for _, targetUserID := range subList {
		//Then get the tribble list of this target user
		tribListKey := util.FormatTribListKey(targetUserID)
		list, _ := ts.lib.GetList(tribListKey)
		//Next get all posts, up to 100
		count := 0
		for i := len(list) - 1; i >= 0; i-- {
			postKey := list[i]
			marshaledTribble, err := ts.lib.Get(postKey)
			if err == nil {
				var tribble tribrpc.Tribble
				json.Unmarshal([]byte(marshaledTribble), &tribble)
				tribbles = append(tribbles, tribble)
				count += 1
			} else {
				fmt.Println("GetTribblesBySubscription error:", err)
			}
			if count >= 100 {
				break
			}
		}
	}
	sort.Sort(ByTime(tribbles))
	length := len(tribbles)
	if length > 100 {
		length = 100
	}
	//fmt.Println("GetTribblesBySubscription returning", length, "results")
	reply.Status = tribrpc.OK
	reply.Tribbles = tribbles[:length]
	return nil
}

type ByTime []tribrpc.Tribble

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//Actually returns greater than since we're getting the most recent ones
func (a ByTime) Less(i, j int) bool { return a[i].Posted.UnixNano() > a[j].Posted.UnixNano() }
