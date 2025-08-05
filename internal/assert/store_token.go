package assert

import (
    "encoding/json"
    "fmt"
)

const keyJSONPath = "json_path"

func StoreToken(body []byte, kw map[string]any) error {
    path, ok := kw[keyJSONPath].(string)
    if !ok {
        return fmt.Errorf("store_token: %s param missing or not a string", keyJSONPath)
    }

    var data map[string]any
    if err := json.Unmarshal(body, &data); err != nil {
        return fmt.Errorf("store_token: %w", err)
    }

    raw, ok := data[path]
    if !ok {
        return fmt.Errorf("store_token: field %s not present in body", path)
    }

    token, ok := raw.(string)
    if !ok {
        return fmt.Errorf("store_token: field %s is not a string", path)
    }

    if Ctx() != nil {
        Ctx().Set("token", token)
    }
    return nil
}

func init() {
    Register("store_token", StoreToken)
}
