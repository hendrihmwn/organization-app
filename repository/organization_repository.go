package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org OrganizationEntity) (OrganizationEntity, error)
	Update(ctx context.Context, org OrganizationEntity) (OrganizationEntity, error)
	Get(ctx context.Context, id int) (OrganizationEntity, error)
	GetByParentArr(ctx context.Context, id []int) ([]OrganizationEntity, error)
	Delete(ctx context.Context, id []int) error
	GetAll(ctx context.Context) ([]OrganizationEntity, error)
}

type OrganizationRepositoryImpl struct {
	Db *sql.DB
}

func (o *OrganizationRepositoryImpl) Create(ctx context.Context, input OrganizationEntity) (output OrganizationEntity, err error) {
	err = o.Db.QueryRowContext(ctx, "INSERT INTO organizations (name, parent_id, level) VALUES ($1, $2, $3) RETURNING id, name, parent_id, created_at, updated_at",
		input.Name, input.ParentId, input.Level,
	).Scan(&output.Id, &output.Name, &output.ParentId, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return
}

func (o *OrganizationRepositoryImpl) Update(ctx context.Context, input OrganizationEntity) (output OrganizationEntity, err error) {
	err = o.Db.QueryRowContext(ctx, "UPDATE organizations SET name = $1, parent_id = $2 WHERE id = $3 AND deleted_at is null RETURNING id, name, parent_id, created_at, updated_at",
		input.Name, input.ParentId, input.Id,
	).Scan(&output.Id, &output.Name, &output.ParentId, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return
}

func (o *OrganizationRepositoryImpl) Get(ctx context.Context, id int) (output OrganizationEntity, err error) {
	err = o.Db.QueryRowContext(ctx, "SELECT id, name, parent_id, level, created_at, updated_at FROM organizations WHERE id = $1 AND deleted_at is null",
		id,
	).Scan(&output.Id, &output.Name, &output.ParentId, &output.Level, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return

}

func (o *OrganizationRepositoryImpl) GetByParentArr(ctx context.Context, id []int) (output []OrganizationEntity, err error) {
	rows, err := o.Db.QueryContext(ctx, "SELECT id, name, parent_id, created_at, updated_at FROM organizations WHERE parent_id = ANY($1) AND deleted_at is null",
		pq.Array(id),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []OrganizationEntity

	// Iterate over the rows
	for rows.Next() {
		var organization OrganizationEntity
		if err := rows.Scan(&organization.Id, &organization.Name, &organization.ParentId, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
			return nil, err
		}
		organizations = append(organizations, organization)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizations, nil
}

func (o *OrganizationRepositoryImpl) Delete(ctx context.Context, id []int) (err error) {
	_, err = o.Db.QueryContext(ctx, "UPDATE organizations SET deleted_at = NOW() WHERE id = any($1)",
		pq.Array(id),
	)
	if err != nil {
		return
	}
	return
}

func (o *OrganizationRepositoryImpl) GetAll(ctx context.Context) (output []OrganizationEntity, err error) {
	rows, err := o.Db.QueryContext(ctx, "SELECT id, name, parent_id, level, created_at, updated_at FROM organizations WHERE deleted_at is null")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []OrganizationEntity

	// Iterate over the rows
	for rows.Next() {
		var organization OrganizationEntity
		if err := rows.Scan(&organization.Id, &organization.Name, &organization.ParentId, &organization.Level, &organization.CreatedAt, &organization.UpdatedAt); err != nil {
			return nil, err
		}
		organizations = append(organizations, organization)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return organizations, nil
}

func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &OrganizationRepositoryImpl{
		Db: db,
	}
}
