package repository

import (
	"database/sql"
	"log"

	"github.com/danilomarques1/godemo/registerMerchant/api/model"
)

type MerchantSqlRepository struct {
	db *sql.DB
}

func NewMerchantSqlRepository(db *sql.DB) *MerchantSqlRepository {
	return &MerchantSqlRepository{db: db}
}

func (mr *MerchantSqlRepository) Save(merchant *model.Merchant) error {
	stmt, err := mr.db.Prepare("insert into merchant values($1, $2, $3, $4, $5);")
	if err != nil {
		log.Printf("Error preparing statement %v\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		merchant.MerchantId,
		merchant.MerchantName,
		merchant.MerchantAddress,
		merchant.SubAcquirerId,
		merchant.SubAcquirerName,
	)

	if err != nil {
		log.Printf("Error exec statement %v\n", err)
		return err
	}

	return nil
}

func (mr *MerchantSqlRepository) FindById(merchantId string) (*model.Merchant, error) {
	log.Printf("Find merchant by id = %v\n", merchantId)
	stmt, err := mr.db.Prepare(`
		select merchant_id, merchant_name, merchant_address, sub_acquirer_id, sub_acquirer_name
		from merchant
		where merchant_id = $1;
	`)
	if err != nil {
		log.Printf("Error preparing statement %v\n", err)
		return nil, err
	}
	defer stmt.Close()

	merchant := &model.Merchant{}
	err = stmt.QueryRow(merchantId).Scan(
		&merchant.MerchantId,
		&merchant.MerchantName,
		&merchant.MerchantAddress,
		&merchant.SubAcquirerId,
		&merchant.SubAcquirerName,
	)

	if err != nil {
		log.Printf("Error scanning %v\n", err)
		return nil, err
	}
	log.Printf("Retornando merchant %v\n", merchant)

	return merchant, nil
}

func (mr *MerchantSqlRepository) FindAll() ([]model.Merchant, error) {
	stmt, err := mr.db.Prepare(`
		select merchant_id, merchant_name, merchant_address, sub_acquirer_id, sub_acquirer_name
		from merchant
	`)
	if err != nil {
		log.Printf("Error preparing statement %v\n", err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Error querying statement %v\n", err)
		return nil, err
	}
	defer rows.Close()

	merchants := []model.Merchant{}
	for rows.Next() {
		merchant := model.Merchant{}
		err := rows.Scan(
			&merchant.MerchantId,
			&merchant.MerchantName,
			&merchant.MerchantAddress,
			&merchant.SubAcquirerId,
			&merchant.SubAcquirerName,
		)
		if err != nil {
			log.Printf("Error scanning statement %v\n", err)
			return nil, err
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}
