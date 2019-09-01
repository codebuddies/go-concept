package db

import (
	"errors"
	"fmt"
	"time"
)

type (
	Resource struct {
		ID          int          `json:"id" db:"id"`
		Title       string       `json:"title" db:"title"`
		Description *string      `json:"description" db:"description"`
		URL         string       `json:"url" db:"url"`
		Referrer    *string      `json:"referrer" db:"referrer"`
		Credit      *string      `json:"credit" db:"credit"`
		PublishedAt *time.Time   `json:"published_at" db:"published_at"`
		CreatedAt   time.Time    `json:"created_at" db:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
		Type        ResourceType `json:"type" db:"type"`
		Tags        []*Tag       `json:"tags" db:"-"`
	}

	Tag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	ResourceType string
)

const (
	ResourceTypeVideo    ResourceType = "video"
	ResourceTypePodcast  ResourceType = "podcast"
	ResourceTypeTalk     ResourceType = "talk"
	ResourceTypeTutorial ResourceType = "tutorial"
	ResourceTypeCourse   ResourceType = "course"
	ResourceTypeBook     ResourceType = "book"
	ResourceTypeBlog     ResourceType = "blog"
)

func (r *Resource) Validate() error {
	if r.Type != ResourceTypeVideo &&
		r.Type != ResourceTypePodcast &&
		r.Type != ResourceTypeTalk &&
		r.Type != ResourceTypeTutorial &&
		r.Type != ResourceTypeCourse &&
		r.Type != ResourceTypeBook &&
		r.Type != ResourceTypeBlog {
		return fmt.Errorf("resource type must be one of (%s, %s, %s, %s, %s, %s, %s)",
			ResourceTypeVideo,
			ResourceTypePodcast,
			ResourceTypeTalk,
			ResourceTypeTutorial,
			ResourceTypeCourse,
			ResourceTypeBook,
			ResourceTypeBlog)
	}
	return nil
}

func (db *DB) GetResources(page, perPage int) ([]*Resource, error) {
	var resources []*Resource
	offset := (page - 1) * perPage

	if err := db.DB.Select(&resources, "SELECT * FROM resources LIMIT $1 OFFSET $2", perPage, offset); err != nil {
		return nil, err
	}

	return resources, nil
}

func (db *DB) InsertResources(resource *Resource) error {
	if resource.Title == "" {
		return errors.New("resource must have a valid title")
	}

	if resource.Description != nil && *resource.Description == "" {
		return errors.New("resource must have a valid description")
	}

	if resource.URL == "" {
		return errors.New("resource must have a valid url")
	}

	if resource.Referrer != nil && *resource.Referrer == "" {
		return errors.New("resource must have a valid referrer")
	}

	if resource.Credit != nil && *resource.Credit == "" {
		return errors.New("resource must have a valid credit")
	}

	if err := resource.Validate(); err != nil {
		return err
	}

	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO resources (title, description, url, referrer, credit, type) VALUES ($1, $2, $3, $4, $5, $6)`,
		resource.Title, resource.Description, resource.URL, resource.Referrer, resource.Credit, resource.Type); err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("sql error: %s, rollback error: %s", err, txErr)
		}

		return err
	}

	if err := tx.Get(resource, `SELECT * FROM resources WHERE id = last_insert_rowid()`); err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("sql error: %s, rollback error: %s", err, txErr)
		}

		return err
	}

	return tx.Commit()
}

func (db *DB) GetResource(id int) (*Resource, error) {
	var resource Resource

	if err := db.DB.Get(&resource, "SELECT * FROM resources WHERE id = $1", id); err != nil {
		return nil, err
	}

	return &resource, nil
}
