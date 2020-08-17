package esl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Headers holds information regarding given header
type Headers struct {
	header map[string]interface{}
	lock   *sync.RWMutex
}

// NewHeaders initialize the given headers
func NewHeaders() Headers {
	headers := Headers{
		header: make(map[string]interface{}),
		lock:   &sync.RWMutex{},
	}

	return headers
}

func (h Headers) String() string {
	h.lock.RLock()
	defer h.lock.RUnlock()

	var headers []string = h.Keys()

	var full []string

	for _, key := range headers {
		val := h.GetString(key)
		full = append(full, fmt.Sprintf("%s=%s", key, val))
	}

	return fmt.Sprintf("%s", strings.Join(full, " | "))
}

// Add a new header, or update an existed one
func (h *Headers) Add(key string, value interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.header[key] = value
}

// Get a header, if not found, return nil
func (h Headers) Get(key string) interface{} {
	h.lock.RLock()
	defer h.lock.RUnlock()

	value, found := h.header[key]
	if !found {
		return nil
	}
	return value
}

// GetString return string value from headers. If not found, returns empty string
func (h Headers) GetString(key string) string {
	value := h.Get(key)

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		return strconv.FormatInt(i, 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		return strconv.FormatUint(i, 10)
	case reflect.Float32, reflect.Float64:
		f := val.Float()
		s := strconv.FormatFloat(f, 'f', 3, 64)
		return s
	default:
		return ""
	}
}

// GetInt retun value as int, if empty, return 0
func (h Headers) GetInt(key string) int64 {
	value := h.Get(key)

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		s := val.String()
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return i
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := val.Int()
		return i
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := val.Uint()
		return int64(i)
	case reflect.Float32, reflect.Float64:
		f := val.Float()
		return int64(f)
	default:
		return 0
	}
}

// Remove a given key
func (h *Headers) Remove(key string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.header, key)
}

// Exists a given key at the headers
func (h *Headers) Exists(key string) bool {
	h.lock.RLock()
	defer h.lock.RUnlock()

	_, exists := h.header[key]
	return exists
}

// Keys returns a list of all existed keys available
func (h *Headers) Keys() []string {
	h.lock.RLock()
	defer h.lock.RUnlock()

	var keys []string

	for key := range h.header {
		keys = append(keys, key)
	}

	return keys
}
