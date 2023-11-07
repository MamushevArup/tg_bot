package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MamushevArup/krisha-scraper/models"
	uuid2 "github.com/google/uuid"
	"github.com/lib/pq"
	"log"
)

func (s *Sql) userExist(username *models.User) (uuid2.UUID, error) {
	query := `select id from users where username = $1`
	var id uuid2.UUID
	err := s.Db.QueryRow(query, username.Username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// user does not exist so it returns zero id.
			return id, err
		} else {
			log.Fatal("Cannot get info in the userExist method ", err)
		}
	}
	// user exist and returns id
	return id, nil
}

func (s *Sql) IntroDataStartCommand(user *models.User) {
	uid, err := s.userExist(user)
	if err != nil {
		s.insertIntro(user)
		return
	}
	s.id = uid
	s.updateIntro(uid, user)
}

func (s *Sql) insertIntro(user *models.User) {
	// if user do not exists
	uuid := uuid2.New()
	s.id = uuid
	query := `insert into users(id, username, buyOrRent, typeItem, city) values ($1, $2, $3, $4, $5)`
	_, err := s.Db.Exec(query, uuid, user.Username, user.BuyOrRent, user.TypeItem, user.City)
	if err != nil {
		log.Fatal("Error with inserting values to the user InsertIntro method ", err)
	}
}

func (s *Sql) updateIntro(id uuid2.UUID, users *models.User) {
	query := `update users set buyOrRent = $2, typeItem = $3, city = $4 where id = $1`
	_, err := s.Db.Exec(query, id, users.BuyOrRent, users.TypeItem, users.City)
	if err != nil {
		log.Println("Error cannot update in t"+
			"he updateInto method ", err)
	}
}

func (s *Sql) UpdateCity(city string) error {
	query := `update users set city = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, city)
	if err != nil {
		log.Println("Cannot update city in the UpdateCity method ", err)
		return err
	}
	return nil
}

func (s *Sql) UpdateRooms(rooms []string) error {
	query := `update users set rooms = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, pq.Array(rooms))
	if err != nil {
		log.Println("Cannot update rooms in the UpdateRooms method ", err)
		return err
	}
	return nil
}

func (s *Sql) UpdateType(types []string) error {
	query := `update users set typeHouse = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, pq.Array(types))
	if err != nil {
		log.Println("Cannot update type in the UpdateType method ", err)
		return err
	}
	return nil
}

func (s *Sql) UpdateBuiltFrom(year uint) error {
	query := `update users set yearBuiltFrom = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, year)
	if err != nil {
		log.Println("Cannot update builtFrom in the UpdateBuiltFrom method ", err)
		return err
	}
	return nil
}

func (s *Sql) UpdateBuiltTo(year uint) error {
	query := `update users set yearBuiltTo = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, year)
	if err != nil {
		log.Println("Cannot update builtTo in the UpdateBuiltTo method ", err)
		return err
	}
	return nil
}

func (s *Sql) UpdatePriceFrom(price uint64) error {
	query := `update users set priceFrom = $2 where id = $1`
	fmt.Println(s.id)
	_, err := s.Db.Exec(query, s.id, price)
	if err != nil {
		log.Println("Cannot update priceFrom in the UpdatePriceFrom method")
		return err
	}
	return nil
}
func (s *Sql) UpdatePriceTo(price uint64) error {
	query := `update users set priceTo = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, price)
	if err != nil {
		log.Println("Cannot update priceTo in the UpdatePriceTo method")
		return err
	}
	return nil
}
func (s *Sql) UpdateFloorFrom(floor uint64) error {
	query := `update users set floorFrom = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, floor)
	if err != nil {
		log.Println("Cannot update floorFrom in the UpdateFloorFrom method")
		return err
	}
	return nil
}
func (s *Sql) UpdateFloorTo(floor uint64) error {
	query := `update users set floorTo = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, floor)
	if err != nil {
		log.Println("Cannot update floorTo in the UpdateFloorTo method")
		return err
	}
	return nil
}
func (s *Sql) UpdateFloorInTheHouseFrom(floorHouse uint64) error {
	query := `update users set floorInTheHouseFrom = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, floorHouse)
	if err != nil {
		log.Println("Cannot update floorHouse in the UpdateFloorInTheHouseFrom method")
		return err
	}
	return nil
}
func (s *Sql) UpdateFloorInTheHouseTo(floorHouse uint64) error {
	query := `update users set floorInTheHouseTo = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, floorHouse)
	if err != nil {
		log.Println("Cannot update floorHouse in the UpdateFloorInTheHouseTo method")
		return err
	}
	return nil
}
func (s *Sql) UpdateAreaFrom(area string) error {
	query := `update users set areaFrom = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, area)
	if err != nil {
		log.Println("Cannot update area in the UpdateAreaFrom method")
		return err
	}
	return nil
}
func (s *Sql) UpdateAreaTo(area string) error {
	query := `update users set areaTo = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, area)
	if err != nil {
		log.Println("Cannot update area in the UpdateAreaTo method")
		return err
	}
	return nil
}
func (s *Sql) UpdateKitchenFrom(kitchen string) error {
	query := `update users set kitchenFrom = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, kitchen)
	if err != nil {
		log.Println("Cannot update kitchen in the UpdateKitchenFrom method")
		return err
	}
	return nil
}
func (s *Sql) UpdateKitchenTo(kitchen string) error {
	query := `update users set kitchenTo = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, kitchen)
	if err != nil {
		log.Println("Cannot update kitchen in the UpdateKitchenTo method")
		return err
	}
	return nil
}
func (s *Sql) UpdateNotFirstFloor(flag bool) error {
	query := `update users set notFirstFloor = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, flag)
	if err != nil {
		log.Println("Cannot update notFirstFloor in the UpdateNotFirstFloor method")
		return err
	}
	return nil
}
func (s *Sql) UpdateNotLastFloor(flag bool) error {
	query := `update users set notLastFloor = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, flag)
	if err != nil {
		log.Println("Cannot update notLastFloor in the UpdateNotLastFloor method")
		return err
	}
	return nil
}
func (s *Sql) UpdateFromOwner(flag bool) error {
	query := `update users set fromOwner = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, flag)
	if err != nil {
		log.Println("Cannot update fromOwner in the UpdateFromOwner method")
		return err
	}
	return nil
}
func (s *Sql) UpdateNewBuilding(flag bool) error {
	query := `update users set newBuilding = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, flag)
	if err != nil {
		log.Println("Cannot update newBuilding in the UpdateNewBuilding method")
		return err
	}
	return nil
}
func (s *Sql) UpdateRealEstate(flag bool) error {
	query := `update users set realEstate = $2 where id = $1`
	_, err := s.Db.Exec(query, s.id, flag)
	if err != nil {
		log.Println("Cannot update realEstate in the UpdateRealEstate method")
		return err
	}
	return nil
}
func (s *Sql) GetAll() (*models.User, error) {
	user := models.User{}
	query := `select * from users where id = $1`
	err := s.Db.Get(&user, query, s.id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Sql) Insert() {
	uuid, _ := uuid2.NewUUID()
	query := `insert into users(id, username, buyorrent, typeitem) values($1, 'what', '1', '2')`
	_, err := s.Db.Exec(query, uuid)
	if err != nil {
		log.Fatal(err)
	}
}
