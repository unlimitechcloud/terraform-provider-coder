package base

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RemoteResourceHandlers struct {
	PreCreate  func(ctx context.Context, d *schema.ResourceData, meta interface{}) error
	PostCreate func(ctx context.Context, d *schema.ResourceData, meta interface{}, resp *LambdaResponse) error
	PreUpdate  func(ctx context.Context, d *schema.ResourceData, meta interface{}) error
	PostUpdate func(ctx context.Context, d *schema.ResourceData, meta interface{}, resp *LambdaResponse) error
	PreDelete  func(ctx context.Context, d *schema.ResourceData, meta interface{}) error
	PostDelete func(ctx context.Context, d *schema.ResourceData, meta interface{}, resp *LambdaResponse) error
}

type LambdaResponse struct {
	ID      string                 `json:"id"`
	Result  map[string]interface{} `json:"result"`
	Store   map[string]interface{} `json:"store"`
	Replace bool                   `json:"replace"`
	Reason  string                 `json:"reason"`
	Attribute  string                 `json:"attribute"`
}

type LambdaPayload struct {
	Resource string                 `json:"resource"`
	Action   string                 `json:"action"`
	Args     map[string]interface{} `json:"args"`
	State    map[string]interface{} `json:"state,omitempty"`
	Store    map[string]interface{} `json:"store,omitempty"`
	Planning bool                   `json:"planning,omitempty"`
}

type RemoteClient struct {
	LambdaName string
	Svc        *lambda.Lambda
	once       sync.Once
}

func terraformError(msg string, args ...interface{}) {
	out := fmt.Sprintf("[ERROR] "+msg, args...)
	log.Print(out)
	fmt.Fprintln(os.Stderr, "Error:", fmt.Sprintf(msg, args...))
}

func InvokeLambda(client *RemoteClient, payload LambdaPayload) (*LambdaResponse, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		terraformError("Failed to marshal Lambda payload: %v", err)
		return nil, err
	}
	log.Printf("[INFO] invoking Lambda %s with payload: %s", client.LambdaName, string(bytes))
	resp, err := client.Svc.Invoke(&lambda.InvokeInput{
		FunctionName: &client.LambdaName,
		Payload:      bytes,
	})
	if err != nil {
		terraformError("Lambda invocation failed: %v", err)
		return nil, err
	}
	if resp.FunctionError != nil {
		terraformError("Lambda returned function error: %s", string(resp.Payload))
		return nil, fmt.Errorf("lambda error: %s", string(resp.Payload))
	}

	log.Printf("[INFO] lambda invocation timestamp: %s", time.Now().Format(time.RFC3339Nano))
	log.Printf("[INFO] lambda response raw: %s", string(resp.Payload))

	var out map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &out); err != nil {
		terraformError("Failed to unmarshal lambda response: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] lambda response as map: %#v", out)

	var resultVal map[string]interface{}
	if res, ok := out["result"]; ok {
		switch v := res.(type) {
		case map[string]interface{}:
			resultVal = v
		case string:
			if err := json.Unmarshal([]byte(v), &resultVal); err != nil {
				terraformError("Could not unmarshal result string: %v", err)
				resultVal = map[string]interface{}{}
			}
		default:
			terraformError("Unexpected result type: %T", v)
			resultVal = map[string]interface{}{}
		}
	} else {
		resultVal = map[string]interface{}{}
	}
	log.Printf("[DEBUG] lambda result extracted: %#v", resultVal)

	var storeVal map[string]interface{}
	if store, ok := out["store"]; ok {
		switch v := store.(type) {
		case map[string]interface{}:
			storeVal = v
		case string:
			_ = json.Unmarshal([]byte(v), &storeVal)
		default:
			storeVal = map[string]interface{}{}
		}
	}
	log.Printf("[DEBUG] lambda store extracted: %#v", storeVal)

	replace, _ := out["replace"].(bool)
	reason, _ := out["reason"].(string)
	attribute, _ := out["attribute"].(string)

	id := ""
	if resultVal != nil {
		if idRaw, ok := resultVal["id"]; ok {
			id, _ = idRaw.(string)
		}
	}
	log.Printf("[DEBUG] lambda id extracted: '%s'", id)

	return &LambdaResponse{
		ID:      id,
		Result:  resultVal,
		Store:   storeVal,
		Replace: replace,
		Reason:  reason,
		Attribute:  attribute,
	}, nil
}

func TransformArgsForLambda(args map[string]interface{}) map[string]interface{} {
	return args
}

func DeepCopyMap(m map[string]interface{}) map[string]interface{} {
	b, _ := json.Marshal(m)
	out := map[string]interface{}{}
	_ = json.Unmarshal(b, &out)
	return out
}

func getStoreFromResource(d interface{ GetOk(string) (interface{}, bool) }) map[string]interface{} {
	store := map[string]interface{}{}
	if v, ok := d.GetOk("store"); ok {
		storeStr, ok := v.(string)
		if ok && storeStr != "" {
			_ = json.Unmarshal([]byte(storeStr), &store)
		}
	}
	return store
}

func setStoreAsJSONString(d *schema.ResourceData, store map[string]interface{}) error {
	if store == nil {
		return d.Set("store", "")
	}
	bs, err := json.Marshal(store)
	if err != nil {
		return err
	}
	return d.Set("store", string(bs))
}

func MergeSchemas(src, dst map[string]*schema.Schema) map[string]*schema.Schema {
	out := map[string]*schema.Schema{}
	for k, v := range src {
		out[k] = v
	}
	for k, v := range dst {
		out[k] = v
	}
	return out
}

func getPreviousArgs(d *schema.ResourceData, inputSchema map[string]*schema.Schema) map[string]interface{} {
	if d == nil || d.IsNewResource() {
		return map[string]interface{}{}
	}
	args := make(map[string]interface{})
	for key := range inputSchema {
		if v, ok := d.GetOkExists(key); ok {
			args[key] = v
		}
	}
	return args
}

func getPreviousArgsDiff(d *schema.ResourceDiff, inputSchema map[string]*schema.Schema) map[string]interface{} {
	args := make(map[string]interface{})
	for key := range inputSchema {
		old, _ := d.GetChange(key)
		if old != nil {
			args[key] = old
		}
	}
	return args
}

func normalizeListsForInternalSchema(schemaMap map[string]*schema.Schema, data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := map[string]interface{}{}
		for key, val := range v {
			if sch, ok := schemaMap[key]; ok {
				if sch.Type == schema.TypeList {
					switch inner := val.(type) {
					case []interface{}:
						if resSchema, ok := sch.Elem.(*schema.Resource); ok {
							outList := make([]interface{}, len(inner))
							for i, item := range inner {
								outList[i] = normalizeListsForInternalSchema(resSchema.Schema, item)
							}
							result[key] = outList
						} else {
							result[key] = inner
						}
					case map[string]interface{}:
						if resSchema, ok := sch.Elem.(*schema.Resource); ok {
							result[key] = []interface{}{normalizeListsForInternalSchema(resSchema.Schema, inner)}
						} else {
							result[key] = []interface{}{inner}
						}
					case nil:
						result[key] = []interface{}{}
					default:
						result[key] = []interface{}{inner}
					}
				} else if sch.Type == schema.TypeMap {
					result[key] = val
				} else if sch.Type == schema.TypeSet {
					result[key] = val
				} else if resSchema, ok := sch.Elem.(*schema.Resource); ok && (sch.Type == schema.TypeSet || sch.Type == schema.TypeList) {
					result[key] = normalizeListsForInternalSchema(resSchema.Schema, val)
				} else {
					result[key] = val
				}
			} else {
				result[key] = val
			}
		}
		return result
	case []interface{}:
		for i, item := range v {
			v[i] = normalizeListsForInternalSchema(schemaMap, item)
		}
		return v
	default:
		return v
	}
}

func setInternalFields(d *schema.ResourceData, lambdaResp map[string]interface{}, internalSchema map[string]*schema.Schema) error {
	for key, sch := range internalSchema {
		val, ok := lambdaResp[key]
		if !ok {
			if sch.Type == schema.TypeList {
				if err := d.Set(key, []interface{}{}); err != nil {
					return fmt.Errorf("failed to set internal field '%s': %w", key, err)
				}
			}
			if sch.Type == schema.TypeString {
				if err := d.Set(key, ""); err != nil {
					return fmt.Errorf("failed to set internal field '%s' to empty string: %w", key, err)
				}
			}
			continue
		}
		if key == "result" && sch.Type == schema.TypeList {
			norm := normalizeListsForInternalSchema(sch.Elem.(*schema.Resource).Schema, val)
			switch v := norm.(type) {
			case nil:
				val = []interface{}{}
			case []interface{}:
				val = v
			case map[string]interface{}:
				val = []interface{}{v}
			default:
				val = []interface{}{v}
			}
		}
		if key == "store" && sch.Type == schema.TypeString {
			if str, ok := val.(string); ok {
				val = str
			} else if val == nil {
				val = ""
			} else if m, ok := val.(map[string]interface{}); ok {
				bs, err := json.Marshal(m)
				if err != nil {
					return fmt.Errorf("failed to marshal store: %w", err)
				}
				val = string(bs)
			} else {
				val = ""
			}
		}
		if err := d.Set(key, val); err != nil {
			return fmt.Errorf("failed to set internal field '%s': %w", key, err)
		}
	}
	return nil
}

func setOutputFields(d *schema.ResourceData, result map[string]interface{}, outputSchema map[string]*schema.Schema) error {
	if result == nil {
		result = map[string]interface{}{}
	}
	normalized := normalizeListsForInternalSchema(outputSchema, result)
	normMap, _ := normalized.(map[string]interface{})
	for key := range outputSchema {
		val, ok := normMap[key]
		if !ok {
			if outputSchema[key].Type == schema.TypeList {
				if err := d.Set(key, []interface{}{}); err != nil {
					return fmt.Errorf("failed to set output field '%s': %w", key, err)
				}
			}
			continue
		}
		if outputSchema[key].Type == schema.TypeList && val == nil {
			val = []interface{}{}
		}
		if err := d.Set(key, val); err != nil {
			return fmt.Errorf("failed to set output field '%s': %w", key, err)
		}
	}
	return nil
}

// Utilidad para transformar claves especificas de listas singleton a mapas
func flattenSingletonLists(data map[string]interface{}, singletonKeys []string) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range data {
		if contains(singletonKeys, k) {
			if list, ok := v.([]interface{}); ok && len(list) == 1 {
				// Recursividad si el valor también es map
				if subMap, ok := list[0].(map[string]interface{}); ok {
					out[k] = flattenSingletonLists(subMap, singletonKeys)
				} else {
					out[k] = list[0]
				}
				continue
			}
		}
		// Si es mapa, revisa sus hijos recursivamente
		if subMap, ok := v.(map[string]interface{}); ok {
			out[k] = flattenSingletonLists(subMap, singletonKeys)
		} else if list, ok := v.([]interface{}); ok && len(list) > 0 {
			var newList []interface{}
			for _, item := range list {
				if itemMap, ok := item.(map[string]interface{}); ok {
					newList = append(newList, flattenSingletonLists(itemMap, singletonKeys))
				} else {
					newList = append(newList, item)
				}
			}
			out[k] = newList
		} else {
			out[k] = v
		}
	}
	return out
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Devuelve true si todos los valores del mapa son zero value ("" o 0 o nil o [] o {} o false)
func isZeroArgs(args map[string]interface{}) bool {
	for _, v := range args {
		switch vv := v.(type) {
		case string:
			if vv != "" {
				return false
			}
		case int:
			if vv != 0 {
				return false
			}
		case []interface{}:
			if len(vv) != 0 {
				return false
			}
		case map[string]interface{}:
			if len(vv) != 0 {
				return false
			}
		case bool:
			if vv {
				return false
			}
		default:
			if v != nil {
				return false
			}
		}
	}
	return true
}

func NewRemoteResource(
	name string,
	inputSchema map[string]*schema.Schema,
	outputSchema map[string]*schema.Schema,
	internalSchema map[string]*schema.Schema,
	clientFactory func(interface{}) *RemoteClient,
	handlers *RemoteResourceHandlers,
) *schema.Resource {
	// Especifica los campos que deben ser des-listificados
	singletonKeys := []string{"coder", "workspace", "owner"}
	return &schema.Resource{
		CreateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			if handlers != nil && handlers.PreCreate != nil {
				if err := handlers.PreCreate(ctx, d, m); err != nil {
					return diag.FromErr(err)
				}
			}
			client := clientFactory(m)
			args := map[string]interface{}{}
			for key := range inputSchema {
				if v, ok := d.GetOk(key); ok {
					args[key] = v
				}
			}
			args = flattenSingletonLists(TransformArgsForLambda(args), singletonKeys)
			store := getStoreFromResource(d)
			res, err := InvokeLambda(client, LambdaPayload{
				Resource: name,
				Action:   "create",
				Args:     args,
				Store:    store,
				Planning: false,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			if res.ID == "" {
				return diag.FromErr(fmt.Errorf("lambda create response missing required 'id' field or returned empty id"))
			}
			d.SetId(res.ID)
			if err := setOutputFields(d, res.Result, outputSchema); err != nil {
				return diag.FromErr(err)
			}
			lambdaResp := map[string]interface{}{
				"result": res.Result,
				"id":     res.ID,
				"store":  res.Store,
			}
			if err := setInternalFields(d, lambdaResp, internalSchema); err != nil {
				return diag.FromErr(err)
			}
			if err := setStoreAsJSONString(d, res.Store); err != nil {
				return diag.FromErr(err)
			}
			if handlers != nil && handlers.PostCreate != nil {
				if err := handlers.PostCreate(ctx, d, m, res); err != nil {
					return diag.FromErr(err)
				}
			}
			return nil
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			client := clientFactory(m)
			args := map[string]interface{}{}
			for key := range inputSchema {
				if v, ok := d.GetOk(key); ok {
					args[key] = v
				}
			}
			args = flattenSingletonLists(TransformArgsForLambda(args), singletonKeys)
			store := getStoreFromResource(d)
			res, err := InvokeLambda(client, LambdaPayload{
				Resource: name,
				Action:   "read",
				Args:     args,
				Store:    store,
				Planning: false,
			})
			if err != nil {
				log.Printf("[ERROR] remote read failed: %v", err)
				return diag.FromErr(fmt.Errorf("remote read failed: %w", err))
			}
			if res.ID == "" {
				log.Printf("[INFO] remote resource no longer exists, clearing ID")
				d.SetId("")
				return nil
			}
			d.SetId(res.ID)
			if err := setOutputFields(d, res.Result, outputSchema); err != nil {
				return diag.FromErr(err)
			}
			lambdaResp := map[string]interface{}{
				"result": res.Result,
				"id":     res.ID,
				"store":  res.Store,
			}
			if err := setInternalFields(d, lambdaResp, internalSchema); err != nil {
				return diag.FromErr(err)
			}
			if err := setStoreAsJSONString(d, res.Store); err != nil {
				return diag.FromErr(err)
			}
			return nil
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			if handlers != nil && handlers.PreUpdate != nil {
				if err := handlers.PreUpdate(ctx, d, m); err != nil {
					return diag.FromErr(err)
				}
			}
			client := clientFactory(m)
			args := map[string]interface{}{}
			for key := range inputSchema {
				if v, ok := d.GetOk(key); ok {
					args[key] = v
				}
			}
			args = flattenSingletonLists(TransformArgsForLambda(args), singletonKeys)
			store := getStoreFromResource(d)
			res, err := InvokeLambda(client, LambdaPayload{
				Resource: name,
				Action:   "update",
				Args:     args,
				Store:    store,
				Planning: false,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			if res.ID == "" {
				return diag.FromErr(fmt.Errorf("lambda update response missing required 'id' field or returned empty id"))
			}
			d.SetId(res.ID)
			if err := setOutputFields(d, res.Result, outputSchema); err != nil {
				return diag.FromErr(err)
			}
			lambdaResp := map[string]interface{}{
				"result": res.Result,
				"id":     res.ID,
				"store":  res.Store,
			}
			if err := setInternalFields(d, lambdaResp, internalSchema); err != nil {
				return diag.FromErr(err)
			}
			if err := setStoreAsJSONString(d, res.Store); err != nil {
				return diag.FromErr(err)
			}
			if handlers != nil && handlers.PostUpdate != nil {
				if err := handlers.PostUpdate(ctx, d, m, res); err != nil {
					return diag.FromErr(err)
				}
			}
			return nil
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			if handlers != nil && handlers.PreDelete != nil {
				if err := handlers.PreDelete(ctx, d, m); err != nil {
					return diag.FromErr(err)
				}
			}
			client := clientFactory(m)
			args := map[string]interface{}{}
			for key := range inputSchema {
				if v, ok := d.GetOk(key); ok {
					args[key] = v
				}
			}
			args = flattenSingletonLists(TransformArgsForLambda(args), singletonKeys)
			store := getStoreFromResource(d)
			res, err := InvokeLambda(client, LambdaPayload{
				Resource: name,
				Action:   "delete",
				Args:     args,
				Store:    store,
				Planning: false,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			if res.ID == "" {
				d.SetId("")
				return nil
			}
			d.SetId(res.ID)
			if err := setOutputFields(d, res.Result, outputSchema); err != nil {
				return diag.FromErr(err)
			}
			lambdaResp := map[string]interface{}{
				"result": res.Result,
				"id":     res.ID,
				"store":  res.Store,
			}
			if err := setInternalFields(d, lambdaResp, internalSchema); err != nil {
				return diag.FromErr(err)
			}
			if err := setStoreAsJSONString(d, res.Store); err != nil {
				return diag.FromErr(err)
			}
			if handlers != nil && handlers.PostDelete != nil {
				if err := handlers.PostDelete(ctx, d, m, res); err != nil {
					return diag.FromErr(err)
				}
			}
			return nil
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			client := clientFactory(meta)
			args := map[string]interface{}{}
			for key := range inputSchema {
				if v, ok := d.GetOk(key); ok {
					args[key] = v
				}
			}
			args = flattenSingletonLists(TransformArgsForLambda(args), singletonKeys)
			store := getStoreFromResource(d)
			oldArgs := getPreviousArgsDiff(d, inputSchema)
			oldArgs = flattenSingletonLists(TransformArgsForLambda(oldArgs), singletonKeys)
			if len(oldArgs) == 0 || isZeroArgs(oldArgs) {
				log.Printf("[INFO] Skipping diff call to Lambda: this is a create operation (no previous state or all zero values).")
				return nil
			}
			res, err := InvokeLambda(client, LambdaPayload{
				Resource: name,
				Action:   "diff",
				Args:     args,
				State:    oldArgs,
				Store:    store,
			})
			if err != nil {
				return err
			}
			if res.Replace {
				log.Printf("[INFO] Lambda requested replace: %s -> %s", res.Attribute, res.Reason)
				// for key := range inputSchema {
					if err := d.ForceNew(res.Attribute); err != nil {
						return fmt.Errorf("failed to mark '%s' for replacement: %w", res.Attribute, err)
					}
				// }
			}
			return nil
		},
		Schema: MergeSchemas(
			inputSchema,
			MergeSchemas(
				outputSchema,
				internalSchema,
			),
		),
	}
}