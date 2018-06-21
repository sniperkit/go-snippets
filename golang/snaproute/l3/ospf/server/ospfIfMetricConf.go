package server

import (
	//"fmt"
	"errors"
	"l3/ospf/config"
)

func (server *OSPFServer) processIfMetricConfig(conf config.IfMetricConf) error {
	intfConfKey := IntfConfKey{
		IPAddr:  conf.IfMetricIpAddress,
		IntfIdx: conf.IfMetricAddressLessIf,
	}

	ent, exist := server.IntfConfMap[intfConfKey]
	if !exist {
		server.logger.Err("No such interface exists for applying Interface Metric Configuration")
		err := errors.New("No such interface exists for applying Interface Metric Configuration")
		return err
	}
	if conf.IfMetricTOS == 0 {
		ent.IfCost = uint32(conf.IfMetricValue)
	}
	ent.IfMetricTOSMap[uint8(conf.IfMetricTOS)] = uint32(conf.IfMetricValue)
	server.IntfConfMap[intfConfKey] = ent
	return nil
}
