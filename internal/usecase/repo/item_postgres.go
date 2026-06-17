package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"technical_test/internal/entity"
)

type ItemRepo struct {
	db *pgxpool.Pool
	sq squirrel.StatementBuilderType
}

func NewItemRepo(pg *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{
		db: pg,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ItemRepo) Insert(ctx context.Context, item entity.Item) (entity.Item, error) {
	sql, args, err := r.sq.Insert("items").
		Columns("company_id", "code", "name", "type", "price", "category_name", "status").
		Values(item.CompanyID, item.Code, item.Name, item.Type, item.Price, item.CategoryName, item.Status).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()

	if err != nil {
		return entity.Item{}, fmt.Errorf("ItemRepo - Insert - r.sq.Insert: %w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return entity.Item{}, entity.ErrDuplicateCode
		}
		return entity.Item{}, fmt.Errorf("ItemRepo - Insert - r.db.QueryRow: %w", err)
	}

	return item, nil
}

func (r *ItemRepo) GetByFilter(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
	query := r.sq.Select("id", "company_id", "code", "name", "type", "price", "category_name", "status", "created_at", "updated_at").
		From("items").
		Where(squirrel.Eq{"company_id": companyID})

	if filter.Status != "" {
		query = query.Where(squirrel.Eq{"status": filter.Status})
	} else {
		query = query.Where(squirrel.NotEq{"status": entity.StatusArchived})
	}

	if filter.Type != "" {
		query = query.Where(squirrel.Eq{"type": filter.Type})
	}

	if filter.Keyword != "" {
		kw := "%" + filter.Keyword + "%"
		query = query.Where(squirrel.Or{
			squirrel.ILike{"code": kw},
			squirrel.ILike{"name": kw},
		})
	}

	sql, args, err := query.OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("ItemRepo - GetByFilter - query.ToSql: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ItemRepo - GetByFilter - r.db.Query: %w", err)
	}
	defer rows.Close()

	var items []entity.Item
	for rows.Next() {
		var i entity.Item
		err = rows.Scan(&i.ID, &i.CompanyID, &i.Code, &i.Name, &i.Type, &i.Price, &i.CategoryName, &i.Status, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ItemRepo - GetByFilter - rows.Scan: %w", err)
		}
		items = append(items, i)
	}

	if items == nil {
		items = []entity.Item{}
	}

	return items, nil
}

func (r *ItemRepo) GetByID(ctx context.Context, companyID, itemID string) (entity.Item, error) {
	sql, args, err := r.sq.Select("id", "company_id", "code", "name", "type", "price", "category_name", "status", "created_at", "updated_at").
		From("items").
		Where(squirrel.Eq{"id": itemID, "company_id": companyID}).
		ToSql()

	if err != nil {
		return entity.Item{}, fmt.Errorf("ItemRepo - GetByID - r.sq.Select: %w", err)
	}

	var i entity.Item
	err = r.db.QueryRow(ctx, sql, args...).Scan(&i.ID, &i.CompanyID, &i.Code, &i.Name, &i.Type, &i.Price, &i.CategoryName, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Item{}, entity.ErrItemNotFound
		}
		return entity.Item{}, fmt.Errorf("ItemRepo - GetByID - r.db.QueryRow: %w", err)
	}

	return i, nil
}

func (r *ItemRepo) Update(ctx context.Context, item entity.Item) (entity.Item, error) {
	sql, args, err := r.sq.Update("items").
		Set("code", item.Code).
		Set("name", item.Name).
		Set("type", item.Type).
		Set("price", item.Price).
		Set("category_name", item.CategoryName).
		Set("status", item.Status).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": item.ID, "company_id": item.CompanyID}).
		Suffix("RETURNING updated_at").
		ToSql()

	if err != nil {
		return entity.Item{}, fmt.Errorf("ItemRepo - Update - r.sq.Update: %w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&item.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.Item{}, entity.ErrDuplicateCode
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Item{}, entity.ErrItemNotFound
		}
		return entity.Item{}, fmt.Errorf("ItemRepo - Update - r.db.QueryRow: %w", err)
	}

	return item, nil
}

func (r *ItemRepo) Archive(ctx context.Context, companyID, itemID string) error {
	sql, args, err := r.sq.Update("items").
		Set("status", entity.StatusArchived).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": itemID, "company_id": companyID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("ItemRepo - Archive - r.sq.Update: %w", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ItemRepo - Archive - r.db.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entity.ErrItemNotFound
	}

	return nil
}
