package coder

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/aws/aws-sdk-go/service/lambda"
// )

// // remoteClient es el cliente que maneja la conexión y los datos para invocar Lambda
// type remoteClient struct {
// 	LambdaName string
// 	Svc        *lambda.Lambda
// 	once       sync.Once
// }

// // lambdaPayload es la estructura del payload enviado a Lambda
// type lambdaPayload struct {
// 	Action   string                 `json:"action"`
// 	Args     map[string]interface{} `json:"args"`
// 	State    map[string]interface{} `json:"state,omitempty"`
// 	Store    map[string]interface{} `json:"store,omitempty"`
// 	Planning bool                   `json:"planning,omitempty"`
// }

// // lambdaResponse es la estructura de la respuesta que devuelve Lambda
// type lambdaResponse struct {
// 	ID      string                 `json:"id"`
// 	Result  map[string]interface{} `json:"result"`
// 	Store   map[string]interface{} `json:"store"`
// 	Replace bool                   `json:"replace"`
// 	Reason  string                 `json:"reason"`
// }

// func terraformError(msg string, args ...interface{}) {
// 	// Terraform reconoce [ERROR] como error, pero también puedes imprimir con el prefijo "Error: "
// 	out := fmt.Sprintf("[ERROR] "+msg, args...)
// 	// Imprime como error estándar
// 	log.Print(out)
// 	// Además, imprime con el prefijo "Error: " para mayor claridad en el output de Terraform
// 	fmt.Fprintln(os.Stderr, "Error:", fmt.Sprintf(msg, args...))
// }

// // invokeLambda realiza la invocación a la función Lambda y procesa la respuesta
// func invokeLambda(client *remoteClient, payload lambdaPayload) (*lambdaResponse, error) {
// 	bytes, err := json.Marshal(payload)
// 	if err != nil {
// 		terraformError("Failed to marshal Lambda payload: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("[INFO] invoking Lambda %s with payload: %s", client.LambdaName, string(bytes))
// 	resp, err := client.Svc.Invoke(&lambda.InvokeInput{
// 		FunctionName: &client.LambdaName,
// 		Payload:      bytes,
// 	})
// 	if err != nil {
// 		terraformError("Lambda invocation failed: %v", err)
// 		return nil, err
// 	}
// 	if resp.FunctionError != nil {
// 		terraformError("Lambda returned function error: %s", string(resp.Payload))
// 		return nil, fmt.Errorf("lambda error: %s", string(resp.Payload))
// 	}

// 	log.Printf("[INFO] lambda invocation timestamp: %s", time.Now().Format(time.RFC3339Nano))
// 	log.Printf("[INFO] lambda response raw: %s", string(resp.Payload))

// 	var out map[string]interface{}
// 	if err := json.Unmarshal(resp.Payload, &out); err != nil {
// 		terraformError("Failed to unmarshal lambda response: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("[DEBUG] lambda response as map: %#v", out)

// 	var resultVal map[string]interface{}
// 	if res, ok := out["result"]; ok {
// 		switch v := res.(type) {
// 		case map[string]interface{}:
// 			resultVal = v
// 		case string:
// 			if err := json.Unmarshal([]byte(v), &resultVal); err != nil {
// 				terraformError("Could not unmarshal result string: %v", err)
// 				resultVal = map[string]interface{}{}
// 			}
// 		default:
// 			terraformError("Unexpected result type: %T", v)
// 			resultVal = map[string]interface{}{}
// 		}
// 	} else {
// 		resultVal = map[string]interface{}{}
// 	}
// 	log.Printf("[DEBUG] lambda result extracted: %#v", resultVal)

// 	var storeVal map[string]interface{}
// 	if store, ok := out["store"]; ok {
// 		switch v := store.(type) {
// 		case map[string]interface{}:
// 			storeVal = v
// 		case string:
// 			_ = json.Unmarshal([]byte(v), &storeVal)
// 		default:
// 			storeVal = map[string]interface{}{}
// 		}
// 	}
// 	log.Printf("[DEBUG] lambda store extracted: %#v", storeVal)

// 	replace, _ := out["replace"].(bool)
// 	reason, _ := out["reason"].(string)

// 	id := ""
// 	if resultVal != nil {
// 		if idRaw, ok := resultVal["id"]; ok {
// 			id, _ = idRaw.(string)
// 		}
// 	}
// 	log.Printf("[DEBUG] lambda id extracted: '%s'", id)

// 	return &lambdaResponse{
// 		ID:      id,
// 		Result:  resultVal,
// 		Store:   storeVal,
// 		Replace: replace,
// 		Reason:  reason,
// 	}, nil
// }