package n2xsdk

type IptvChannelPool struct {
 Handler string
}

func(np *IptvChannelPool) SetRandomWatchTime () error {
 //parameters: ScenarioHandle GroupPoolHandle WatchTimeRange
 //AgtIptvChannelPool SetRandomWatchTime
 return nil
}

func(np *IptvChannelPool) GetRandomWatchTime ()(string, error) {
 //parameters: ScenarioHandle GroupPoolHandle
 //AgtIptvChannelPool GetRandomWatchTime
 return "", nil
}

func(np *IptvChannelPool) SetWatchTime () error {
 //parameters: ScenarioHandle GroupPoolHandle WatchTime
 //AgtIptvChannelPool SetWatchTime
 return nil
}

func(np *IptvChannelPool) GetWatchTime ()(string, error) {
 //parameters: ScenarioHandle GroupPoolHandle
 //AgtIptvChannelPool GetWatchTime
 return "", nil
}

func(np *IptvChannelPool) SetStbDelay () error {
 //parameters: ScenarioHandle GroupPoolHandle StbDelay
 //AgtIptvChannelPool SetStbDelay
 return nil
}

func(np *IptvChannelPool) GetStbDelay ()(string, error) {
 //parameters: ScenarioHandle GroupPoolHandle
 //AgtIptvChannelPool GetStbDelay
 return "", nil
}

func(np *IptvChannelPool) AddSourcePools () error {
 //parameters: ScenarioHandle GroupPoolHandle SourcePoolHandles
 //AgtIptvChannelPool AddSourcePools
 return nil
}

func(np *IptvChannelPool) RemoveSourcePools () error {
 //parameters: ScenarioHandle GroupPoolHandle SourcePoolHandles
 //AgtIptvChannelPool RemoveSourcePools
 return nil
}

func(np *IptvChannelPool) RemoveAllSourcePools () error {
 //parameters: ScenarioHandle GroupPoolHandle
 //AgtIptvChannelPool RemoveAllSourcePools
 return nil
}

func(np *IptvChannelPool) ListSourcePools ()(string, error) {
 //parameters: ScenarioHandle GroupPoolHandle
 //AgtIptvChannelPool ListSourcePools
 return "", nil
}

func(np *IptvChannelPool) SetGroupPoolFilter () error {
 //parameters: ScenarioHandle GroupPoolHandle GroupPoolFilter
 //AgtIptvChannelPool SetGroupPoolFilter
 return nil
}

func(np *IptvChannelPool) GetGroupPoolFilter ()(string, error) {
 //parameters: ScenarioHandle GroupPoolHandle
 //AgtIptvChannelPool GetGroupPoolFilter
 return "", nil
}

