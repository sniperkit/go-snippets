package n2xsdk

type MosEstimation struct {
 Handler string
}

func(np *MosEstimation) IsMosEstimationSupported () error {
 //parameters: PortHandle
 //AgtMosEstimation IsMosEstimationSupported
 return nil
}

func(np *MosEstimation) SetFactors () error {
 //parameters: Count psaStreamGroupHandles EquipmentImpairmentFactor PacketLossRobustnessFactor
 //AgtMosEstimation SetFactors
 return nil
}

func(np *MosEstimation) GetFactors ()(string, error) {
 //parameters: Count psaStreamGroupHandles
 //AgtMosEstimation GetFactors
 return "", nil
}

func(np *MosEstimation) ClearFactors () error {
 //parameters: Count psaStreamGroupHandles
 //AgtMosEstimation ClearFactors
 return nil
}

func(np *MosEstimation) IsFactorSet () error {
 //parameters: Count psaStreamGroupHandles
 //AgtMosEstimation IsFactorSet
 return nil
}

