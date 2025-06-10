package database

import (
	"database/sql"
	"testing"

	"github.com/Torkel-Aannestad/coop-backend/services/social-media-aggregator-api/sql/migrations"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//The goal of these tests are to test the actual queries to the database with a test instance of the db running the same migrations.
//For every test we open a new connection to the test database and truncate the tables.

func setupTestDB(t *testing.T) *sql.DB {
	dsnTestDB := "postgres://postgres:postgres@localhost:5455/social_media_aggregator?sslmode=disable"
	db, err := sql.Open("postgres", dsnTestDB)
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// err = Migrate(db, "../../sql/migrations")
	err = MigrateFS(db, migrations.FS, ".")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE messages CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables %v", err)
	}

	return db
}

func TestInsertMessage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	messageModel := MessageModel{
		DB: db,
	}

	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{

		{
			name: "valid message",
			message: &Message{
				ExternalId: "100",
				Author:     "testUsername1",
				Body:       "This is a mastodon post",
				Platform:   "mastodon",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := messageModel.Insert(tt.message)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// By adding a messageModel.GetMessageByID() method we can test a lot of actual values being created or retrieved from the database. Same applies if we return the created values from messageModel.Insert(tt.message)
			// retrieved, err := messageModel.GetMessageByID(tt.message.ID)
			// require.NoError(t, err)

			// assert.Equal(t, tt.message.ID, retrieved.ID)
			// assert.Equal(t, tt.message.Body, tt.message.Body)

		})
	}
}
