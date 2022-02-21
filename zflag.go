// Package kozflag implements a koanf.Provider that reads commandline
// parameters as conf maps using gowarden/zflag, a POSIX compliant
// alternative to Go's stdlib flag package.
package kozflag

import (
	"errors"

	"github.com/gowarden/zflag"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/maps"
)

// KZFlag implements a zflag command line provider.
type KZFlag struct {
	delim   string
	flagset *zflag.FlagSet
	ko      *koanf.Koanf
	cb      func(key string, value string) (string, interface{})
	flagCB  func(f *zflag.Flag) (string, interface{})
}

// Provider returns a commandline flags provider that returns
// a nested map[string]interface{} of environment variable where the
// nesting hierarchy of keys are defined by delim. For instance, the
// delim "." will convert the key `parent.child.key: 1`
// to `{parent: {child: {key: 1}}}`.
//
// It takes an optional (but recommended) Koanf instance to see if
// the flags defined have been set from other providers, for instance,
// a config file. If they are not, then the default values of the flags
// are merged. If they do exist, the flag values are not merged but only
// the values that have been explicitly set in the command line are merged.
func Provider(f *zflag.FlagSet, delim string, ko *koanf.Koanf) *KZFlag {
	return &KZFlag{
		flagset: f,
		delim:   delim,
		ko:      ko,
	}
}

// ProviderWithValue works exactly the same as Provider except the callback
// takes a (key, value) with the variable name and value and allows their modification.
// This is useful for cases where complex types like slices separated by
// custom separators.
// Returning "" for the key causes the particular flag to be disregarded.
func ProviderWithValue(f *zflag.FlagSet, delim string, ko *koanf.Koanf, cb func(key string, value string) (string, interface{})) *KZFlag {
	return &KZFlag{
		flagset: f,
		delim:   delim,
		ko:      ko,
		cb:      cb,
	}
}

// ProviderWithFlag takes zflag.FlagSet and a callback that takes *zflag.Flag
// and applies the callback to all items in the flagset. It does not parse
// zflag.Flag values and expects the callback to process the keys and values
// from *zflag.Flag however. FlagVal() can be used in the callback to avoid
// repeating the type-switch block for parsing values.
// Returning "" for the key causes the particular flag to be disregarded.
//
// Example:
//
//  p := koanf_zflag.ProviderWithFlag(flagset, ".", ko, func(f *zflag.Flag) (string, interface{})) {
//     // Transform the key in whatever manner.
//     key := f.Name
//
//     // Use FlagVal() and then transform the value, or don't use it at all
//     // and add custom logic to parse the value.
//     val := koanf_zflag.FlagVal(flagset, f)
//
//     return key, val
//  })
func ProviderWithFlag(f *zflag.FlagSet, delim string, ko *koanf.Koanf, cb func(f *zflag.Flag) (string, interface{})) *KZFlag {
	return &KZFlag{
		flagset: f,
		delim:   delim,
		ko:      ko,
		flagCB:  cb,
	}
}

// Read reads the flag variables and returns a nested conf map.
func (p *KZFlag) Read() (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	p.flagset.VisitAll(func(f *zflag.Flag) {
		var (
			key = f.Name
			val interface{}
		)

		// If there is a (key, value) callback set, pass the key and string
		// value from the flagset to it and use the results.
		if p.cb != nil {
			key, val = p.cb(key, f.Value.String())
		} else if p.flagCB != nil {
			// If there is a zflag.Flag callback set, pass the flag as-is
			// to it and use the results from the callback.
			key, val = p.flagCB(f)
		} else {
			// There are no callbacks set. Use the in-built flag value parser.
			val = FlagVal(p.flagset, f)
		}

		if key == "" {
			return
		}

		// If the default value of the flag was never changed by the user,
		// it should not override the value in the conf map (if it exists in the first place).
		if !f.Changed {
			if p.ko != nil {
				if p.ko.Exists(key) {
					return
				}
			} else {
				return
			}
		}

		// No callback. Use the key and value as-is.
		mp[key] = val
	})

	return maps.Unflatten(mp, p.delim), nil
}

// ReadBytes is not supported by the env koanf.
func (p *KZFlag) ReadBytes() ([]byte, error) {
	return nil, errors.New("zflag provider does not support this method")
}

// Watch is not supported.
func (p *KZFlag) Watch(cb func(event interface{}, err error)) error {
	return errors.New("koanf-zflag provider does not support this method")
}

// FlagVal examines a zflag.Flag and returns a typed value as an interface{}
// from the types that zflag supports. If it is of a type that isn't known
// for any reason, the value is returned as a string.
func FlagVal(fs *zflag.FlagSet, f *zflag.Flag) interface{} {
	var (
		val interface{}
	)

	if v, hasGetter := f.Value.(zflag.Getter); hasGetter {
		val = v.Get()
	}

	return val
}
