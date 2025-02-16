package core

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/rs/zerolog/log"
)

var SensorQ chan []byte = make(chan []byte)

func SensorWorker() {
	hourly, err := LoadHourlySensors(time.Now())
	if err != nil {
		log.Error().Err(err).Msg("failed to load counters")
	}

	for {
		payload := <-SensorQ

		data := make(map[string]Counters)
		err := json.Unmarshal(payload, &data)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal payload")
			continue
		}

		now := time.Now()
		hour := now.Format("2006010215")
		if hour != hourly.Hour {
			hourly.Closed = true
			err = hourly.Save()
			if err != nil {
				log.Error().Err(err).Msg("failed to dump counters")
			}
			hourly = NewHourlySensors(now)
			log.Info().Msg("counter hourly rotation")
		}
		for k, v := range data {
			values := hourly.Sensors[k]
			if values == nil {
				values = make(Counters)
			}
			for inK, inV := range v {
				values[inK] += inV
			}
			hourly.Sensors[k] = values
		}
		hourly.UpdatedAt = time.Now()
		err = hourly.Save()
		if err != nil {
			log.Error().Err(err).Msg("failed to dump sensors")
		}
	}
}

type Counters map[string]float64

type HourlySensors struct {
	Sensors   map[string]Counters `json:"sensors"`
	Hour      string              `json:"hour"`
	UpdatedAt time.Time           `json:"updated_at"`
	Closed    bool                `json:"closed"`
}

func NewHourlySensors(t time.Time) *HourlySensors {
	return &HourlySensors{
		Hour:    t.Format("2006010215"),
		Sensors: make(map[string]Counters),
	}
}

// Save saves a representation of v to the file at path.
func (hc *HourlySensors) Save() error {
	filename := "/tmp/sensors." + hc.Hour + ".json"
	filenameWriting := "/tmp/sensors." + hc.Hour + ".json.writing"

	f, err := os.Create(filenameWriting) //nolint:gosec
	if err != nil {
		return errors.Wrapf(err, "failed to create new file")
	}
	defer f.Close()
	b, err := json.Marshal(hc)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(b))
	if err != nil {
		return errors.Wrapf(err, "failed to write new file")
	}
	err = os.Rename(filenameWriting, filename)
	if err != nil {
		return errors.Wrapf(err, "failed to rename file")
	}

	return err
}

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func LoadHourlySensors(t time.Time) (*HourlySensors, error) {
	f, err := os.Open("/tmp/sensors." + t.Format("2006010215") + ".json")
	if os.IsNotExist(err) {
		return NewHourlySensors(t), nil
	} else if err != nil {
		return NewHourlySensors(t), err
	}
	defer f.Close()

	res := HourlySensors{}
	err = json.NewDecoder(f).Decode(&res)
	if err != nil {
		return NewHourlySensors(t), errors.Wrapf(err, "failed to unmarshal counter file")
	}

	log.Info().Int("len", len(res.Sensors)).Msgf("bc-sensors loaded from disk")

	return &res, nil
}
