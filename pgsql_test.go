package sample_test

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sawadashota/sql-test-sample"
)

type migrationSources struct {
	migrations []migrate.MigrationSource
}

func (mss *migrationSources) FindMigrations() ([]*migrate.Migration, error) {
	var migrations []*migrate.Migration

	for _, ms := range mss.migrations {
		migration, err := ms.FindMigrations()
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration...)
	}

	return migrations, nil
}

func migrateUp(t *testing.T) {
	migrateAndSeed(t, migrate.Up)
}

func migrateDown(t *testing.T) {
	migrateAndSeed(t, migrate.Down)
}

func migrateAndSeed(t *testing.T, dir migrate.MigrationDirection) {
	t.Helper()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mc := &migrationSources{
		migrations: []migrate.MigrationSource{
			&migrate.FileMigrationSource{
				Dir: "sql/migrations",
			},
			&migrate.FileMigrationSource{
				Dir: "sql/testdata",
			},
		},
	}

	_, err = migrate.Exec(db, "postgres", mc, dir)
	if err != nil {
		t.Fatal(err)
	}
}

func database(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestGetUser(t *testing.T) {
	sample.DB = database(t)
	defer sample.DB.Close()

	type args struct {
		id int
	}
	cases := map[string]struct {
		args    args
		want    *sample.User
		wantErr bool
	}{
		"find exists user": {
			args: args{
				id: 1,
			},
			want: &sample.User{
				Id:   1,
				Name: "Bob",
				Sex:  "male",
			},
			wantErr: false,
		},
		"find un exists user": {
			args: args{
				id: 9999999,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			migrateUp(t)
			defer migrateDown(t)

			got, err := sample.GetUser(c.args.id)
			if (err != nil) != c.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("GetUser() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	sample.DB = database(t)
	defer sample.DB.Close()

	type args struct {
		user *sample.User
	}
	cases := map[string]struct {
		args    args
		want    *sample.User
		wantErr bool
	}{
		"normal": {
			args: args{
				user: &sample.User{
					Name: "Mary",
					Sex:  "female",
				},
			},
			want: &sample.User{
				Id:   2,
				Name: "Mary",
				Sex:  "female",
			},
			wantErr: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			migrateUp(t)
			defer migrateDown(t)

			if err := sample.InsertUser(c.args.user); (err != nil) != c.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, c.wantErr)
			}
			if c.wantErr {
				return
			}

			u, err := sample.GetUser(c.want.Id)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(u, c.want) {
				t.Errorf("GetUser() = %v, want %v", u, c.want)
			}
		})
	}
}
