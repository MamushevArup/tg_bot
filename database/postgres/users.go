package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	uuid2 "github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)

func (s *Sql) GetUser(user *models.User) {
	query := `select id, username from users where username = $1`
	var id uuid2.UUID
	var username string
	err := s.Db.QueryRow(query, user.Username).Scan(&id, &username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.insertUser(user)
		} else {
			log.Fatal("Cannot get some info error is ", err)
		}
	} else {
		s.updateUser(user)
	}
}

func (s *Sql) CheckForStart(user *models.User) (string, string) {
	query := `select buyOrRent, typeItem from users where username = $1`
	var buy, typeI string
	_ = s.Db.QueryRow(query, user.Username).Scan(&buy, &typeI)
	return buy, typeI
}

func (s *Sql) insertUser(user *models.User) {
	uuid := uuid2.New()
	uuidFrom := uuid2.New()
	uuidTo := uuid2.New()
	uuidCheck := uuid2.New()
	in := user.UserChoice

	tn, err := s.Db.Begin()
	if err != nil {
		log.Println("Error with start a transaction", err)
		return
	}

	queryFrom := `insert into datafrom values 
                         ($1, $2, $3, $4, $5, $6, $7 )`
	_, err = tn.Exec(queryFrom, uuidFrom, in.YearOfBuiltFrom, in.PriceFrom,
		in.FloorFrom, in.FloorInTheHouseFrom, in.AreaFrom, in.KitchenAreaFrom)
	if err != nil {
		tn.Rollback()
		log.Println("Cannot insert to the datafrom table", err)
	}
	queryTo := `insert into datato values 
                       ($1,$2,$3,$4,$5,$6,$7)`
	_, err = tn.Exec(queryTo, uuidTo, in.YearOfBuiltTo, in.PriceTo,
		in.FloorTo, in.FloorInTheHouseTo, in.AreaTo, in.KitchenAreaTo)
	if err != nil {
		tn.Rollback()
		log.Println("Cannot insert to the datato table", err)
	}
	queryCheck := `insert into checkbox values 
                         ($1,$2,$3,$4,$5,$6)`
	_, err = tn.Exec(queryCheck, uuidCheck, in.CheckboxNotFirstFloor,
		in.CheckboxNotLastFloor, in.CheckboxFromOwner, in.CheckboxNewBuilding,
		in.CheckRealEstate)
	if err != nil {
		log.Println("Cannot insert to the check table", err)
	}
	query := ` insert into users values 
	                      ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = tn.Exec(query, uuid, user.Username, in.BuyOrRent, in.TypeItem,
		in.City, pq.Array(in.Rooms), pq.Array(in.TypeHouse), uuidFrom, uuidTo, uuidCheck)
	if err != nil {
		tn.Rollback()
		log.Println("Cannot insert to the users table", err)

	}
	err = tn.Commit()
	if err != nil {
		log.Println("Cannot confirm a commit command", err)
		return
	}
	fmt.Println("Succesfully end")
}

func (s *Sql) updateUser(user *models.User) error {
	var err error
	query, param := s.updateUsersTableQuery(user)
	queryFrom, paramFrom := s.updateDataFromTable(user)
	queryTo, paramsTo := s.updateDataToTable(user)
	queryCheck, paramsCheck := s.updateCheckTable(user)
	tn, err := s.Db.Begin()
	if err != nil {
		log.Println("Error with start a transaction ", err)
		return err
	}
	execTransaction(query, param, tn, "users")
	execTransaction(queryFrom, paramFrom, tn, "dataFrom")
	execTransaction(queryTo, paramsTo, tn, "dataTo")
	execTransaction(queryCheck, paramsCheck, tn, "checkbox")

	err = tn.Commit()
	if err != nil {
		log.Println("Cannot confirm that commit is happen ", err)
		return err
	}
	return nil
}
func execTransaction(query string, param []interface{}, tn *sql.Tx, tableName string) {
	_, err := tn.Exec(query, param...)
	if err != nil {
		tn.Rollback()
		log.Printf("Error with updating %v table %s", tableName, err)
		return
	}
}

func (s *Sql) updateUsersTableQuery(user *models.User) (string, []interface{}) {
	hmap := map[string]interface{}{
		"buyOrRent": user.BuyOrRent,
		"typeItem":  user.TypeItem,
		"city":      user.City,
		"rooms":     pq.StringArray(user.Rooms),
		"typeHouse": pq.StringArray(user.TypeHouse),
	}

	var query strings.Builder
	var params []interface{}
	params = append(params, user.Username)

	for column, value := range hmap {
		if value != nil {
			if arr, ok := value.(pq.StringArray); ok {
				query.WriteString(column + " = $" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, arr)
			} else if strValue, _ := value.(string); strValue != "" {
				query.WriteString(column + " = $" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, strValue)
			}
		}
	}
	queryString := strings.TrimSuffix(query.String(), ", ")
	finalQuery := "UPDATE users SET " + queryString + " WHERE username = $1"
	return finalQuery, params
}
func (s *Sql) updateDataFromTable(user *models.User) (string, []interface{}) {
	hmap := map[string]interface{}{
		"yearBuiltFrom":       user.YearOfBuiltFrom,
		"priceFrom":           user.PriceFrom,
		"floorFrom":           user.FloorFrom,
		"floorInTheHouseFrom": user.FloorInTheHouseFrom,
		"areaFrom":            user.AreaFrom,
		"kitchenFrom":         user.KitchenAreaFrom,
	}
	var query strings.Builder
	var params []interface{}
	params = append(params, user.Username)
	for k, v := range hmap {
		switch val := v.(type) {
		// uint8 uint64 string
		case uint:
			if val != 0 {
				query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, val)
			}
		case uint64:
			if val != 0 {
				query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, val)
			}
		case string:
			if val != "" {
				query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, val)
			}
		}
	}
	queryString := strings.TrimSuffix(query.String(), ", ")
	finalQuery := "UPDATE datafrom SET " + queryString + " WHERE id = (select idFrom from users where username = $1)"
	return finalQuery, params
}
func (s *Sql) updateDataToTable(user *models.User) (string, []interface{}) {
	hmap := map[string]interface{}{
		"yearBuiltTo":       user.YearOfBuiltTo,
		"priceTo":           user.PriceTo,
		"floorTo":           user.FloorTo,
		"floorInTheHouseTo": user.FloorInTheHouseTo,
		"areaTo":            user.AreaTo,
		"kitchenTo":         user.KitchenAreaTo,
	}
	var query strings.Builder
	var params []interface{}
	params = append(params, user.Username)
	for k, v := range hmap {
		switch val := v.(type) {
		// uint8 uint64 string
		case uint:
			if val != 0 {
				query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, val)
			}
		case uint64:
			if val != 0 {
				query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, val)
			}
		case string:
			if val != "" {
				query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
				params = append(params, val)
			}
		}
	}
	queryString := strings.TrimSuffix(query.String(), ", ")
	finalQuery := "UPDATE datato SET " + queryString + " WHERE id = (select idTo from users where username = $1)"
	return finalQuery, params
}

func (s *Sql) updateCheckTable(user *models.User) (string, []interface{}) {
	hmap := map[string]interface{}{
		"notFirstFloor": user.CheckboxNotFirstFloor,
		"notLastFloor":  user.CheckboxNotLastFloor,
		"fromOwner":     user.CheckboxFromOwner,
		"newBuilding":   user.CheckboxNewBuilding,
		"realEstate":    user.CheckRealEstate,
	}
	var query strings.Builder
	var params []interface{}
	params = append(params, user.Username)
	for k, v := range hmap {
		if v.(bool) {
			query.WriteString(k + "=$" + strconv.Itoa(len(params)+1) + ", ")
			params = append(params, v)
		}
	}
	queryString := strings.TrimSuffix(query.String(), ", ")
	finalQuery := "UPDATE checkbox SET " + queryString + " WHERE id = (select idCheck from users where username = $1)"
	return finalQuery, params
}
