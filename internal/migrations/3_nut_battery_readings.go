package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		for _, name := range []string{"nut_devices", "nut_stats"} {
			collection, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				return err
			}
			collection.Fields.Add(
				&core.BoolField{Name: "battery_charge_reported"},
				&core.BoolField{Name: "battery_runtime_reported"},
			)
			if err := app.Save(collection); err != nil {
				return err
			}
			// Older agents omitted both unknown values and measured zeroes.
			// Only a positive historical value is evidence of a reading.
			_, err = app.DB().NewQuery(`
				UPDATE {{` + name + `}}
				SET battery_charge_reported = (battery_charge > 0),
				    battery_runtime_reported = (battery_runtime > 0)
				WHERE battery_charge > 0 OR battery_runtime > 0
			`).Execute()
			if err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error {
		for _, name := range []string{"nut_devices", "nut_stats"} {
			collection, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				return err
			}
			collection.Fields.RemoveByName("battery_charge_reported")
			collection.Fields.RemoveByName("battery_runtime_reported")
			if err := app.Save(collection); err != nil {
				return err
			}
		}
		return nil
	})
}
