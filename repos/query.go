package repos

import (
	"encoding/json"

	"github.com/supabase-community/postgrest-go"
)

func execeuteSelect(query *postgrest.FilterBuilder, dest any) (int64, error) {
	data, count, err := query.Execute()
	if err != nil {
		return count, err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return count, err
	}

	return count, nil
}
