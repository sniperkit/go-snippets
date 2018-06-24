package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/giagiannis/data-profiler/core"
)

// /datasets/
func controllerDatasetList(w http.ResponseWriter, r *http.Request) Model {
	return modelDatasetsList()
}

// /datasets/<id>/
func controllerDatasetView(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	m := modelDatasetGet(id)
	return m
}

// /datasets/<id>/delete
func controllerDatasetDelete(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	modelDatasetDelete(id)
	http.Redirect(w, r, "/datasets/", 307)
	return nil
}

// /download/
func controllerDownload(w http.ResponseWriter, r *http.Request) Model {
	fileType := r.URL.Query().Get("type")
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	var filePath string
	if fileType == "datafile" {
		m := modelDatasetGetInfo(id)
		if m != nil {
			filePath = m.Path + "/" + name
		}
	} else if fileType == "sm" {
		m := modelSimilarityMatrixGet(id)
		if m != nil {
			filePath = m.Path
		}
	} else if fileType == "coord" {
		m := modelCoordinatesGet(id)
		if m != nil {
			filePath = m.Path
		}
	} else if fileType == "operator" {
		m := modelOperatorGet(id)
		if m != nil {
			filePath = m.Path
		}
	} else if fileType == "scores" {
		m := modelOperatorGet(id)
		if m != nil {
			filePath = m.ScoresFile
		}
	} else if fileType == "samples" {
		m := modelDatasetModelGet(id)
		if m != nil {
			filePath = m.SamplesPath
		}
		log.Println(m)
	} else if fileType == "appx" {
		m := modelDatasetModelGet(id)
		if m != nil {
			filePath = m.AppxValuesPath
		}
	}

	if filePath != "" {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Disposition",
			"attachment; filename="+name)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Write(content)
	} else {
		w.WriteHeader(404)
	}

	return nil
}

// /tasks/
func controllerTasksList(w http.ResponseWriter, r *http.Request) Model {
	return TEngine.Tasks
}

// /sm/<id>/visual
func controllerSMVisual(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	m := modelSimilarityMatrixGet(id)
	ret := modelDatasetGetInfo(m.DatasetID)
	ret.Files = modelDatasetGetFiles(m.DatasetID)
	return ret
}

// /sm/<id>/csv/
func controllerSMtoCSV(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	m := modelSimilarityMatrixGet(id)
	sm := new(core.DatasetSimilarityMatrix)
	cnt, err := ioutil.ReadFile(m.Path)
	if err != nil {
		log.Println(err)
	}
	sm.Deserialize(cnt)
	w.Write([]byte("x,y,value\n"))
	files := modelDatasetGetFiles(m.DatasetID)
	for i := 0; i < sm.Capacity(); i++ {
		for j := 0; j < sm.Capacity(); j++ {
			w.Write([]byte(fmt.Sprintf("%s,%s,%.5f\n", files[i], files[j], sm.Get(i, j))))
		}
	}
	return nil
}

// /sm/<id>/delete
func controllerSMDelete(w http.ResponseWriter, r *http.Request) Model {
	datasetID := r.URL.Query().Get("datasetID")
	_, id, _ := parseURL(r.URL.Path)
	modelSimilarityMatrixDelete(id)
	http.Redirect(w, r, "/datasets/"+datasetID+"#sm", 307)
	return nil
}

// /operator/<id>/delete
func controllerOperatorDelete(w http.ResponseWriter, r *http.Request) Model {
	datasetID := r.URL.Query().Get("datasetID")
	_, id, _ := parseURL(r.URL.Path)
	modelOperatorDelete(id)
	http.Redirect(w, r, "/datasets/"+datasetID+"#operators", 307)
	return nil
}

// TASK BASED urls

// /datasets/<id>/newsm
func controllerDatasetNewSM(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	action := r.URL.Query().Get("action")
	if action != "submit" {
		m := modelDatasetGetInfo(id)
		scripts := make(map[string]string)
		for _, f := range Conf.Scripts.Analysis {
			scripts[f] = path.Base(f)
		}
		sms := modelSimilarityMatrixGetByDataset(id)
		matrices := make(map[string]string)
		for _, s := range sms {
			matrices[s.EstimatorPath] = s.Filename
		}

		return struct {
			ID       string
			Scripts  map[string]string
			Matrices map[string]string
		}{ID: m.ID, Scripts: scripts, Matrices: matrices}
	}
	conf := make(map[string]string)
	//	err := r.ParseForm()
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		log.Println(err)
	} else {
		f, h, err := r.FormFile("script")
		if err != nil {
			log.Println(err)
		}
		tempF, err := ioutil.TempFile("/tmp", "userscript")
		if err != nil {
			log.Println(err)
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
		}
		tempF.Write(buf)
		tempF.Close()
		os.Chmod(tempF.Name(), 0700)
		conf["script"] = tempF.Name()
		conf["script-name"] = h.Filename
	}

	for k, v := range r.PostForm {
		if len(v) > 0 {
			if _, ok := conf[k]; !ok {
				conf[k] = v[0]
			}
		} else {
			conf[k] = ""
		}
	}
	TEngine.Submit(NewSMComputationTask(id, conf))
	http.Redirect(w, r, "/tasks/", 307)
	return nil
}

// /mds/<id>/run
func controllerMDSRun(w http.ResponseWriter, r *http.Request) Model {
	datasetID := r.URL.Query().Get("datasetID")
	action := r.URL.Query().Get("action")
	_, id, _ := parseURL(r.URL.Path)
	if action != "submit" { // render the form
		return struct{ ID, DatasetID string }{id, datasetID}
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	conf := map[string]string{"k": r.PostFormValue("k")}
	task := NewMDSComputationTask(id, datasetID, conf)
	TEngine.Submit(task)
	http.Redirect(w, r, "/tasks/", 307)
	return nil

}

// /operator/<id>/run
func controllerOperatorRun(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	task := NewOperatorRunTask(id)
	TEngine.Submit(task)
	http.Redirect(w, r, "/tasks/", 307)
	return nil
}

// /coords/<matrixid>
func controllerCoordsView(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	return modelCoordinatesGetByMatrix(id)
}

// /coords/<id>/visual
func controllerCoordsVisual(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	m := modelCoordinatesGet(id)
	if m == nil {
		log.Println("Coordinates file not found")
		return nil
	}
	cnt, err := ioutil.ReadFile(m.Path)
	if err != nil {
		log.Println(err)
	}
	sm := m.SimilarityMatrix
	//	sm := modelSimilarityMatrixGet(m.matrixID)
	if sm == nil {
		log.Println("SM not found")
		return nil
	}
	datasetID := sm.DatasetID
	fileNames := ""
	files := modelDatasetGetFiles(datasetID)
	for i, n := range files {
		fileNames += n
		if i < len(files)-1 {
			fileNames += "\n"
		}
	}
	operators := modelOperatorGetByDataset(datasetID)
	return struct {
		Coordinates, Labels string
		Operators           []*ModelOperator
	}{
		Coordinates: string(cnt),
		Labels:      fileNames,
		Operators:   operators}
}

// /dataset/<id>/newop
func controllerDatasetNewOP(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	action := r.URL.Query().Get("action")
	if action != "submit" { // render stuff for the form
		return modelDatasetGetInfo(id)
	}
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		log.Println(err)
	}
	f, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	operatorDescription := r.PostForm["description"][0]
	operatorCnt, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	operatorFilename := header.Filename

	//	log.Println(operatorName, operatorDescription, len(operatorCnt), operatorFilename)
	modelOperatorInsert(id, operatorDescription, operatorFilename, operatorCnt)

	http.Redirect(w, r, "/datasets/"+id, 307)
	return nil
}

// /datasets/new/new
func controllerDatasetNew(w http.ResponseWriter, r *http.Request) Model {
	action := r.URL.Query().Get("action")
	if action != "submit" { // render stuff for the form
		return nil
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	datasetName := r.PostForm["name"][0]
	datasetDescription := r.PostForm["description"][0]
	datasetPath := r.PostForm["path"][0]
	id := modelDatasetInsert(datasetName, datasetDescription, datasetPath)
	http.Redirect(w, r, "/datasets/"+id, 307)
	return nil
}

// /modeling/<id>/comparison/
func controllerModelComparison(w http.ResponseWriter, r *http.Request) Model {
	r.ParseForm()
	xLabel := r.PostForm["xlabel"][0]
	yLabel := r.PostForm["ylabel"][0]
	result := make([]struct{ Key, Value string }, 0)
	for _, modelID := range r.PostForm["ids"] {
		mod := modelDatasetModelGet(modelID)
		var x string
		if xLabel == "SR" {
			x = fmt.Sprintf("%.2f", mod.SamplingRate)
		} else if xLabel == "k" {
			x = mod.Coordinates.K
		} else {
			x = mod.Coordinates.SimilarityMatrix.Configuration[xLabel]
		}
		result = append(result, struct{ Key, Value string }{x, mod.Errors[yLabel]})
	}
	return struct {
		Data           []struct{ Key, Value string }
		XLabel, YLabel string
	}{result, xLabel, yLabel}
}

// /modeling/new/new
func controllerModelNew(w http.ResponseWriter, r *http.Request) Model {
	action := r.URL.Query().Get("action")
	_, id, _ := parseURL(r.URL.Path)
	if action != "submit" {
		// just render form
		operators := modelOperatorGetByDataset(id)
		coordinates := modelCoordinatesGetByDataset(id)
		matrices := modelSimilarityMatrixGetByDataset(id)
		return struct {
			Operators   []*ModelOperator
			Coordinates []*ModelCoordinates
			Matrices    []*ModelSimilarityMatrix
			MLScripts   map[string]string
			DatasetID   string
		}{
			Operators:   operators,
			Coordinates: coordinates,
			MLScripts:   Conf.Scripts.ML,
			DatasetID:   id,
			Matrices:    matrices,
		}
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	operator := r.Form["operatorid"][0]
	mlScript := r.Form["script"][0]
	modelType := r.Form["modeltype"][0]
	matrixid := r.Form["matrixid"][0]
	k := r.Form["k"][0]
	dataset := r.Form["datasetid"][0]
	coordinates := r.Form["coordinatesid"][0]
	sr, err := strconv.ParseFloat(r.Form["sr"][0], 64)
	if err != nil {
		log.Println(err)
	}

	log.Println(operator, mlScript, modelType, matrixid, k, dataset, coordinates)
	TEngine.Submit(
		NewModelTrainTask(
			dataset, operator, sr,
			modelType,
			coordinates, mlScript,
			matrixid, k))
	http.Redirect(w, r, "/tasks/", 307)
	return nil
}

// /scores/<id>/text
func controllerScoresText(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	m := modelOperatorGet(id)
	if m == nil {
		return nil
	}
	cnt, err := ioutil.ReadFile(m.ScoresFile)
	if err != nil {
		log.Println(err)
	}
	s := core.NewDatasetScores()
	s.Deserialize(cnt)
	buffer := new(bytes.Buffer)
	for k, v := range s.Scores {
		buffer.WriteString(fmt.Sprintf("%s:%.5f\n", k, v))
	}
	w.Write(buffer.Bytes())
	return nil
}

// /modeling/<id>/visual
func controllerModelVisual(w http.ResponseWriter, r *http.Request) Model {
	_, id, _ := parseURL(r.URL.Path)
	m := modelDatasetModelGet(id)
	samples, _ := ioutil.ReadFile(m.SamplesPath)
	apprx, _ := ioutil.ReadFile(m.AppxValuesPath)
	coordinates, _ := ioutil.ReadFile(m.Coordinates.Path)
	files := modelDatasetGetFiles(m.Dataset.ID)
	fileStr := ""
	for _, f := range files {
		fileStr += f + "\n"
	}
	scoresID := ""
	if m.Operator.ScoresFile != "" {
		scoresID = m.Operator.ID
	}

	return struct {
		Labels             string
		Samples            string
		ApproximatedValues string
		Coordinates        string
		ScoresID           string
		Errors             map[string]string
	}{fileStr, string(samples), string(apprx), string(coordinates), scoresID, m.Errors}
}

func controllerModelDelete(w http.ResponseWriter, r *http.Request) Model {
	datasetID := r.URL.Query().Get("datasetID")
	_, id, _ := parseURL(r.URL.Path)
	modelDatasetModelDelete(id)
	http.Redirect(w, r, "/datasets/"+datasetID+"#model", 307)
	return nil
}
