package stitch

import (
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/fs"
)

/*
	Reduce a formula to a slice of []UnpackSpec, ready to be used
	invoking Assembler.Run().

	Typically the filters arg will be either `api.Filter_NoMutation` or
	`api.Filter_LowPriv`, depending if you're using repeatr or rio respectively,
	though other values are of course valid.

	Whether the action, outputs, or saveUrls are set is irrelevant;
	they will be ignored completely.
*/
func FormulaToUnpackSpecs(frm api.Formula, frmCtx api.FormulaContext, filters api.FilesetFilters) (parts []UnpackSpec) {
	for path, wareID := range frm.Inputs {
		warehouses, _ := frmCtx.FetchUrls[path]
		parts = append(parts, UnpackSpec{
			Path:       fs.MustAbsolutePath(string(path)),
			WareID:     wareID,
			Filters:    filters,
			Warehouses: warehouses,
		})
	}
	return
}

/*
	Reduce a formula to a slice of []PackSpec, ready to be used
	invoking PackMulti().

	The filters given will be applied to all *unset* fields in the filters
	already given by the formula outputs; set fields are not changed.
	Typically the filters arg will be `api.Filter_DefaultFlatten`
	(both `rio pack` and repeatr outputs default to this),
	though other values are of course valid.

	Whether the action, inputs, or fetchUrls are set is irrelevant;
	they will be ignored completely.
*/
func FormulaToPackSpecs(frm api.Formula, frmCtx api.FormulaContext, filters api.FilesetFilters) (parts []PackSpec) {
	for path, output := range frm.Outputs {
		warehouse, _ := frmCtx.SaveUrls[path]
		parts = append(parts, PackSpec{
			Path:      fs.MustAbsolutePath(string(path)),
			PackType:  output.PackType,
			Filters:   apiutil.MergeFilters(output.Filters, filters),
			Warehouse: warehouse,
		})
	}
	return
}
