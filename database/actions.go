package database

type PhoneNumber struct {
	ID     int
	Number string
}

func (pdb PhoneDb) InsertNumbers(numbers []string) error {
	const insert = `
	INSERT INTO phones (number)
	VALUES ($1)`

	for _, p := range numbers {
		_, err := pdb.sqlDb.Exec(insert, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdb PhoneDb) UpdateNumbers(phone PhoneNumber) error {
	const update = `
	UPDATE phones
	SET number = $1
	WHERE id = $2`

	_, err := pdb.sqlDb.Exec(update, phone.Number, phone.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pdb PhoneDb) GetNumbers() ([]PhoneNumber, error) {
	const query = `SELECT id, number FROM phones`

	rows, err := pdb.sqlDb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	numbers := make([]PhoneNumber, 0)
	for rows.Next() {
		var pn PhoneNumber

		err = rows.Scan(&pn.ID, &pn.Number)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, pn)

	}

	return numbers, nil
}

func (pdb PhoneDb) Close() error {
	err := pdb.sqlDb.Close()
	return err
}
