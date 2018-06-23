package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/gophergala/food-plan-organizer/_third_party/github.com/rubenv/sql-migrate"
	"github.com/gophergala/food-plan-organizer/_third_party/golang.org/x/text/encoding"
	"github.com/gophergala/food-plan-organizer/_third_party/golang.org/x/text/encoding/charmap"
	"github.com/gophergala/food-plan-organizer/_third_party/golang.org/x/text/transform"

	"flag"

	_ "github.com/gophergala/food-plan-organizer/_third_party/github.com/mattn/go-sqlite3"
	"github.com/gophergala/food-plan-organizer/cmd/parse/etl"
	"github.com/gophergala/food-plan-organizer/models"
)

//go:generate go-bindata -pkg main -o bindata.go migrations

var (
	dataDirectory = flag.String("data.dir", "", "sr27 data directory")
	database      = flag.String("database", "", "path to database file")
)

type Handler func(interface{}, *sql.DB)

func InsertFood(f *models.Food, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO foods VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`, f.ID, f.FoodGroupID, f.Name, f.ShortName, f.CommonName, f.ScientificName, f.NitrogenFactor, f.ProteinFactor, f.FatFactor, f.CarbonhydrateFactor); err != nil {
		return err
	}
	return nil
}

func persistFood(i interface{}, tx *sql.DB) {
	if fo, ok := i.(etl.Food); ok {
		var dbFood = models.Food{
			ID:                  fo.ID,
			FoodGroupID:         fo.FoodGroupID,
			Name:                fo.Name,
			ShortName:           fo.ShortName,
			CommonName:          fo.CommonName,
			ScientificName:      fo.ScientificName,
			NitrogenFactor:      fo.NitrogenFactor,
			ProteinFactor:       fo.ProteinFactor,
			FatFactor:           fo.FatFactor,
			CarbonhydrateFactor: fo.CarbohydrateFactor,
		}
		if ins := InsertFood(&dbFood, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.Food")
	}
}

func InsertFoodGroup(fg *models.FoodGroup, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO food_groups VALUES ($1,$2);`, fg.ID, fg.Name); err != nil {
		return err
	}
	return nil
}

func persistFoodGroup(i interface{}, tx *sql.DB) {
	if fg, ok := i.(etl.FoodGroup); ok {
		var dbGroup = models.FoodGroup{
			ID:   fg.GroupID,
			Name: fg.Name,
		}
		if ins := InsertFoodGroup(&dbGroup, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.FoodGroup")
	}
}

func InsertLanguageDescription(ld *models.LanguageDescription, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO language_descriptions VALUES ($1, $2);`, ld.Code, ld.Description); err != nil {
		return err
	}
	return nil
}

func persistLangDescription(i interface{}, tx *sql.DB) {
	if ld, ok := i.(etl.LanguageDescription); ok {
		var dbLangDesc = models.LanguageDescription{
			Code:        ld.FactorCode,
			Description: ld.Description,
		}
		if ins := InsertLanguageDescription(&dbLangDesc, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.LangDescription")
	}
}

func InsertLanguage(l *models.Language, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO languages VALUES ($1, $2);`, l.NutrientID, l.FactorCode); err != nil {
		return err
	}
	return nil
}

func persistLang(i interface{}, tx *sql.DB) {
	if l, ok := i.(etl.Lang); ok {
		var dbL = models.Language{
			NutrientID: l.NutrientID,
			FactorCode: l.FactorCode,
		}
		if ins := InsertLanguage(&dbL, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.Lang")
	}
}

func InsertWeight(w *models.Weight, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO weights VALUES ($1, $2, $3, $4, $5);`, w.NutrientID, w.Seq, w.Amount, w.Description, w.GramWeight); err != nil {
		return err
	}
	return nil
}

func persistWeight(i interface{}, tx *sql.DB) {
	if w, ok := i.(etl.Weight); ok {
		var dbW = models.Weight{
			NutrientID:  w.NutrientID,
			Seq:         w.Seq,
			Amount:      w.Amount,
			Description: w.Description,
			GramWeight:  w.GramWeight,
		}
		if ins := InsertWeight(&dbW, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.Weight")
	}
}

func InsertNutrientDefinition(nd *models.NutrientDefinition, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO nutrient_definitions VALUES ($1,$2,$3,$4,$5);`, nd.NutrientID, nd.Units, nd.Tagname, nd.Description, nd.DecimalPlaces); err != nil {
		return err
	}
	return nil
}

func persistNutrientDefinitions(i interface{}, tx *sql.DB) {
	if nd, ok := i.(etl.NutrientDefinition); ok {
		var dbND = models.NutrientDefinition{
			NutrientID:    nd.NutrientID,
			Units:         nd.Units,
			Tagname:       nd.Tagname,
			Description:   nd.Description,
			DecimalPlaces: nd.DecimalPlaces,
		}
		if ins := InsertNutrientDefinition(&dbND, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.NutrientDefinition")
	}
}

func InsertNutrient(n *models.Nutrient, tx *sql.DB) error {
	if _, err := tx.Exec(`INSERT INTO nutrients VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`, n.FoodID, n.NutrientID, n.NutritionValue, n.Min, n.Max, n.DegreesOfFreedom, n.LowerErrorBound, n.UpperErrorBound); err != nil {
		return err
	}
	return nil
}

func persistNutrient(i interface{}, tx *sql.DB) {
	if n, ok := i.(etl.Nutrient); ok {
		var dbN = models.Nutrient{
			FoodID:           n.FoodID,
			NutrientID:       n.NutrientID,
			NutritionValue:   n.NutritionValue,
			Min:              n.Min,
			Max:              n.Max,
			DegreesOfFreedom: n.DegreesOfFreedom,
			LowerErrorBound:  n.LowerErrorBound,
			UpperErrorBound:  n.UpperErrorBound,
		}
		if ins := InsertNutrient(&dbN, tx); ins != nil {
			panic(ins)
		}
	} else {
		panic("expected etl.Nutrient")
	}
}

type file struct {
	Name        string
	Encoding    *encoding.Encoding
	Extractor   etl.Extract
	h           Handler
	Description string
}

var files = []file{
	file{
		Name:        "FOOD_DES.txt",
		Description: "Food Description",
		Extractor:   etl.FoodExtractor{},
		h:           persistFood,
	},
	file{
		Name:        "FD_GROUP.txt",
		Description: "Food Group Description",
		Extractor:   etl.FoodGroupExtractor{},
		h:           persistFoodGroup,
	},
	file{
		Name:        "LANGDESC.txt",
		Description: "Langual Factors Description",
		Extractor:   etl.LangDescriptionExtractor{},
		h:           persistLangDescription,
	},
	file{
		Name:        "LANGUAL.txt",
		Description: "Langual Factor",
		Extractor:   etl.LangExtractor{},
		h:           persistLang,
	},
	file{
		Name:        "NUTR_DEF.txt",
		Encoding:    &charmap.ISO8859_3,
		Description: "Nutrient Definition",
		Extractor:   etl.NutrientDefinitionExtractor{},
		h:           persistNutrientDefinitions,
	},
	file{
		Name:        "NUT_DATA.txt",
		Description: "Nutrient Data",
		Extractor:   etl.NutrientExtractor{},
		h:           persistNutrient,
	},
	file{
		Name:        "DERIV_CD.txt",
		Description: "Data Derivation Code Description",
	},
	file{
		Name:        "SRC_CD.txt",
		Description: "Source Code",
	},
	file{
		Name:        "WEIGHT.txt",
		Description: "Weight",
		Extractor:   etl.WeightExtractor{},
		h:           persistWeight,
	},
	file{
		Name:        "DATA_SRC.txt",
		Description: "Sources of Data",
	},
	file{
		Name:        "DATSRCLN.txt",
		Description: "Sources of Data Link",
	},
	file{
		Name:        "FOOTNOTE.txt",
		Description: "Footnote",
	},
}

func runTx(tx *sql.DB, sql string) {
	if _, err := tx.Exec(sql); err != nil {
		log.Fatalf("Error executing '%v': %v", sql, err)
	}
}

var db *sql.DB

func runMigrations(db *sql.DB) {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	if n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		log.Printf("unable to migrate: %v", err)
	} else {
		log.Printf("Applied %d migrations!\n", n)
	}
}

func main() {
	flag.Parse()

	if *dataDirectory == "" {
		fmt.Printf("missing data.dir argument!\n")
		os.Exit(1)
	}
	if *database == "" {
		fmt.Printf("missing database argument!\n")
		os.Exit(1)
	}

	if _, err := os.Stat(*database); err == nil {
		os.Remove(*database)
	}

	var err error
	db, err = sql.Open("sqlite3", *database)
	db.SetMaxOpenConns(10)
	if err != nil {
		fmt.Printf("Failed opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	runMigrations(db)

	// health check
	for _, file := range files {
		var f = path.Join(*dataDirectory, "/", file.Name)
		var _, err = os.Stat(f)
		if err != nil {
			fmt.Printf("missing required file: %v\n", f)
			os.Exit(2)
		}
	}

	for _, fi := range files {
		if fi.Extractor != nil {

			var f, _ = os.Open(path.Join(*dataDirectory, "/", fi.Name))

			var r io.Reader = f
			if fi.Encoding != nil {
				r = transform.NewReader(r, (*fi.Encoding).NewDecoder())
			}

			var ch = make(chan interface{})
			go func() {
				if err := fi.Extractor.Parse(r, ch); err != nil {
					fmt.Printf("Failed to extract %v: %v\n", fi.Description, err)
				}
				close(ch)
			}()

			for i := range ch {
				fi.h(i, db)
			}

		}
	}

}
