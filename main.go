package finitefilesystem

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"

	"github.com/mitchellh/hashstructure"
)

type Object map[string][]string
type Instance map[string]string

var types = make(map[string]map[string][]string)

func Register(name string, key string, generator []string) {
	if _, ok := types[name]; !ok { // value not already in memory map
		types[name] = make(map[string][]string)
	}

	if _, ok := types[name][key]; !ok { // key not already added (not dealing with multiple adds on same key, first register wins)
		types[name][key] = generator
	}
}

func Generate(name string) error {
	if _, ok := types[name]; !ok {
		return errors.New("No matching type registered")
	}

	object := types[name]

	if len(object) < 1 {
		return errors.New("type: " + name + " has no registered definitions, cannot generate")
	}

	// get sorted list of keys
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return generateLoop(name, keys, object, Instance{})
}

func generateLoop(name string, keys []string, object Object, instance Instance) error {
	for _, val := range object[keys[0]] {
		inst := Instance{}

		// Copy the object over to avoid corruption
		for key, val := range instance {
			inst[key] = val
		}

		// Set this key=val
		inst[keys[0]] = val

		if len(keys) > 1 { // cannot output yet, we must on!
			err := generateLoop(name, keys[1:], object, inst)

			if err != nil {
				return err
			}
		} else { // done with the permutations for this
			_, err := Store(inst)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Note: hashes will always point to the right thing but the object may be serialized differently in different runs
// this is apparently how gob works but it has no real ramifications other then potential confusion at seeing the data files change
func Store(instance Instance) (hash string, err error) {
	hash, err = hashInstance(instance)

	log.Print("hash: " + hash)

	if err != nil {
		return
	}

	if !fileExists(hash) {
		err = storeFile(hash, instance)
	}

	return
}

func Get(hash string) (string, error) {
	instance, err := getFile(hash)

	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func Remove(hash string) error {
	return os.Remove(getPath(hash))
}

// Internal functions
func hashInstance(instance Instance) (string, error) {
	i, err := hashstructure.Hash(instance, nil)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(i)), nil
}

func storeFile(hash string, instance Instance) error {
	// serialize instance and write it out under the hash
	p := getPath(hash)

	log.Print("Store at " + p)
	file, err := os.Create(p)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(instance)

		log.Print("Store success")
	}
	file.Close()
	return err
}

func getPath(hash string) string {
	return path.Join(getDirectory(), "data", hash)
}

func getFile(hash string) (Instance, error) {
	// lookup file using hash and unserialize it
	if !fileExists(hash) {
		return nil, errors.New("Could not find file")
	}

	instance := make(Instance)

	inst := &instance
	file, err := os.Open(getPath(hash))
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(inst)
	}
	file.Close()
	return *inst, err
}

func fileExists(hash string) bool {
	b, err := exists(getPath(hash))

	if err != nil {
		return false
	}

	return b
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getDirectory() string {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		fmt.Println("Could not figure out position of directory")
		os.Exit(1)
	}

	return path.Dir(filename)
}
