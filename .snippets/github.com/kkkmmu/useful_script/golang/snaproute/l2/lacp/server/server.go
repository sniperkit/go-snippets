package server

import (
	"asicd/asicdCommonDefs"
	"fmt"
	//"infra/sysd/sysdCommonDefs"
	"l2/lacp/protocol/drcp"
	"l2/lacp/protocol/lacp"
	"l2/lacp/protocol/utils"
	"utils/commonDefs"
	//"utils/keepalive"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
)

type LaConfigMsgType int8

const (
	LAConfigMsgCreateLaPortChannel LaConfigMsgType = iota + 1
	LAConfigMsgDeleteLaPortChannel
	LAConfigMsgUpdateLaPortChannelLagHash
	LAConfigMsgUpdateLaPortChannelSystemIdMac
	LAConfigMsgUpdateLaPortChannelSystemPriority
	LAConfigMsgUpdateLaPortChannelLagType
	LAConfigMsgUpdateLaPortChannelAdminState
	LAConfigMsgUpdateLaPortChannelAggMode
	LAConfigMsgUpdateLaPortChannelPeriod
	LAConfigMsgCreateLaAggPort
	LAConfigMsgDeleteLaAggPort
	LAConfigMsgUpdateLaAggPortAdminState
	LAConfigMsgCreateDistributedRelay
	LAConfigMsgDeleteDistributedRelay
	LAConfigMsgAggregatorCreated
	LAConfigMsgCreateConversationId
	LAConfigMsgUpdateConversationId
	LAConfigMsgDeleteConversationId
	LAConfigMsgAddL3IntfType
	LAConfigMsgAddL2IntfType
)

type LAConfig struct {
	Msgtype LaConfigMsgType
	Msgdata interface{}
}

type LAServer struct {
	logger           *logging.Writer
	ConfigCh         chan LAConfig
	AsicdSubSocketCh chan commonDefs.AsicdNotifyMsg
	eventDbHdl       *dbutils.DBUtil
}

func NewLAServer(logger *logging.Writer) *LAServer {
	return &LAServer{
		logger:           logger,
		ConfigCh:         make(chan LAConfig),
		AsicdSubSocketCh: make(chan commonDefs.AsicdNotifyMsg),
	}
}

func (server *LAServer) InitServer() {
	utils.ConstructPortConfigMap()

	err := server.initializeEvents()
	if err != nil {
		utils.GetLaLogger().Err("Error initializing Event Db")
	}
	// TODO
	//go server.ListenToClientStateChanges()
	server.StartLaConfigNotificationListener()
	drcp.GetAllCVIDConversations()
}

func (server *LAServer) initializeEvents() error {
	logger := utils.GetLaLogger()
	server.eventDbHdl = dbutils.NewDBUtil(logger)
	err := server.eventDbHdl.Connect()
	if err != nil {
		utils.GetLaLogger().Err("Failed to create the DB handle")
		return err
	}

	return eventUtils.InitEvents("LACPD", server.eventDbHdl, server.eventDbHdl, logger, 1000)
}

/*
TODO
func (server *LAServer) ListenToClientStateChanges() {
	clientStatusListener := keepalive.InitDaemonStatusListener()
	if clientStatusListener != nil {
		go clientStatusListener.StartDaemonStatusListner()
		for {
			select {
			case clientStatus := <-clientStatusListener.DaemonStatusCh:
				svr.logger.Info(fmt.Sprintln("Received client status: ", clientStatus.Name, clientStatus.Status))
				if svr.IsReady() {
					switch clientStatus.Status {
					case sysdCommonDefs.STOPPED, sysdCommonDefs.RESTARTING:
						go svr.DisconnectFromClient(clientStatus.Name)
					case sysdCommonDefs.UP:
						go svr.ConnectToClient(clientStatus.Name)
					}
				}
			}
		}
	}
}
*/
// StartSTPSConfigNotificationListener
func (s *LAServer) StartLaConfigNotificationListener() {
	//server.InitDone <- true
	go func(svr *LAServer) {
		svr.logger.Info("Starting LA Config Event Listener")
		for {
			select {
			case laConf, ok := <-svr.ConfigCh:
				if ok {
					svr.processLaConfig(laConf)
				} else {
					// channel was closed
					return
				}
			case msg := <-svr.AsicdSubSocketCh:
				svr.processAsicdNotification(msg)
			}
		}
	}(s)
}

func (s *LAServer) processLaConfig(conf LAConfig) {

	switch conf.Msgtype {
	case LAConfigMsgCreateLaPortChannel:
		s.logger.Info("CONFIG: Create Link Aggregation Group / Port Channel")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		lacp.CreateLaAgg(config)

	case LAConfigMsgDeleteLaPortChannel:
		s.logger.Info("CONFIG: Delete Link Aggregation Group / Port Channel")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		lacp.DeleteLaAgg(config.Id)

	case LAConfigMsgUpdateLaPortChannelLagHash:
		s.logger.Info("CONFIG: Link Aggregation Group / Port Channel Lag Hash Mode")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		lacp.SetLaAggHashMode(config.Id, config.HashMode)

	case LAConfigMsgUpdateLaPortChannelSystemIdMac:
		s.logger.Info("CONFIG: Link Aggregation Group / Port Channel SystemId MAC")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		var a *lacp.LaAggregator
		if lacp.LaFindAggById(config.Id, &a) {
			// configured ports
			for _, pId := range a.PortNumList {
				lacp.SetLaAggPortSystemInfo(uint16(pId), config.Lacp.SystemIdMac, config.Lacp.SystemPriority)
			}
		}

	case LAConfigMsgUpdateLaPortChannelSystemPriority:
		s.logger.Info("CONFIG: Link Aggregation Group / Port Channel System Priority")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		var a *lacp.LaAggregator
		if lacp.LaFindAggById(config.Id, &a) {
			// configured ports
			for _, pId := range a.PortNumList {
				lacp.SetLaAggPortSystemInfo(uint16(pId), config.Lacp.SystemIdMac, config.Lacp.SystemPriority)
			}
		}

	case LAConfigMsgUpdateLaPortChannelLagType, LAConfigMsgUpdateLaPortChannelAggMode:
		s.logger.Info("CONFIG: Link Aggregation Group / Port Channel System Lag Type")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		var a *lacp.LaAggregator
		var p *lacp.LaAggPort
		if lacp.LaFindAggById(config.Id, &a) {

			if config.Type == lacp.LaAggTypeSTATIC {
				// configured ports
				for _, pId := range a.PortNumList {
					if lacp.LaFindPortById(uint16(pId), &p) {
						lacp.SetLaAggPortLacpMode(uint16(pId), lacp.LacpModeOn)
					}
				}
			} else {
				for _, pId := range a.PortNumList {
					if lacp.LaFindPortById(uint16(pId), &p) {
						lacp.SetLaAggPortLacpMode(uint16(pId), int(config.Lacp.Mode))
					}
				}
			}
		}

	case LAConfigMsgUpdateLaPortChannelAdminState:
		s.logger.Info("CONFIG: Link Aggregation Group / Port Channel System Lag Type")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		if config.Enabled {
			lacp.EnableLaAgg(config.Id)
		} else {
			lacp.DisableLaAgg(config.Id)
		}
	case LAConfigMsgUpdateLaPortChannelPeriod:
		s.logger.Info("CONFIG: Link Aggregation Group / Port Channel System Period")
		config := conf.Msgdata.(*lacp.LaAggConfig)
		var a *lacp.LaAggregator
		if lacp.LaFindAggById(config.Id, &a) {
			// configured ports
			for _, pId := range a.PortNumList {
				lacp.SetLaAggPortLacpPeriod(uint16(pId), config.Lacp.Interval)
			}
		}
	case LAConfigMsgCreateLaAggPort:
		s.logger.Info("CONFIG: Create Link Aggregation Port")
		config := conf.Msgdata.(*lacp.LaAggPortConfig)
		lacp.CreateLaAggPort(config)

	case LAConfigMsgDeleteLaAggPort:
		s.logger.Info("CONFIG: Delete Link Aggregation Port")
		config := conf.Msgdata.(*lacp.LaAggPortConfig)
		lacp.DeleteLaAggPort(config.Id)

	case LAConfigMsgCreateDistributedRelay:
		s.logger.Info("CONFIG: Create Distributed Relay")
		config := conf.Msgdata.(*drcp.DistributedRelayConfig)
		drcp.CreateDistributedRelay(config)

	case LAConfigMsgDeleteDistributedRelay:
		s.logger.Info("CONFIG: Delete Distributed Relay")
		config := conf.Msgdata.(*drcp.DistributedRelayConfig)
		drcp.DeleteDistributedRelay(config.GetKey())

	case LAConfigMsgCreateConversationId:
		s.logger.Info("CONFIG: Create Conversation Id")
		config := conf.Msgdata.(*drcp.DRConversationConfig)
		drcp.CreateConversationId(config)

	case LAConfigMsgDeleteConversationId:
		s.logger.Info("CONFIG: Delete Conversation Id")
		config := conf.Msgdata.(*drcp.DRConversationConfig)
		drcp.DeleteConversationId(config, false)

	case LAConfigMsgUpdateConversationId:
		s.logger.Info("CONFIG: Update Conversation Id")
		config := conf.Msgdata.(*drcp.DRConversationConfig)
		drcp.UpdateConversationId(config)

	case LAConfigMsgAddL2IntfType:
		s.logger.Info("CONFIG: Update L2 Intf")
		config := conf.Msgdata.(*commonDefs.IPv4L3IntfStateNotifyMsg)
		lacp.UpdateIntfType(int(config.IfIndex), "L2")

	case LAConfigMsgAddL3IntfType:
		s.logger.Info("CONFIG: Update L3 Intf")
		config := conf.Msgdata.(*commonDefs.IPv4L3IntfStateNotifyMsg)
		lacp.UpdateIntfType(int(config.IfIndex), "L3")
	}
}

func (s *LAServer) processLinkDownEvent(linkId int) {
	s.logger.Info(fmt.Sprintln("LA EVT: Link Down", linkId))
	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(linkId), &p) {
		p.DeleteRxTx()
		p.LinkOperStatus = false
		lacp.DisableLaAggPort(uint16(linkId))
	} else {
		for _, ipp := range drcp.DRCPIppDBList {
			if int(ipp.Id) == linkId {
				go ipp.DrIppLinkDown()
			}
		}
	}
}

func (s *LAServer) processLinkUpEvent(linkId int) {
	s.logger.Info(fmt.Sprintln("LA EVT: Link Up", linkId))
	var p *lacp.LaAggPort
	if lacp.LaFindPortById(uint16(linkId), &p) {
		p.CreateRxTx()
		p.LinkOperStatus = true
		lacp.EnableLaAggPort(uint16(linkId))

	} else {
		for _, ipp := range drcp.DRCPIppDBList {
			if int(ipp.Id) == linkId {
				go ipp.DrIppLinkUp()
			}
		}
	}
}

func (s *LAServer) processVlanEvent(vlanMsg commonDefs.VlanNotifyMsg) {

	// need to determine whether the message is a new vlan or if
	// ports were updated are any of them part of any of the
	// aggregators which exist
	msgtype := LAConfigMsgUpdateConversationId
	if vlanMsg.MsgType == commonDefs.NOTIFY_VLAN_CREATE {
		msgtype = LAConfigMsgCreateConversationId
	} else if vlanMsg.MsgType == commonDefs.NOTIFY_VLAN_DELETE {
		msgtype = LAConfigMsgDeleteConversationId
	}
	var agg *lacp.LaAggregator
	for lacp.LaGetAggNext(&agg) {
		portList := vlanMsg.TagPorts
		for _, p := range vlanMsg.UntagPorts {
			portList = append(portList, p)
		}

		for _, uifindex := range portList {
			id := asicdCommonDefs.GetIntfIdFromIfIndex(uifindex)
			for _, aggport := range agg.PortNumList {
				if id == int(aggport) {
					var dr *drcp.DistributedRelay
					if drcp.DrFindByAggregator(int32(agg.AggId), &dr) {
						// gateway message
						cfg := drcp.DRConversationConfig{
							DrniName: dr.DrniName,
							Idtype:   drcp.GATEWAY_ALGORITHM_CVID,
							Cvlan:    uint16(vlanMsg.VlanId),
						}

						s.ConfigCh <- LAConfig{
							Msgtype: msgtype,
							Msgdata: cfg,
						}
					}
				}
			}
		}
	}
}

func (s *LAServer) processL3IntEvent(msg commonDefs.IPv4L3IntfStateNotifyMsg) {

	ifindex := msg.IfIndex
	iftype := commonDefs.GetIfTypeName(asicdCommonDefs.GetIntfTypeFromIfIndex(ifindex))
	if iftype == "Lag" {
		if msg.MsgType == commonDefs.NOTIFY_IPV4INTF_CREATE {
			s.ConfigCh <- LAConfig{
				Msgtype: LAConfigMsgAddL3IntfType,
				Msgdata: msg,
			}
		} else if msg.MsgType == commonDefs.NOTIFY_IPV4INTF_DELETE {
			s.ConfigCh <- LAConfig{
				Msgtype: LAConfigMsgAddL2IntfType,
				Msgdata: msg,
			}
		} else if msg.MsgType == commonDefs.NOTIFY_IPV6INTF_CREATE {
			s.ConfigCh <- LAConfig{
				Msgtype: LAConfigMsgAddL3IntfType,
				Msgdata: msg,
			}
		} else if msg.MsgType == commonDefs.NOTIFY_IPV6INTF_DELETE {
			s.ConfigCh <- LAConfig{
				Msgtype: LAConfigMsgAddL2IntfType,
				Msgdata: msg,
			}
		}
	}
}

func (s *LAServer) processAsicdNotification(msg commonDefs.AsicdNotifyMsg) {
	switch msg.(type) {
	case commonDefs.L2IntfStateNotifyMsg:
		l2Msg := msg.(commonDefs.L2IntfStateNotifyMsg)
		s.logger.Info(fmt.Sprintf("Msg linkstatus = %d msg port = %d\n", l2Msg.IfState, l2Msg.IfIndex))
		if l2Msg.IfState == asicdCommonDefs.INTF_STATE_DOWN {
			s.processLinkDownEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l2Msg.IfIndex)) //asicd always sends out link State events for PHY ports
		} else {
			s.processLinkUpEvent(asicdCommonDefs.GetIntfIdFromIfIndex(l2Msg.IfIndex))
		}
	case commonDefs.VlanNotifyMsg:
		vlanMsg := msg.(commonDefs.VlanNotifyMsg)
		s.logger.Info(fmt.Sprintln("Msg vlan = ", vlanMsg))
		s.processVlanEvent(vlanMsg)

	case commonDefs.IPv4L3IntfStateNotifyMsg:
		l3intfMsg := msg.(commonDefs.IPv4L3IntfStateNotifyMsg)
		s.logger.Info(fmt.Sprintln("Msg l3intf = ", l3intfMsg))
		s.processL3IntEvent(l3intfMsg)
	}
}
