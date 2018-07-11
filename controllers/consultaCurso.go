package controllers

import (
	"database/sql"
	"fmt"
//	"strconv"
	_ "github.com/mattn/go-oci8"
	"os"
	"strings"
)
// swagger:response retornoCurso
type RetornoCurso struct{
	// Dados da consulta
	// in: body
	User 		string 	`json:"user"`
	//Campo referente ao username do colaborador
	IdCurso		string 	`json:"idcurso"`
	//Id do curso pesquisado
	Conclusao	string 	`json:"conclusao"`
	//Status de conclusao
	Status		string 	`json:"status"`
	//Status para acesso ao aplicativo
	StatusCode	int 	`json:"statusCode"`
	// id do status

}



// Consulta curso consulta a tabela de colaboradores que finalizaram o curso no banco
func ConsultaCurso(login string, curso string) RetornoCurso {
	connect := os.Getenv("DB")
	var usuario RetornoCurso
	db, err:= sql.Open("oci8", connect);
	
	defer db.Close()

	login = strings.ToLower(login)
	
	query := fmt.Sprintf("SELECT USERID, IDCURSO, CONCLUSAO FROM MAG_T_CAD_FUNC_CURSOS WHERE USERID = '%s' AND IDCURSO = '%s'", login,curso);
	fmt.Println(query);

	row, err := db.Query(query);
	if err != nil {
		panic(err)
	}
	defer row.Close();

	
	for row.Next() {
		var user, idcurso, conclusao sql.NullString
		if err := row.Scan(&user, &idcurso, &conclusao); err != nil {
			panic(err)
		}
		usuario.User = user.String
		usuario.IdCurso = idcurso.String
		usuario.Conclusao = conclusao.String
		
		if usuario.Conclusao == "completed"  {
			usuario.Status = "liberado"
			usuario.StatusCode = 0
		}else{
			usuario.Status = "bloqueado"
			usuario.StatusCode = 1
		}

	}

	return usuario

}

//ConsultaDesbloqueado consulta se o mesmo está desbloqueado
func ConsultaDesbloqueado(login string, idcurso string) RetornoCurso{
	connect := os.Getenv("DB")
	var usuario RetornoCurso
	db, err:= sql.Open("oci8", connect);
	
	defer db.Close()

	login = strings.ToUpper(login)
	
	query := fmt.Sprintf("SELECT USERID FROM MAG_T_CAD_CURSO_CDC WHERE USERID = '%s'", login);
	fmt.Println(query);

	row, err := db.Query(query);
	if err != nil {
		panic(err)
	}
	defer row.Close();

	
	for row.Next() {
		var user sql.NullString
		if err := row.Scan(&user); err != nil {
			panic(err)
		}
		usuario.User = strings.ToLower(user.String)
		usuario.IdCurso = idcurso
		usuario.Conclusao = "uncompleted"
		usuario.Status = "desbloqueado manualmente"
		usuario.StatusCode = 3
		
	}

	return usuario

}