package n2xsdk

type MplsLabeledRoutePool struct {
 Handler string
}

func(np *MplsLabeledRoutePool) EnableLabeling () error {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool EnableLabeling
 return nil
}

func(np *MplsLabeledRoutePool) DisableLabeling () error {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool DisableLabeling
 return nil
}

func(np *MplsLabeledRoutePool) SetLabelingMode () error {
 //parameters: SessionHandle RoutePoolHandle LabelingMode
 //AgtMplsLabeledRoutePool SetLabelingMode
 return nil
}

func(np *MplsLabeledRoutePool) SetUserLabelValue () error {
 //parameters: SessionHandle RoutePoolHandle LabelValue
 //AgtMplsLabeledRoutePool SetUserLabelValue
 return nil
}

func(np *MplsLabeledRoutePool) IsLabelingEnabled () error {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool IsLabelingEnabled
 return nil
}

func(np *MplsLabeledRoutePool) GetLabelingMode ()(string, error) {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool GetLabelingMode
 return "", nil
}

func(np *MplsLabeledRoutePool) GetUserLabelValue ()(string, error) {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool GetUserLabelValue
 return "", nil
}

func(np *MplsLabeledRoutePool) GetLabels ()(string, error) {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool GetLabels
 return "", nil
}

func(np *MplsLabeledRoutePool) GetState ()(string, error) {
 //parameters: SessionHandle RoutePoolHandle
 //AgtMplsLabeledRoutePool GetState
 return "", nil
}

