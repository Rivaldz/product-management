package repo

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"technical_test/internal/entity"
)

type CompanyRepo struct {
	db *pgxpool.Pool
	sq squirrel.StatementBuilderType
}

func NewCompanyRepo(pg *pgxpool.Pool) *CompanyRepo {
	return &CompanyRepo{
		db: pg,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *CompanyRepo) GetList(ctx context.Context) ([]entity.Company, error) {
	sql, args, err := r.sq.Select("id", "name").
		From("companies").
		Where(squirrel.Eq{"status": "ACTIVE"}).
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CompanyRepo - GetList - r.sq.Select: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("CompanyRepo - GetList - r.db.Query: %w", err)
	}
	defer rows.Close()

	var companies []entity.Company
	for rows.Next() {
		var c entity.Company
		err = rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, fmt.Errorf("CompanyRepo - GetList - rows.Scan: %w", err)
		}
		companies = append(companies, c)
	}

	if companies == nil {
		companies = []entity.Company{}
	}

	return companies, nil
}
