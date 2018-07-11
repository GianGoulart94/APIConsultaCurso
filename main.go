// Documentation API ConsultaCurso
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - Authorization:
//
//     SecurityDefinitions:
//     Authorization:
//          type: md5
//          name: Authorization
//          in: header
//
// swagger:meta
//go:generate swagger generate spec
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"strconv"
	"bitbucket.org/magazine-ondemand/trava-app-cdc/controllers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-oci8"
	_ "github.com/qodrorid/godaemon"
	"github.com/rs/cors"

)

type CheckUser struct{
	Status 			string 	`json:"status"`
	User			string	`json:"user"`
	Progress		int		`json:"progress"`
	GuidTraining 	string	`json:"guidTraining"`
	TrainingTitle	string	`json:"trainingTitle"`
	Message			string	`json:"message"`
}


func main(){

		
	rotas := mux.NewRouter().SkipClean(true)
// swagger:operation GET /consultaCurso/{login}/{curso}/{idcurso} statusCurso 
// ---
// summary: Informa se o usuário passado por parametro concluiu o curso pesquisado.
// description: Se não for encontrado nem o usuário ou o curso, Error Not Found (404) será retornado.
// parameters:
// - name: login
//   in: path
//   description: username do colaborador
//   type: string
//   required: true
//
// - name: curso
//   in: path
//   description: id do curso na base de dados 
//   type: string
//   required: true
//
// - name: idcurso
//   in: path
//   description: chave do curso na api da crossknowledge 
//   type: string
//   required: true
//
// responses:
//   "200":
//     "$ref": "#/responses/consultaCurso"
//   "404":
//     "$ref": "#/responses/notFound"
	rotas.HandleFunc("/consultaCurso/{login}/{curso}/{idcurso}", consultaCurso).Methods("GET")
	rotas.HandleFunc("/getToken/{app}", controllers.GenerateToken).Methods("GET")
	rotas.HandleFunc("/health", controllers.HealthCheck).Methods("GET")
    
	Port, _ := strconv.Atoi(os.Getenv("PORT"))

    if Port == 0 {
        Port = 3003
    }
	fmt.Println("Server running in port:", Port)
	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"},
	})
	
	http.ListenAndServe(fmt.Sprintf(":%d", Port), c.Handler(rotas))
}

func consultaCurso(w http.ResponseWriter, r *http.Request){

	token := r.Header.Get("Authorization");

	auth := os.Getenv("TOKEN_API");

	fmt.Println(auth)
	fmt.Println(token)

	if auth != token || token == "" {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}



	var usuario controllers.RetornoCurso
	var desbloqueado controllers.RetornoCurso
	vars := mux.Vars(r);

	login := vars["login"];
	curso := vars["curso"];
	idcurso := vars["idcurso"]

	
	usuario = controllers.ConsultaCurso(login, curso)
	fmt.Println(usuario)
	
	
	if usuario.User == "" || usuario.Conclusao == "incomplete"{
		
		desbloqueado = controllers.ConsultaDesbloqueado(login, curso)

		if desbloqueado.User != "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(desbloqueado);
			
		}else{

			cross := checkUser(login, idcurso, curso);
			if cross.StatusCode == 4{
				if usuario.Conclusao == "incomplete"{
				
					cross.Conclusao = "uncompleted";
					cross.Status = "bloqueado";
					cross.StatusCode = 1
				}
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cross);
		}
		
	}else{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuario)
	}
	
}


func checkUser(login string, idcurso string, curso string ) controllers.RetornoCurso {
	var url = os.Getenv("URL_CROSS")
	var token = os.Getenv("TOKEN_CROSS")
	var usuario controllers.RetornoCurso
	var checkUser CheckUser
	
	url = fmt.Sprintf(url+"%s/%s", login, idcurso);

		// Build the request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	fmt.Println(url)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	
	}
	
	defer resp.Body.Close()
	
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode);
		
	// Fill the record with the data from the JSON
	if resp.StatusCode == 200 || resp.StatusCode == 202 {
		jsonErr := json.Unmarshal(data, &checkUser)
		if jsonErr != nil {
			panic(jsonErr)
		}
		fmt.Println(checkUser)
			
		if checkUser.Status != "1" {
			usuario.User = login;
			usuario.IdCurso = curso;
			usuario.Conclusao = "uncompleted"
			usuario.Status	= "bloqueado"
			usuario.StatusCode = 1 
			
			return	usuario;
		}else{
			usuario.User = login;
			usuario.IdCurso = curso;
			usuario.Conclusao = "completed"
			usuario.Status	= "liberado"
			usuario.StatusCode = 0 
			return	usuario;
		}

	}else{
	
		usuario.User = login;
		usuario.IdCurso = curso;
		usuario.Conclusao = "not found"
		usuario.Status = "not found"
		usuario.StatusCode = 4
		return	usuario;
			
	}
	

}		

