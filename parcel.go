package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()
	res, err := db.Exec("INSERT INTO parcel(client, status, address, created_at)"+
		"VALUES(:client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()

	// верните идентификатор последней добавленной записи
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return p, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE "+
		"number = :number", sql.Named("number", number))
	err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT number, client, status, address, created_at from parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	var status string

	row := db.QueryRow("SELECT status from parcel WHERE number = :number", sql.Named("number", number))
	err = row.Scan(&status)
	if err != nil {
		return err
	}

	if status != ParcelStatusRegistered {
		return nil
	}

	_, err = db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("number", number),
		sql.Named("address", address))

	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	var status string

	row := db.QueryRow("SELECT status from parcel WHERE number = :number", sql.Named("number", number))
	err = row.Scan(&status)
	if err != nil {
		return err
	}

	if status != ParcelStatusRegistered {
		return nil
	}

	_, err = db.Exec("DELETE FROM parcel WHERE number = :number",
		sql.Named("number", number))
	return err
}
