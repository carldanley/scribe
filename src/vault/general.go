package vault

import (
	"log"
	"strings"
	"time"

	"github.com/carldanley/scribe/src/compendium"
	"github.com/carldanley/scribe/src/instruments"
)

func (v *Vault) RegisterTomesFromCompendium(c *compendium.Compendium) {
	if v.TomeCache == nil {
		v.TomeCache = map[*compendium.TomeSpec]*Tome{}
	}

	for tomeKey := range c.Tomes {
		tomeSpec := c.Tomes[tomeKey]

		// if a cache for this tome does not exist, create it
		if _, ok := v.TomeCache[&tomeSpec]; !ok {
			tome := &Tome{
				Spec:       &tomeSpec,
				Instrument: instruments.CreateInstrument(tomeSpec.Instrument, &tomeSpec),
				Secrets:    map[string]*Secret{},
			}

			for secretKey := range tomeSpec.Secrets {
				secretSpec := tomeSpec.Secrets[secretKey]
				secret := Secret{
					Spec:           &secretSpec,
					LastCacheCheck: 0,
				}

				if secretSpec.WatchInterval == 0 {
					secretSpec.WatchInterval = 5
				}

				if secretSpec.WatchForChanges == true {
					v.ShouldWatchForChanges = true
				}

				tome.Secrets[secretSpec.Path] = &secret
			}

			v.TomeCache[&tomeSpec] = tome
		}
	}
}

func (v *Vault) GetExpiredSecrets() *map[string][]*Tome {
	expiredSecrets := map[string][]*Tome{}
	currentTime := int32(time.Now().Unix())

	for _, tome := range v.TomeCache {
		for path, secret := range tome.Secrets {
			isExpired := false

			if secret.Cache == nil {
				isExpired = true
			} else if secret.Spec.WatchForChanges == true {
				if (currentTime - secret.LastCacheCheck) >= secret.Spec.WatchInterval {
					isExpired = true
				}
			}

			if isExpired == true {
				if _, ok := expiredSecrets[path]; !ok {
					expiredSecrets[path] = []*Tome{}
				}

				expiredSecrets[path] = append(expiredSecrets[path], tome)
			}
		}
	}

	return &expiredSecrets
}

func (v *Vault) FetchCacheForPath(path string) *map[string]interface{} {
	cache, err := v.GetClient().Logical().Read(path)

	if err != nil || cache == nil {
		log.Println("Skipping cache retrieval for path (error):", path, "...")
		log.Println(err)
		return &map[string]interface{}{}
	} else if cache == nil {
		log.Println("Skipping cache retrieval for path (no cache):", path, "...")
		return &map[string]interface{}{}
	}

	return &cache.Data
}

func (v *Vault) TranscribeSecretsForPath(path string, secrets *map[string]interface{}, tome *Tome) {
	for _, secret := range tome.Secrets {
		if secret.Spec.Path != path {
			continue
		}

		included := map[string]compendium.SecretField{}
		omitted := map[string]compendium.SecretField{}

		for _, field := range secret.Spec.Fields {
			if field.Omit == true {
				omitted[field.Name] = field
			} else {
				included[field.Name] = field
			}
		}

		newCache := map[string]string{}
		for key, value := range *secrets {
			if len(included) == 0 && len(omitted) == 0 {
				newCache[key] = value.(string)
			} else if len(omitted) > 0 {
				if _, ok := omitted[key]; !ok {
					newCache[key] = value.(string)
				}
			} else if len(included) > 0 {
				if field, ok := included[key]; ok {
					if field.MapTo != "" {
						key = field.MapTo
					}

					if value == "" {
						newCache[key] = field.DefaultValue
					}

					if field.ForceUpper == true {
						newCache[key] = strings.ToUpper(value.(string))
					} else if field.ForceLower == true {
						newCache[key] = strings.ToLower(value.(string))
					} else {
						newCache[key] = value.(string)
					}
				}
			}
		}

		secret.Cache = newCache
		secret.LastCacheCheck = int32(time.Now().Unix())
	}
}

func (v *Vault) GetAllTomeTranscriptions() *map[*Tome]map[string]string {
	transcriptions := map[*Tome]map[string]string{}

	for _, tome := range v.TomeCache {
		transcription := map[string]string{}

		for _, secret := range tome.Secrets {
			for key, value := range secret.Cache {
				transcription[key] = value
			}
		}

		transcriptions[tome] = transcription
	}

	return &transcriptions
}

func (v *Vault) TranscriptionsAreDifferent(existing map[string]string, new map[string]string) bool {
	// check new keys and updated values
	for key, value := range new {
		if old, ok := existing[key]; ok {
			if old != value {
				return true
			}
		} else {
			return true
		}
	}

	// check removed keys
	for key := range existing {
		if _, ok := new[key]; !ok {
			return true
		}
	}

	return false
}

func (v *Vault) Update() {
	// cache the existing transcriptions for every tome
	existingTranscriptions := v.GetAllTomeTranscriptions()

	// iterate through the expired secrets, fetching a new cache for each path
	for path, tomes := range *v.GetExpiredSecrets() {
		secrets := v.FetchCacheForPath(path)

		// make sure we actually retrieved secrets for this path
		if len(*secrets) == 0 {
			continue
		}

		// perform a transcription of secrets on each tome for this path
		for _, tome := range tomes {
			v.TranscribeSecretsForPath(path, secrets, tome)
		}
	}

	// identify if any transcriptions have changed; if so, use the tome's instrument
	// to write the new composition of secrets
	for tome, newTranscription := range *v.GetAllTomeTranscriptions() {
		if existingTranscription, ok := (*existingTranscriptions)[tome]; ok {
			if v.TranscriptionsAreDifferent(existingTranscription, newTranscription) == true {
				(*tome.Instrument).Write(newTranscription)
			}
		}
	}

	// check to see if any of the secrets need to be watched
	if v.ShouldWatchForChanges == true {
		time.Sleep(time.Second * 1)
		v.Update()
	}
}
