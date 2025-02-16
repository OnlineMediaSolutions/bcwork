package kvdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/rs/zerolog/log"
)

var (
	DB     = make(map[string]string)
	lock   sync.RWMutex
	Saving int
)

func Get(key string) string {
	lock.RLock()
	defer lock.RUnlock()

	return DB[key]
}

func Set(key string, value string) {
	lock.Lock()
	defer lock.Unlock()
	DB[key] = value
}

func Count() int {
	lock.RLock()
	defer lock.RUnlock()

	return len(DB)
}

func Scan() {
	lock.RLock()
	defer lock.RUnlock()
	for k := range DB {
		fmt.Println(k)
	}
}

func Load(filename string) error {
	readFile, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msgf("Load failed")

		return errors.Wrapf(err, "failed to open file")
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	defer readFile.Close()

	kCount := 0
	for fileScanner.Scan() {
		kCount++
		line := fileScanner.Text()
		toks := strings.Split(line, "  -->  ")
		if len(toks) != 2 {
			log.Error().Str("line", line).Msg("corrupted line")
		}
		Set(strings.TrimSpace(toks[0]), strings.TrimSpace(toks[1]))
		if kCount%10000 == 0 {
			fmt.Println(kCount)
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Error().Err(err).Msgf("something bad happened in the line %v: %v", kCount, err)
	}

	return nil
}

func Save(filename string) (int, error) {
	if Saving > 0 {
		return Saving, nil
	}
	Saving++
	lock.RLock()
	defer lock.RUnlock()
	start := time.Now().Unix()
	f, err := os.Create(filename)
	if err != nil {
		log.Error().Err(err).Msgf("Save failed")
		return Saving, errors.Wrapf(err, "failed to create file")
	}
	w := bufio.NewWriter(f)

	for k, v := range DB {
		line := fmt.Sprintf("%s  -->  %s\n", k, v)
		_, err = w.WriteString(line)
		fmt.Println(line)
		if err != nil {
			Saving = 0
			log.Error().Err(err).Msgf("Save failed")

			return Saving, errors.Wrapf(err, "failed to create file")
		}
		Saving++
	}
	// Very important to invoke after writing a large number of lines
	err = w.Flush()
	if err != nil {
		Saving = 0
		log.Error().Err(err).Msgf("Save failed")

		return Saving, errors.Wrapf(err, "failed to flush file")
	}

	Saving = 0
	log.Info().Msgf("Save completed after %d seconds", time.Now().Unix()-start)

	return 0, nil
}
