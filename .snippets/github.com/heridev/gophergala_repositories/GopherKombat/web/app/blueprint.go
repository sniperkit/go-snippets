package app

import (
	"bufio"
	"github.com/gophergala/GopherKombat/common/user"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func BlueprintHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetCurrentUser(r)
	data := make(map[string]interface{})
	data["loggedIn"] = ok
	if ok {
		data["user"] = user
	}
	render(w, "blueprint", data)

}

func BlueprintSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user, ok := GetCurrentUser(r)
		data := make(map[string]interface{})
		if ok {
			r.ParseForm()
			code := r.PostFormValue("code")
			validated, message := validate(user, code)
			data["success"] = validated
			data["message"] = message
			if !validated {
				remove(user)
			}
		} else {
			data["success"] = ok
			data["message"] = "You are not logged in."
		}
		renderJson(w, r, data)
	} else {
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func BlueprintGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		user, ok := GetCurrentUser(r)
		data := make(map[string]interface{})
		if ok {
			data["success"] = ok
			code, err := get(user)
			data["code"] = code
			data["message"] = "Here is your blueprint. Maybe you should tweak it a little?"
			if err != nil {
				data["success"] = false
				data["message"] = "Looks like you dont have a blueprint yet. Here is one for example to get started."
			}
		} else {
			data["success"] = ok
			data["message"] = "You are not logged in."
		}
		renderJson(w, r, data)
	} else {
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

func validate(user *user.User, code string) (bool, string) {
	err := save(user, code)
	if err != nil {
		return false, "Failed to save file."
	}
	file := os.Getenv("GOPATH") + "/src/blueprints/" + user.Name + "/main.go"
	out, err := exec.Command("go", "build", file).CombinedOutput()
	if err != nil {
		//panic(err)
	}
	message := string(out)
	if message == "" {
		message = "All good! Your blueprint is ready to be tested in GopherKombat."
	}
	return true, message
}

func save(user *user.User, code string) error {
	dir := os.Getenv("GOPATH") + "/src/blueprints/" + user.Name

	err := os.Chdir(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
	}
	if err != nil {
		panic(err)
		return err
	}
	f, err := os.Create("main.go")
	if os.IsExist(err) {
		err = os.Remove("main.go")
		f, err = os.Create("main.go")
	}
	if err != nil {
		panic(err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(code)

	if err != nil {
		panic(err)
		return err
	}
	w.Flush()
	return nil
}

func get(user *user.User) (string, error) {
	file := os.Getenv("GOPATH") + "/src/blueprints/" + user.Name + "/main.go"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		example := os.Getenv("GOPATH") + "/src/blueprints/example.go"
		data, _ := ioutil.ReadFile(example)
		return string(data), err
	}
	return string(data), nil
}

func remove(user *user.User) {
	dir := os.Getenv("GOPATH") + "/src/blueprints/" + user.Name

	err := os.Chdir(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
	}
	if err != nil {
		panic(err)
	}
	err = os.Remove("main.go")
	if err != nil {
		panic(err)
	}
}
