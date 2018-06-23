// automatically generated!

package main

import (
	"fmt"

	"github.com/go-hep/croot"
)

type NTExtraVars struct {
	Nvtx                  uint32
	Metjet                float32
	Metsjet               float32
	Metcell               float32
	Metel                 float32
	Metmu                 float32
	Metph                 float32
	Mettrack              float32
	Mettrack_phi          float32
	MT2                   float32
	Thrust                float32
	Thrust_axis_phi       float32
	Sphericity            float32
	MCT                   float32
	NISRJet               int32
	IsrjetPt              float32
	MT2_noISR             float32
	Thrust_noISR          float32
	Thrust_axis_phi_noISR float32
	Sphericity_noISR      float32
	MCT_noISR             float32
}

type NTVars struct {
	RunNumber        uint32
	EventNumber      uint32
	Veto             uint32
	EventWeight      float32
	PileupWeight     float32
	PileupWeightUp   float32
	PileupWeightDown float32
	GenWeight        float32
	NJet             uint32
	Jet1Pt           float32
	Jet2Pt           float32
	Jet3Pt           float32
	Jet4Pt           float32
	Jet5Pt           float32
	Jet6Pt           float32
	Jet1Eta          float32
	Jet2Eta          float32
	Jet3Eta          float32
	Jet4Eta          float32
	Jet5Eta          float32
	Jet6Eta          float32
	Jet1Phi          float32
	Jet2Phi          float32
	Jet3Phi          float32
	Jet4Phi          float32
	Jet5Phi          float32
	Jet6Phi          float32
	Jet1M            float32
	Jet2M            float32
	Jet3M            float32
	Jet4M            float32
	Jet5M            float32
	Jet6M            float32
	Met              float32
	MetPhi           float32
	DPhi             float32
	DPhiR            float32
	Meff1Jet         float32
	Meff2Jet         float32
	Meff3Jet         float32
	Meff4Jet         float32
	Meff5Jet         float32
	Meff6Jet         float32
	MeffInc          float32
	Hardproc         int32
	NBJet            uint32
	BTagWeight       float32
	BTagWeightBUp    float32
	BTagWeightBDown  float32
	BTagWeightCUp    float32
	BTagWeightCDown  float32
	BTagWeightLUp    float32
	BTagWeightLDown  float32
	NormWeight       float32
	NormWeightUp     float32
	NormWeightDown   float32
	Cleaning         uint32
	Timing           float32
	Jet1Emf          float32
	Jet2Emf          float32
	Jet1Chf          float32
	Jet2Chf          float32
	PdfId1           int32
	PdfId2           int32
	TauN             uint32
	TauJetBDTLoose   uint32
	TauMt            float32
}

type DataReader struct {
	NTExtraVars NTExtraVars
	NTVars      NTVars

	// branches
	b_NTExtraVars croot.Branch
	b_NTVars      croot.Branch

	Tree croot.Tree
}

func NewDataReader(tree croot.Tree) (*DataReader, error) {
	dr := &DataReader{}
	err := dr.Init(tree)
	if err != nil {
		return nil, err
	}
	return dr, nil
}

func (dr *DataReader) Init(tree croot.Tree) error {
	var err error
	var o int32
	dr.Tree = tree

	o = dr.Tree.SetBranchAddress("NTExtraVars", &dr.NTExtraVars)
	if o < 0 {
		return fmt.Errorf("invalid branch: [NTExtraVars] (got %d)", o)
	}

	o = dr.Tree.SetBranchAddress("NTVars", &dr.NTVars)
	if o < 0 {
		return fmt.Errorf("invalid branch: [NTVars] (got %d)", o)
	}

	return err
}

func (dr *DataReader) GetEntry(entry int64) int {
	if dr.Tree == nil {
		return 0
	}
	return dr.Tree.GetEntry(entry, 1)
}

func init() {
	// register all generated types with CRoot
	croot.RegisterType(&NTExtraVars{})
	croot.RegisterType(&NTVars{})

}
