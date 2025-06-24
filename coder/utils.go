package coder

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/xeipuuv/gojsonschema"
// )

// // --- Deep copy helpers for safe argument transformation ---

// // DeepCopyMap realiza una copia profunda de un map[string]interface{}
// func DeepCopyMap(src map[string]interface{}) map[string]interface{} {
// 	if src == nil {
// 		return nil
// 	}
// 	out := make(map[string]interface{}, len(src))
// 	for k, v := range src {
// 		switch vTyped := v.(type) {
// 		case map[string]interface{}:
// 			out[k] = DeepCopyMap(vTyped)
// 		case []interface{}:
// 			out[k] = DeepCopySlice(vTyped)
// 		default:
// 			out[k] = v
// 		}
// 	}
// 	return out
// }

// // DeepCopySlice realiza una copia profunda de un []interface{}
// func DeepCopySlice(src []interface{}) []interface{} {
// 	if src == nil {
// 		return nil
// 	}
// 	out := make([]interface{}, len(src))
// 	for i, v := range src {
// 		switch vTyped := v.(type) {
// 		case map[string]interface{}:
// 			out[i] = DeepCopyMap(vTyped)
// 		case []interface{}:
// 			out[i] = DeepCopySlice(vTyped)
// 		default:
// 			out[i] = v
// 		}
// 	}
// 	return out
// }

// // Merges src into dst (modifies dst)
// func deepMerge(dst, src map[string]interface{}) {
// 	for k, v := range src {
// 		if vmap, ok := v.(map[string]interface{}); ok {
// 			if dmap, ok := dst[k].(map[string]interface{}); ok {
// 				deepMerge(dmap, vmap)
// 				continue
// 			}
// 		}
// 		dst[k] = v
// 	}
// }

// // Helper to flatten map[string]interface{} for result
// func flattenMapValues(input map[string]interface{}) map[string]interface{} {
// 	out := make(map[string]interface{}, len(input))
// 	for k, v := range input {
// 		switch v.(type) {
// 		case map[string]interface{}, []interface{}:
// 			encoded, err := json.Marshal(v)
// 			if err == nil {
// 				out[k] = string(encoded)
// 			} else {
// 				out[k] = fmt.Sprintf("%v", v)
// 			}
// 		default:
// 			out[k] = v
// 		}
// 	}
// 	return out
// }

// func mapStringValues(input map[string]interface{}) map[string]interface{} {
// 	out := make(map[string]interface{}, len(input))
// 	for k, v := range input {
// 		switch val := v.(type) {
// 		case string:
// 			out[k] = val
// 		default:
// 			out[k] = fmt.Sprintf("%v", val)
// 		}
// 	}
// 	return out
// }

// func isPlanning() bool {
// 	return os.Getenv("TF_LOG") == "TRACE" && os.Getenv("TF_IN_AUTOMATION") == "1"
// }

// func setStoreAsJSONString(d Settable, store map[string]interface{}) error {
// 	if store == nil {
// 		return d.Set("store", "")
// 	}
// 	bytes, err := json.Marshal(store)
// 	if err != nil {
// 		return fmt.Errorf("could not marshal store as JSON: %w", err)
// 	}
// 	return d.Set("store", string(bytes))
// }

// // Settable is a minimal interface for Set usage (for easier testing/util usage)
// type Settable interface {
// 	Set(string, interface{}) error
// }

// func validateWithSchema(schema map[string]interface{}, doc interface{}, side string) error {
// 	if schema == nil || len(schema) == 0 {
// 		log.Printf("[INFO] Skipping schema validation for %s: no schema provided by Lambda.", side)
// 		return nil
// 	}
// 	log.Printf("[INFO] Validating %s with JSON schema...", side)
// 	schemaLoader := gojsonschema.NewGoLoader(schema)
// 	docLoader := gojsonschema.NewGoLoader(doc)
// 	result, err := gojsonschema.Validate(schemaLoader, docLoader)
// 	if err != nil {
// 		log.Printf("[ERROR] JSON schema validation error (%s): %v", side, err)
// 		return fmt.Errorf("jsonschema validation error (%s): %w", side, err)
// 	}
// 	if !result.Valid() {
// 		msg := fmt.Sprintf("%s failed schema validation:\n", side)
// 		for _, desc := range result.Errors() {
// 			msg += "- " + desc.String() + "\n"
// 		}
// 		log.Printf("[ERROR] %s", msg)
// 		return fmt.Errorf(msg)
// 	}
// 	log.Printf("[INFO] %s passed JSON schema validation.", side)
// 	return nil
// }

// // --- Helper: parse JSON string (args) into map ---
// func parseArgsJSON(argsInput interface{}) (map[string]interface{}, error) {
// 	switch arr := argsInput.(type) {
// 	case []interface{}: // expecting an array of strings
// 		final := map[string]interface{}{}
// 		for i, el := range arr {
// 			s, ok := el.(string)
// 			if !ok {
// 				log.Printf("[WARN] Args index %d is not a string: %v", i, el)
// 				continue
// 			}
// 			var m map[string]interface{}
// 			if err := json.Unmarshal([]byte(s), &m); err != nil {
// 				log.Printf("[WARN] Failed to parse args[%d]: %v\nInput: %s", i, err, s)
// 				continue
// 			}
// 			deepMerge(final, m)
// 		}
// 		return final, nil
// 	case string:
// 		// Fallback for old usage: single string
// 		var m map[string]interface{}
// 		if err := json.Unmarshal([]byte(arr), &m); err != nil {
// 			return nil, fmt.Errorf("failed to parse args as JSON object: %w\nInput was:\n%s", err, arr)
// 		}
// 		return m, nil
// 	default:
// 		return nil, fmt.Errorf("args must be either string or array of strings")
// 	}
// }

// // --- NUEVAS FUNCIONES PARA TRANSFORMAR ARGUMENTOS ANTES DE ENVIAR A LAMBDA ---

// // Flattens block_devices list to map keyed by "name".
// // Recursively flattens coder, workspace, owner blocks (takes first element of each list).
// // Ahora NO muta los args originales: siempre trabaja sobre una copia profunda.
// func transformArgsForLambda(args map[string]interface{}) map[string]interface{} {
// 	argsCopy := DeepCopyMap(args)
// 	out := make(map[string]interface{})
// 	for k, v := range argsCopy {
// 		switch k {
// 		case "block_devices":
// 			// v is []interface{}
// 			deviceMap := make(map[string]interface{})
// 			if devices, ok := v.([]interface{}); ok {
// 				for _, device := range devices {
// 					if dmap, ok := device.(map[string]interface{}); ok {
// 						if name, ok := dmap["name"].(string); ok {
// 							// Optionally, remove "name" from device map
// 							dcopy := make(map[string]interface{})
// 							for dk, dv := range dmap {
// 								if dk != "name" {
// 									dcopy[dk] = dv
// 								}
// 							}
// 							deviceMap[name] = dcopy
// 						}
// 					}
// 				}
// 			}
// 			out[k] = deviceMap
// 		case "coder":
// 			// flatten first element
// 			out[k] = flattenFirstMap(v)
// 		default:
// 			out[k] = v
// 		}
// 	}
// 	// Flatten nested coder.workspace etc.
// 	if coder, ok := out["coder"].(map[string]interface{}); ok {
// 		if ws, ok := coder["workspace"]; ok {
// 			coder["workspace"] = flattenFirstMap(ws)
// 			// flatten owner
// 			if workspace, ok := coder["workspace"].(map[string]interface{}); ok {
// 				if owner, ok := workspace["owner"]; ok {
// 					workspace["owner"] = flattenFirstMap(owner)
// 				}
// 			}
// 		}
// 	}
// 	return out
// }

// // Helper: If v is a list with at least one map element, return the first map element.
// // Otherwise, return v as is.
// func flattenFirstMap(v interface{}) interface{} {
// 	if lst, ok := v.([]interface{}); ok && len(lst) > 0 {
// 		if m, ok := lst[0].(map[string]interface{}); ok {
// 			return DeepCopyMap(m)
// 		}
// 	}
// 	return v
// }

// func mapLambdaResultToTerraformOutput(result map[string]interface{}) map[string]interface{} {
// 	out := map[string]interface{}{}

// 	// id directo
// 	if v, ok := result["id"]; ok {
// 		out["id"] = v
// 	}

// 	// instance: map a lista de 1 elemento
// 	if instRaw, ok := result["instance"].(map[string]interface{}); ok {
// 		out["instance"] = []interface{}{instRaw}
// 	}

// 	// volumes: map a lista de 1 elemento, y sus hijos igual
// 	if volumesRaw, ok := result["volumes"].(map[string]interface{}); ok {
// 		volumes := map[string]interface{}{}
// 		if rootRaw, ok := volumesRaw["root"].(map[string]interface{}); ok {
// 			volumes["root"] = []interface{}{rootRaw}
// 		}
// 		if wsRaw, ok := volumesRaw["workspace"].(map[string]interface{}); ok {
// 			volumes["workspace"] = []interface{}{wsRaw}
// 		}
// 		out["volumes"] = []interface{}{volumes}
// 	}

// 	return out
// }

