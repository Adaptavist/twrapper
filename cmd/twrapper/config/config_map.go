package config

import "fmt"

type ConfigMap map[string]interface{}

func (b ConfigMap) Set(k string, v interface{}) {
	b[k] = v
}

func (b ConfigMap) Get(k string) (interface{}, error) {
	val, ok := b[k]
	if !ok {
		return nil, fmt.Errorf("key %s not found", k)
	}
	return val, nil
}

func (b ConfigMap) GetString(k string) (string, error) {
	val, err := b.Get(k)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func (b ConfigMap) GetBool(k string) (bool, error) {
	val, err := b.Get(k)
	if err != nil {
		return false, err
	}
	return val.(bool), nil
}

func (b ConfigMap) GetInt(k string) (int, error) {
	val, err := b.Get(k)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}
