package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/giagiannis/data-profiler/core"
	_ "github.com/mattn/go-sqlite3"
)

// Models

// Model is the interface returned by all the model functions
type Model interface{}

// ModelDataset is the struct that represents a set of datasets
type ModelDataset struct {
	ID          string
	Name        string
	Path        string
	Description string
	Files       []string
	Operators   []*ModelOperator
	Matrices    []*ModelSimilarityMatrix
	Models      []*ModelDatasetModel
}

// ModelOperator is the struct that represents an operator for a given dataset
type ModelOperator struct {
	ID          string
	Path        string
	Description string
	Name        string
	DatasetID   string
	ScoresFile  string
}

// ModelSimilarityMatrix represents a similarity matrix
type ModelSimilarityMatrix struct {
	ID            string
	Path          string
	Filename      string
	Configuration map[string]string
	DatasetID     string
	EstimatorPath string
}

// ModelCoordinates represents a set of coordinates
type ModelCoordinates struct {
	ID               string
	Path             string
	Filename         string
	K                string
	GOF              string
	Stress           string
	SimilarityMatrix *ModelSimilarityMatrix
}

// ModelDatasetModel represents a model of an operator for a given stuff
type ModelDatasetModel struct {
	ID             string
	Coordinates    *ModelCoordinates
	Operator       *ModelOperator
	Dataset        *ModelDataset
	SamplingRate   float64
	Configuration  map[string]string
	Errors         map[string]string
	SamplesPath    string
	AppxValuesPath string
}

// FUNCTIONS

// dbConnect is responsible to establish the connection with the DB backend.
// Written as separate function to increase modularity between the different
// DB backends.
func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", Conf.Database)
	if err != nil {
		log.Println(err)
	}
	return db
}

func modelDatasetsList() []*ModelDataset {
	db := dbConnect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM datasets")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	result := make([]*ModelDataset, 0)
	for rows.Next() {
		obj := new(ModelDataset)
		rows.Scan(&obj.ID, &obj.Path, &obj.Name, &obj.Description)
		result = append(result, obj)
	}
	return result
}

func modelDatasetDelete(id string) *ModelDataset {
	m := modelDatasetGet(id)
	deleteByID("datasets", id)
	return m
}

func modelDatasetInsert(name, description, path string) string {
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO datasets(name,path,description) " +
			"VALUES(?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(name,
		path,
		description)
	if err != nil {
		log.Println(err)
	}
	resultInt, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d", resultInt)
}

func modelDatasetGetInfo(id string) *ModelDataset {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM datasets WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelDataset)
		rows.Scan(&obj.ID, &obj.Path, &obj.Name, &obj.Description)
		return obj
	}
	return nil
}

func modelDatasetGet(id string) *ModelDataset {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM datasets WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelDataset)
		rows.Scan(&obj.ID, &obj.Path, &obj.Name, &obj.Description)
		obj.Matrices = modelSimilarityMatrixGetByDataset(obj.ID)
		obj.Models = modelDatasetModelGetByDataset(obj.ID)
		obj.Operators = modelOperatorGetByDataset(obj.ID)
		obj.Files = modelDatasetGetFiles(obj.ID)
		return obj
	}
	return nil
}
func modelDatasetGetFiles(id string) []string {
	var results []string
	m := modelDatasetGetInfo(id)
	if m == nil {
		return nil
	}
	path := m.Path
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	for _, f := range fs {
		if !f.IsDir() {
			results = append(results, f.Name())
		}
	}
	return results
}

func modelSimilarityMatrixGetByDataset(id string) []*ModelSimilarityMatrix {
	db := dbConnect()
	defer db.Close()
	var results []*ModelSimilarityMatrix
	rows, err := db.Query("SELECT id, path, filename,configuration,estimatorpath " +
		" FROM matrices WHERE datasetid == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		obj := new(ModelSimilarityMatrix)
		confString := ""
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename, &confString, &obj.EstimatorPath)
		obj.Configuration = stringToJSON(confString)
		results = append(results, obj)
	}
	return results
}

// modelSimilarityMatrixInsert inserts a new SM and returns the newly created Id
func modelSimilarityMatrixInsert(datasetID string, smBuffer, estBuffer []byte, conf map[string]string) *ModelSimilarityMatrix {
	dts := modelDatasetGetInfo(datasetID)
	smPath := writeBufferToFile(dts, "matrices", smBuffer)
	var estPath string
	if estBuffer != nil {
		estPath = writeBufferToFile(dts, "estimators", estBuffer)
	}
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO matrices(path,filename,configuration,datasetid,estimatorpath) " +
			"VALUES(?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(smPath,
		path.Base(smPath),
		jsonToString(conf),
		dts.ID, estPath)
	if err != nil {
		log.Println(err)
	}
	resultInt, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	id := fmt.Sprintf("%d", resultInt)
	return &ModelSimilarityMatrix{
		Configuration: conf,
		DatasetID:     datasetID,
		EstimatorPath: estPath,
		Path:          smPath,
		Filename:      path.Base(smPath),
		ID:            id,
	}
}

func modelSimilarityMatrixGet(id string) *ModelSimilarityMatrix {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id,path,filename,configuration,datasetid" +
		" FROM matrices WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelSimilarityMatrix)
		confString := ""
		rows.Scan(
			&obj.ID,
			&obj.Path,
			&obj.Filename,
			&confString,
			&obj.DatasetID)
		obj.Configuration = make(map[string]string)
		json.Unmarshal([]byte(confString), &obj.Configuration)
		return obj
	}
	return nil
}

func modelSimilarityMatrixDelete(id string) *ModelSimilarityMatrix {
	m := modelSimilarityMatrixGet(id)
	if m != nil {
		os.Remove(m.Path)
		os.Remove(m.EstimatorPath)
	}
	deleteByID("matrices", id)
	return nil
}

func modelCoordinatesGet(id string) *ModelCoordinates {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id, path, filename, k, gof, stress, matrixid" +
		" FROM coordinates WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelCoordinates)
		matrixID := ""
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename, &obj.K,
			&obj.GOF, &obj.Stress, &matrixID)
		obj.SimilarityMatrix = modelSimilarityMatrixGet(matrixID)
		return obj
	}
	return nil
}

func modelCoordinatesGetByDataset(datasetID string) []*ModelCoordinates {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT coordinates.*" +
		" FROM coordinates,matrices WHERE matrices.id == coordinates.matrixid AND datasetid == " + datasetID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	var result []*ModelCoordinates
	for rows.Next() {
		obj := new(ModelCoordinates)
		matrixID := ""
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename, &obj.K,
			&obj.GOF, &obj.Stress, &matrixID)
		obj.SimilarityMatrix = modelSimilarityMatrixGet(matrixID)
		result = append(result, obj)
	}
	return result
}

func modelCoordinatesGetByMatrix(matrixID string) []*ModelCoordinates {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT *" +
		" FROM coordinates WHERE matrixid == " + matrixID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	var result []*ModelCoordinates
	for rows.Next() {
		obj := new(ModelCoordinates)
		matrixID := ""
		rows.Scan(&obj.ID, &obj.Path, &obj.Filename, &obj.K,
			&obj.GOF, &obj.Stress, &matrixID)
		obj.SimilarityMatrix = modelSimilarityMatrixGet(matrixID)
		result = append(result, obj)
	}
	return result
}

func modelCoordinatesInsert(coordinates []core.DatasetCoordinates, datasetID, K, GOF, Stress, matrixID string) *ModelCoordinates {
	dts := modelDatasetGetInfo(datasetID)
	var coords [][]float64
	for _, c := range coordinates {
		coords = append(coords, c)
	}
	buffer := serializeCSVFile(coords)
	filePath := writeBufferToFile(dts, "coords", buffer)

	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO coordinates(path,filename,k,gof,stress,matrixid) " +
			"VALUES(?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(filePath,
		path.Base(filePath),
		K,
		GOF,
		Stress,
		matrixID)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func modelOperatorInsert(datasetID, description, filename string, content []byte) *ModelOperator {
	dts := modelDatasetGetInfo(datasetID)
	filePath := writeBufferToFile(dts, "operators", content)
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO operators(name,description,path,datasetid) " +
			"VALUES(?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(filename,
		description,
		filePath,
		datasetID)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func modelOperatorGet(id string) *ModelOperator {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, description, path, datasetid, scoresfile" +
		" FROM operators WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	if rows.Next() {
		obj := new(ModelOperator)
		rows.Scan(&obj.ID, &obj.Name, &obj.Description,
			&obj.Path, &obj.DatasetID, &obj.ScoresFile)
		return obj
	}
	return nil
}

func modelOperatorGetByDataset(id string) []*ModelOperator {
	db := dbConnect()
	defer db.Close()
	var results []*ModelOperator

	rows, err := db.Query("SELECT id, name, description, path, datasetid, scoresfile" +
		" FROM operators WHERE datasetid == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		obj := new(ModelOperator)
		rows.Scan(&obj.ID, &obj.Name, &obj.Description,
			&obj.Path, &obj.DatasetID, &obj.ScoresFile)
		results = append(results, obj)
	}
	return results
}

func modelOperatorDelete(id string) *ModelOperator {
	op := modelOperatorGet(id)
	if op != nil {
		os.Remove(op.Path)
		os.Remove(op.ScoresFile)
	}
	deleteByID("operators", id)
	return nil
}

func modelOperatorScoresInsert(operatorID string, content []byte) *ModelOperator {
	op := modelOperatorGet(operatorID)
	dts := modelDatasetGetInfo(op.DatasetID)
	filePath := writeBufferToFile(dts, "scores", content)

	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"UPDATE operators SET scoresfile=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(filePath,
		operatorID)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func modelDatasetModelInsert(
	coordinatesID, operatorID, datasetID string,
	samples, appxValues []byte,
	conf, errors map[string]string,
	samplingRate float64) *ModelDatasetModel {
	dts := modelDatasetGetInfo(datasetID)
	samplesPath := writeBufferToFile(dts, "samples", samples)
	appxValuesPath := writeBufferToFile(dts, "appx", appxValues)
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare(
		"INSERT INTO models(coordinatesid, operatorid, datasetid, samplingrate, " +
			"configuration, samplespath, appxvaluespath, errors) " +
			"VALUES(?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(
		coordinatesID,
		operatorID,
		datasetID,
		samplingRate,
		jsonToString(conf),
		samplesPath,
		appxValuesPath,
		jsonToString(errors),
	)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func modelDatasetModelDelete(id string) *ModelDatasetModel {
	m := modelDatasetModelGet(id)
	if m != nil {
		os.Remove(m.SamplesPath)
		os.Remove(m.AppxValuesPath)
	}
	deleteByID("models", id)
	return nil
}

func modelDatasetModelGet(id string) *ModelDatasetModel {
	db := dbConnect()
	defer db.Close()

	rows, err := db.Query(
		"SELECT id, coordinatesid, operatorid, datasetid, samplingrate, " +
			"configuration, samplespath, appxvaluespath, errors " +
			"FROM models WHERE id == " + id)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	confString, errorsString := "", ""
	if rows.Next() {
		obj := new(ModelDatasetModel)
		coordinatesID, operatorID, datasetID := "", "", ""
		rows.Scan(&obj.ID,
			&coordinatesID,
			&operatorID,
			&datasetID,
			&obj.SamplingRate,
			&confString,
			&obj.SamplesPath,
			&obj.AppxValuesPath,
			&errorsString)
		obj.Errors = stringToJSON(errorsString)
		obj.Configuration = stringToJSON(confString)
		obj.Coordinates = modelCoordinatesGet(coordinatesID)
		obj.Operator = modelOperatorGet(operatorID)
		obj.Dataset = modelDatasetGetInfo(datasetID)
		return obj
	}
	return nil
}

func modelDatasetModelGetByDataset(datasetID string) []*ModelDatasetModel {
	db := dbConnect()
	defer db.Close()
	var results []*ModelDatasetModel

	rows, err := db.Query("SELECT id, coordinatesid, operatorid, datasetid, samplingrate, " +
		"configuration, samplespath, appxvaluespath,errors " +
		"FROM models WHERE datasetid == " + datasetID)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	confString := ""
	for rows.Next() {
		obj := new(ModelDatasetModel)
		coordinatesID, operatorID, datasetID, errorsString := "", "", "", ""
		rows.Scan(&obj.ID,
			&coordinatesID,
			&operatorID,
			&datasetID,
			&obj.SamplingRate,
			&confString,
			&obj.SamplesPath,
			&obj.AppxValuesPath,
			&errorsString)
		obj.Configuration = stringToJSON(confString)
		obj.Coordinates = modelCoordinatesGet(coordinatesID)
		obj.Operator = modelOperatorGet(operatorID)
		obj.Dataset = modelDatasetGetInfo(datasetID)
		obj.Errors = stringToJSON(errorsString)
		results = append(results, obj)
	}
	return results
}

// utility functions
func writeBufferToFile(dts *ModelDataset, prefix string, buffer []byte) string {
	dstDir := dts.Path + "/" + prefix
	_, err := os.Stat(dstDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dstDir, 0777)
		if err != nil {
			log.Println(err)
		}
	}
	dstPath := dstDir + "/" + prefix + currentTimeSuffix()
	err = ioutil.WriteFile(dstPath, buffer, 0777)
	if err != nil {
		log.Println(err)
	}
	return dstPath
}

func currentTimeSuffix() string {
	t := time.Now()
	y, m, d := t.Year(), int(t.Month()), t.Day()
	h, min, sec := t.Hour(), t.Minute(), t.Second()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", y, m, d, h, min, sec)
}

func serializeCSVFile(coords [][]float64) []byte {
	buffer := new(bytes.Buffer)
	for i := range coords {
		for j := range coords[i] {
			buffer.WriteString(fmt.Sprintf("%.5f", coords[i][j]))
			if j < len(coords[i])-1 {
				buffer.WriteString(",")
			}
		}
		if i < len(coords)-1 {
			buffer.WriteString("\n")
		}
	}
	return buffer.Bytes()
}

func deserializeCSVFile(buffer []byte) [][]float64 {
	fileCont := string(buffer)
	lines := strings.Split(fileCont, "\n")
	var result [][]float64
	for _, line := range lines {
		var tuple []float64
		vals := strings.Split(line, ",")
		for _, v := range vals {
			floatVal, err := strconv.ParseFloat(v, 32)
			if err != nil {
				log.Println(err)
			}
			tuple = append(tuple, floatVal)
		}
		result = append(result, tuple)
	}
	return result
}

func deleteByID(matrix, id string) {
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM " + matrix + " WHERE id == ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err)
	}
}

func jsonToString(conf map[string]string) string {
	confString, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
	}
	return string(confString)
}

func stringToJSON(confString string) map[string]string {
	conf := make(map[string]string)
	err := json.Unmarshal([]byte(confString), &conf)
	if err != nil {
		log.Println(err)
	}
	return conf

}
