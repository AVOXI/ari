package ari

import (
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/log15"
	"strings"
)

type JsonFormatter struct{}

func prepareForMarshal(x any) any {
	switch x := x.(type) {
	case nil:
		return nil
	case error:
		// NOTE: this is the same thing log15 does, and lets us log top level errors since
		// json.Marshal turns errors into a nice helpful "{}"
		// however, this does not handle errors inside structs/etc, and there does not
		// seem to be a way to handle this case without making a JSON marshaler ourselves
		return x.Error()
	// prioritize json.Marshaler over fmt.Stringer
	case json.Marshaler:
		return x
	// same for TextMarshaler (json.Marshal handles this case)
	case encoding.TextMarshaler:
		return x
	case fmt.Stringer:
		// using fmt instead of x.String() since fmt wraps panics
		return fmt.Sprint(x)
	default:
		return x
	}
}

func marshalValue(v any) (msg json.RawMessage, err error) {
	// it is possible for marshal to panic (e.g., a broken MarshalJSON implementation)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic while marshaling value: %v", r)
			// fmt already wraps panics
			msg, _ = json.Marshal(fmt.Sprintf("%#v", v))
		}
	}()

	msg, err = json.Marshal(prepareForMarshal(v))
	if err != nil {
		// use the GoString representation if the JSON representation fails, can represent things
		// like NaN or other non-JSON values
		msg, _ = json.Marshal(fmt.Sprintf("%#v", v))
	}
	return
}

func (*JsonFormatter) Format(r *log15.Record) []byte {
	props := map[string]json.RawMessage{}

	errors := []string{}

	if v, err := json.Marshal(r.Time); err != nil {
		errors = append(errors, fmt.Sprintf("error marshaling time: %v", err))
	} else {
		props[r.KeyNames.Time] = v
	}
	if v, err := json.Marshal(r.Lvl.String()); err != nil {
		errors = append(errors, fmt.Sprintf("unable to marshal log level into JSON: %v", err))
	} else {
		props[r.KeyNames.Lvl] = v
	}
	if v, err := json.Marshal(r.Msg); err != nil {
		errors = append(errors, fmt.Sprintf("unable to marshal log message into JSON: %v", err))
	} else {
		props[r.KeyNames.Msg] = v
	}

	for i := 0; i < len(r.Ctx); i += 2 {
		key := ""
		if k, ok := r.Ctx[i].(string); ok {
			key = k
		} else {
			errors = append(errors, fmt.Sprintf("key with index %v (%#v) is not a string key", i, r.Ctx[i]))
			key = fmt.Sprintf("invalid-key-%d", i)
		}

		value, err := marshalValue(r.Ctx[i+1])
		if err != nil {
			errors = append(errors, fmt.Sprintf("error marshaling value with index %v and name %#v: %v", i+1, key, err))
		}
		props[key] = value
	}

	if len(errors) > 0 {
		if v, err := json.Marshal(strings.Join(errors, "\n")); err != nil {
			panic("unable to marshal log errors into JSON: " + err.Error())
		} else {
			props["LOG15_ERROR"] = v
		}
	}

	b, err := json.Marshal(props)
	if err != nil {
		panic("unable to marshal log record into JSON: " + err.Error())
	}
	return append(b, '\n')
}
