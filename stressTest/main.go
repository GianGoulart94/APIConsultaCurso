package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"bufio"
)

type User struct{
	User	string	`json:"user"`
}


func main() {
	var listUser []User
	var resp []string

	arquivo, err := os.Open("users.csv")
	if err != nil {
		fmt.Println("[main] Houve um erro ao abrir o arquivo. Erro: ", err.Error())
		return
	}

	leitorCsv := csv.NewReader(arquivo)
	conteudo, err := leitorCsv.ReadAll()
	if err != nil {
		fmt.Println("[main] Houve um erro ao ler o arquivo com leitor CSV. Erro: ", err.Error())
		return
	}
	for i := range conteudo {
		listUser = append(listUser, User{User:strings.Join(conteudo[i],"")}) 
	}

	
	for i := range listUser {
		resp = append(resp,callApi(listUser[i].User)) 
	}

	criaRequest(resp)
	arquivo.Close()
}

func criaRequest(calls []string ) error{
	
	chamadas, err := os.Create("chamadas.txt")
	if err != nil {
		fmt.Println("[main] Houve um erro ao criar o arquivo TXT. Erro: ", err.Error())
		return err
	}
	defer chamadas.Close()

	escritor := bufio.NewWriter(chamadas)
	
	for _,call:= range calls{
		fmt.Fprintln(escritor,call)
	} 	

	return escritor.Flush()
}


func callApi(user string ) string{
	var url = "GET http://crossknowledge.magazineluiza.com.br/checkUser/"
	
	url = fmt.Sprintf(url+"%s/%s", user, "7137EA6C-C1B1-CA61-78D0-5444E8482E03");

		// Build the request
	
	return url	
}