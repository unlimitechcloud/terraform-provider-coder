package coder

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// // Helper para obtener los argumentos previos desde el estado del recurso.
// func getPreviousArgs(d *schema.ResourceData) map[string]interface{} {
// 	if d == nil || d.IsNewResource() {
// 		return map[string]interface{}{}
// 	}
// 	args := make(map[string]interface{})
// 	for key := range remoteResourceInputSchema() {
// 		if v, ok := d.GetOkExists(key); ok {
// 			args[key] = v
// 		}
// 	}
// 	return args
// }

// // Helper variante para ResourceDiff (CustomizeDiff)
// func getPreviousArgsDiff(d *schema.ResourceDiff) map[string]interface{} {
// 	args := make(map[string]interface{})
// 	for key := range remoteResourceInputSchema() {
// 		old, _ := d.GetChange(key)
// 		if old != nil {
// 			args[key] = old
// 		}
// 	}
// 	return args
// }

// // Settea los outputs explícitos definidos en outputs.go a partir del campo result de Lambda
// func setOutputFields(d *schema.ResourceData, result map[string]interface{}) error {
// 	tfOutput := mapLambdaResultToTerraformOutput(result)
// 	for key := range remoteResourceOutputSchema() {
// 		if v, ok := tfOutput[key]; ok {
// 			if err := d.Set(key, v); err != nil {
// 				return fmt.Errorf("failed to set output field '%s': %w", key, err)
// 			}
// 		}
// 	}
// 	return nil
// }

// // Merge de dos esquemas
// func mergeSchemas(src, dst map[string]*schema.Schema) map[string]*schema.Schema {
// 	out := map[string]*schema.Schema{}
// 	for k, v := range src {
// 		out[k] = v
// 	}
// 	for k, v := range dst {
// 		out[k] = v
// 	}
// 	return out
// }

// func resourceRemote() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceRemoteCreate,
// 		ReadContext:   resourceRemoteRead,
// 		UpdateContext: resourceRemoteUpdate,
// 		DeleteContext: resourceRemoteDelete,
// 		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
// 			client := meta.(*remoteClient)
// 			args := make(map[string]interface{})
// 			for key := range remoteResourceInputSchema() {
// 				if v, ok := d.GetOk(key); ok {
// 					args[key] = v
// 				}
// 			}
// 			argsTransformed := transformArgsForLambda(args)
// 			store := map[string]interface{}{}
// 			if v, ok := d.GetOk("store"); ok {
// 				storeStr, ok := v.(string)
// 				if ok && storeStr != "" {
// 					_ = json.Unmarshal([]byte(storeStr), &store)
// 				}
// 			}
// 			oldArgs := getPreviousArgsDiff(d)
// 			oldArgsTransformed := transformArgsForLambda(oldArgs)
// 			isCreate := len(oldArgs) == 0
// 			if isCreate {
// 				log.Printf("[INFO] Skipping diff call to Lambda: this is a create operation (no previous state).")
// 				return nil
// 			}
// 			res, err := invokeLambda(client, lambdaPayload{
// 				Action: "diff",
// 				Args:   argsTransformed,
// 				State:  DeepCopyMap(oldArgsTransformed),
// 				Store:  store,
// 			})
// 			if err != nil {
// 				return err
// 			}
// 			if res.Replace {
// 				log.Printf("[INFO] Lambda requested replace: %s", res.Reason)
// 				if err := d.ForceNew("name"); err != nil {
// 					return fmt.Errorf("failed to mark 'name' for replacement: %w", err)
// 				}
// 			}
// 			return nil
// 		},
// 		Schema: mergeSchemas(
// 			remoteResourceInputSchema(),
// 			mergeSchemas(
// 				remoteResourceOutputSchema(),
// 				remoteResourceInternalSchema(),
// 			),
// 		),
// 	}
// }

// func resourceRemoteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*remoteClient)
// 	args := make(map[string]interface{})
// 	for key := range remoteResourceInputSchema() {
// 		if v, ok := d.GetOk(key); ok {
// 			args[key] = v
// 		}
// 	}
// 	argsTransformed := transformArgsForLambda(args)
// 	store := map[string]interface{}{}
// 	if v, ok := d.GetOk("store"); ok {
// 		storeStr, ok := v.(string)
// 		if ok && storeStr != "" {
// 			_ = json.Unmarshal([]byte(storeStr), &store)
// 		}
// 	}
// 	res, err := invokeLambda(client, lambdaPayload{Action: "create", Args: argsTransformed, Store: store, Planning: isPlanning()})
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if res.ID == "" {
// 		return diag.FromErr(fmt.Errorf("lambda create response missing required 'id' field or returned empty id"))
// 	}
// 	d.SetId(res.ID)
// 	if err := setOutputFields(d, res.Result); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	// Set 'result' como lista de un elemento (estructura compatible con schema TypeList en internals.go)
// 	tfResult := mapLambdaResultToTerraformOutput(res.Result)
// 	if err := d.Set("result", []interface{}{tfResult}); err != nil {
// 		return diag.FromErr(fmt.Errorf("failed to set result: %w", err))
// 	}
// 	if err := setStoreAsJSONString(d, res.Store); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return nil
// }

// func resourceRemoteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*remoteClient)
// 	args := make(map[string]interface{})
// 	for key := range remoteResourceInputSchema() {
// 		if v, ok := d.GetOk(key); ok {
// 			args[key] = v
// 		}
// 	}
// 	argsTransformed := transformArgsForLambda(args)
// 	store := map[string]interface{}{}
// 	if v, ok := d.GetOk("store"); ok {
// 		storeStr, ok := v.(string)
// 		if ok && storeStr != "" {
// 			_ = json.Unmarshal([]byte(storeStr), &store)
// 		}
// 	}
// 	res, err := invokeLambda(client, lambdaPayload{Action: "read", Args: argsTransformed, Store: store, Planning: isPlanning()})
// 	if err != nil {
// 		log.Printf("[ERROR] remote read failed: %v", err)
// 		return diag.FromErr(fmt.Errorf("remote read failed: %w", err))
// 	}
// 	if res.ID == "" {
// 		log.Printf("[INFO] remote resource no longer exists, clearing ID")
// 		d.SetId("")
// 		return nil
// 	}
// 	d.SetId(res.ID)
// 	if err := setOutputFields(d, res.Result); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	tfResult := mapLambdaResultToTerraformOutput(res.Result)
// 	if err := d.Set("result", []interface{}{tfResult}); err != nil {
// 		return diag.FromErr(fmt.Errorf("failed to set result: %w", err))
// 	}
// 	if err := setStoreAsJSONString(d, res.Store); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return nil
// }

// func resourceRemoteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*remoteClient)
// 	args := make(map[string]interface{})
// 	for key := range remoteResourceInputSchema() {
// 		if v, ok := d.GetOk(key); ok {
// 			args[key] = v
// 		}
// 	}
// 	argsTransformed := transformArgsForLambda(args)
// 	store := map[string]interface{}{}
// 	if v, ok := d.GetOk("store"); ok {
// 		storeStr, ok := v.(string)
// 		if ok && storeStr != "" {
// 			_ = json.Unmarshal([]byte(storeStr), &store)
// 		}
// 	}
// 	res, err := invokeLambda(client, lambdaPayload{
// 		Action:   "update",
// 		Args:     argsTransformed,
// 		Store:    store,
// 		Planning: isPlanning(),
// 	})
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if res.ID == "" {
// 		return diag.FromErr(fmt.Errorf("lambda update response missing required 'id' field or returned empty id"))
// 	}

// 	d.SetId(res.ID)
// 	if err := setOutputFields(d, res.Result); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	tfResult := mapLambdaResultToTerraformOutput(res.Result)
// 	if err := d.Set("result", []interface{}{tfResult}); err != nil {
// 		return diag.FromErr(fmt.Errorf("failed to set result: %w", err))
// 	}
// 	if err := setStoreAsJSONString(d, res.Store); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return nil
// }

// func resourceRemoteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*remoteClient)
// 	args := make(map[string]interface{})
// 	for key := range remoteResourceInputSchema() {
// 		if v, ok := d.GetOk(key); ok {
// 			args[key] = v
// 		}
// 	}
// 	argsTransformed := transformArgsForLambda(args)
// 	store := map[string]interface{}{}
// 	if v, ok := d.GetOk("store"); ok {
// 		storeStr, ok := v.(string)
// 		if ok && storeStr != "" {
// 			_ = json.Unmarshal([]byte(storeStr), &store)
// 		}
// 	}
// 	res, err := invokeLambda(client, lambdaPayload{Action: "delete", Args: argsTransformed, Store: store, Planning: isPlanning()})
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if res.ID == "" {
// 		d.SetId("")
// 		return nil
// 	}
// 	d.SetId(res.ID)
// 	if err := setOutputFields(d, res.Result); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	tfResult := mapLambdaResultToTerraformOutput(res.Result)
// 	if err := d.Set("result", []interface{}{tfResult}); err != nil {
// 		return diag.FromErr(fmt.Errorf("failed to set result: %w", err))
// 	}
// 	if err := setStoreAsJSONString(d, res.Store); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return nil
// }