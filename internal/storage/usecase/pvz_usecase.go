package usecase

import (
	"database/sql"
	"fmt"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type PVZUsecase struct {
	pvzStorage       *storage.PVZPostgresStorage
	receptionStorage *storage.ReceptionPostgresStorage //
	productStorage   *storage.ProductPostgresStorage   //
}

func NewPVZUsecase(pvzStorage *storage.PVZPostgresStorage, receptionStorage *storage.ReceptionPostgresStorage, productStorage *storage.ProductPostgresStorage) *PVZUsecase {
	return &PVZUsecase{pvzStorage: pvzStorage, receptionStorage: receptionStorage, productStorage: productStorage}
}

func (p *PVZUsecase) CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error) {

	pvz, err := p.pvzStorage.GetPVZById(id)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("pvz exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check PVZ existence: %w", err)
	}

	pvz, err = p.pvzStorage.CreatePVZ(id, user_id, city, date)
	if err != nil {
		return nil, err
	}
	return pvz, nil
}

// TODO: плохо реализовал, нужно переделать
/* func (p *PVZUsecase) GetListPVZ(user_id uuid.UUID) ([]entity.ListPVZ, error) {
	var output []entity.ListPVZ

	listPVZ, err := p.pvzStorage.GetListPVZByUserId(user_id)
	if err != nil {
		return nil, err
	}
	for _, PVZ := range listPVZ {
		listReception, err := p.receptionStorage.GetListReceptionsByPVZId(PVZ.ID)
		if err != nil {
			return nil, err
		}
		for _, reception := range listReception {
			listProduct, err := p.productStorage.GetListProductsByReceptionId(reception.ID)
			if err != nil {
				return nil, err
			}
			reception.Products = append(reception.Products, listProduct...)
		}
		output = append(output, entity.ListPVZ{Pvz: PVZ, Receptions: listReception})
	}
	return output, nil
} */

//time.Now().UTC().Format(time.RFC3339)
