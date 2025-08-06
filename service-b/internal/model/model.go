package model

type ViaCEPResponse struct {
	CEP 			string `json:"cep"`
	Logradouro 		string `json:"logradouro"`
	Complemento 	string `json:"complemento"`
	Bairro 			string `json:"bairro"`
	Localidade 		string `json:"localidade"`
	UF 				string `json:"uf"`
	Erro 			string `json:"erro"` 			
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherResponse struct {
	City 			string `json:"city"`
	TempC 			float64 `json:"temp_c"`
	TempF 			float64 `json:"temp_f"`
	TempK 			float64 `json:"temp_k"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}