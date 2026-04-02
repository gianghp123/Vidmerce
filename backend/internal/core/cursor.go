package core

import (
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func EncodeCursor(lastKey map[string]types.AttributeValue) (string, error) {
	if len(lastKey) == 0 {
		return "", nil
	}

	m := make(map[string]string, len(lastKey))
	for k, v := range lastKey {
		if sv, ok := v.(*types.AttributeValueMemberS); ok {
			m[k] = sv.Value
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeCursor(cursor string) (map[string]types.AttributeValue, error) {
	if cursor == "" {
		return nil, nil
	}

	b, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	key := make(map[string]types.AttributeValue, len(m))
	for k, v := range m {
		key[k] = &types.AttributeValueMemberS{Value: v}
	}
	return key, nil
}
