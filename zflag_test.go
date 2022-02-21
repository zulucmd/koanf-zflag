package kozflag_test

import (
	"reflect"
	"strings"
	"testing"

	kozflag "github.com/gowarden/koanf-zflag"
	"github.com/gowarden/zflag"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
)

func posflagCallback(key string, value string) (string, interface{}) {
	return strings.ReplaceAll(key, "-", "_"), value
}

func TestLoad(t *testing.T) {
	assert := func(t *testing.T, k *koanf.Koanf) {
		assertEq(t, k.String("key.one-example"), "val1")
		assertEq(t, k.String("key.two_example"), "val2")
		assertEq(t, k.Strings("key.strings"), []string{"1", "2", "3"})
		assertEq(t, k.Int("key.int"), int(123))
		assertEq(t, k.Ints("key.ints"), []int{1, 2, 3})
		assertEq(t, k.Float64("key.float"), 123.123)
	}

	fs := &zflag.FlagSet{}
	fs.String("key.one-example", "val1", "")
	fs.String("key.two_example", "val2", "")
	fs.StringSlice("key.strings", []string{"1", "2", "3"}, "")
	fs.Int("key.int", 123, "")
	fs.IntSlice("key.ints", []int{1, 2, 3}, "")
	fs.Float64("key.float", 123.123, "")

	k := koanf.New(".")
	err := k.Load(kozflag.Provider(fs, ".", k), nil)
	assertNoErr(t, err)
	assert(t, k)

	// Test load with a custom flag callback.
	k = koanf.New(".")
	p := kozflag.ProviderWithFlag(fs, ".", k, func(f *zflag.Flag) (string, interface{}) {
		return f.Name, kozflag.FlagVal(fs, f)
	})
	err = k.Load(p, nil)
	assertNoErr(t, err)
	assert(t, k)

	// Test load with a custom key, val callback.
	k = koanf.New(".")
	p = kozflag.ProviderWithValue(fs, ".", k, func(key, val string) (string, interface{}) {
		if key == "key.float" {
			return "", val
		}
		return key, val
	})
	err = k.Load(p, nil)
	assertNoErr(t, err)

	assertEq(t, k.String("key.one-example"), "val1")
	assertEq(t, k.String("key.two_example"), "val2")
	assertEq(t, k.String("key.int"), "123")
	assertEq(t, k.String("key.ints"), "[1 2 3]")
	assertEq(t, k.String("key.float"), "")
}

func TestIssue90(t *testing.T) {
	exampleKeys := map[string]interface{}{
		"key.one_example": "a struct value",
		"key.two_example": "b struct value",
	}

	fs := &zflag.FlagSet{}
	fs.String("key.one-example", "a posflag value", "")
	fs.String("key.two_example", "a posflag value", "")

	k := koanf.New(".")

	err := k.Load(confmap.Provider(exampleKeys, "."), nil)
	assertNoErr(t, err)

	err = k.Load(kozflag.ProviderWithValue(fs, ".", k, posflagCallback), nil)
	assertNoErr(t, err)

	assertEq(t, exampleKeys, k.All())
}

func TestIssue100(t *testing.T) {
	var err error
	f := &zflag.FlagSet{}
	f.StringToString("string", map[string]string{"k": "v"}, "")
	f.StringToInt("int", map[string]int{"k": 1}, "")
	f.StringToInt64("int64", map[string]int64{"k": 2}, "")

	k := koanf.New(".")

	err = k.Load(kozflag.Provider(f, ".", k), nil)

	assertNoErr(t, err)

	type Maps struct {
		String map[string]string
		Int    map[string]int
		Int64  map[string]int64
	}
	maps := new(Maps)

	err = k.Unmarshal("", maps)
	assertNoErr(t, err)

	assertEq(t, map[string]string{"k": "v"}, maps.String)
	assertEq(t, map[string]int{"k": 1}, maps.Int)
	assertEq(t, map[string]int64{"k": 2}, maps.Int64)
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`expected no error, got %q`, err)
	}
}

func assertEq(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(`expected %q, got %q`, expected, actual)
	}
}
