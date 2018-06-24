// pdf extracts metadata from pdf files
package pdf

import (
	"bytes"
	"github.com/datatogether/xmp"
	unicommon "github.com/unidoc/unidoc/common"
	unilicense "github.com/unidoc/unidoc/license"
	pdf "github.com/unidoc/unidoc/pdf"
	"io"
	"os"
)

// MetadataForFile generates metadata from a filepath
func MetadataForFile(file string) (map[string]interface{}, error) {
	// r, err := pdf.Open(file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return extract(f)
}

// MetadataForFile generates metadata from a byte slice
func MetadataForBytes(data []byte) (map[string]interface{}, error) {
	return extract(bytes.NewReader(data))
}

func extract(r io.ReadSeeker) (map[string]interface{}, error) {
	p, err := pdf.NewParser(r)
	if err != nil {
		return nil, err
	}

	// infoRdr, _ := pdf.NewPdfReader(r)
	// info, _ := infoRdr.Inspect()
	// fmt.Println(info)

	i := 1
	for {
		// fmt.Println(i)
		o, err := p.LookupByNumber(i)
		if err != nil || o.String() == "null" {
			break
		}

		iobj, isIndirect := o.(*pdf.PdfIndirectObject)
		if isIndirect {
			// fmt.Printf("IND OOBJ %d: %s\n", xref.objectNumber, iobj)
			dict, isDict := iobj.PdfObject.(*pdf.PdfObjectDictionary)
			if isDict {
				// Check if has Type parameter.
				if ot, has := (*dict)["Type"].(*pdf.PdfObjectName); has {
					otype := string(*ot)
					// fmt.Printf("---> Obj type: %s\n", otype)
					if otype == "Catalog" {
						// fmt.Println(dict.String())
						for key, obj := range *dict {
							// TODO - check pdf spec, is only one metadata entry allowed?
							if key.String() == "Metadata" {
								oNum := obj.(*pdf.PdfObjectReference).ObjectNumber
								obj, err := p.LookupByNumber(int(oNum))
								if err != nil {
									return nil, err
								}
								if sobj, isStream := obj.(*pdf.PdfObjectStream); isStream {
									packet, err := xmp.Unmarshal(sobj.Stream)
									if err != nil {
										return nil, err
									}
									return packet.AsPOD().AsObject()
								}
							}
						}
					}
				}
				//     else if ot, has := (*dict)["Subtype"].(*pdf.PdfObjectName); has {
				// 	// Check if subtype
				// 	otype := string(*ot)
				// 	// fmt.Printf("---> Obj subtype: %s\n", otype)

				// }
				// if val, has := (*dict)["S"].(*pdf.PdfObjectName); has && *val == "JavaScript" {

				// }

			}
		}
		//   else if sobj, isStream := o.(*pdf.PdfObjectStream); isStream {
		// 	// if otype, ok := (*(sobj.PdfObjectDictionary))["Type"].(*pdf.PdfObjectName); ok {
		// 	// 	// fmt.Printf("--> Stream object type: %s\n", *otype)
		// 	// 	// if otype.String() == "Metadata" {
		// 	// 	// 	fmt.Println(string(sobj.Stream))
		// 	// 	// }
		// 	// }
		// } else if dict, isDict := o.(*pdf.PdfObjectDictionary); isDict {
		// 	ot, isName := (*dict)["Type"].(*pdf.PdfObjectName)
		// 	if isName {
		// 		// otype := string(*ot)
		// 		// fmt.Println("object type:", otype)
		// 	}
		// } else {
		// 	fmt.Println(o)
		// }
		i++
		// break
	}

	// fmt.Println(pg.GetPageAsIndirectObject())

	// fmt.Println(p.Inspect())
	// fmt.Println(p.PageList)
	// fmt.Println(pg)
	// fmt.Println(pg.GetPageAsIndirectObject())
	return nil, nil
}

// extract pulls metadata from a pdf reader
// func extract(r *pdf.Reader) (map[string]interface{}, error) {
// 	for i := 1; i <= r.NumPage(); i++ {
// 		fmt.Printf("interpret page %d\n", i)

// 		pdf.Interpret(r.Page(i).Resources(), func(stk *pdf.Stack, op string) {
// 			fmt.Println(op)
// 			fmt.Println(stk)
// 		})
// 		fmt.Printf("interpreted page %d\n", i)
// 	}
// 	return nil, nil
// }

func init() {
	initUniDoc("")
}

func initUniDoc(licenseKey string) error {
	if len(licenseKey) > 0 {
		err := unilicense.SetLicenseKey(licenseKey)
		if err != nil {
			return err
		}
	}

	// To make the library log we just have to initialise the logger which satisfies
	// the unicommon.Logger interface, unicommon.DummyLogger is the default and
	// does not do anything. Very easy to implement your own.
	unicommon.SetLogger(unicommon.DummyLogger{})

	return nil
}
