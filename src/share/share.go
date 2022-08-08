package share

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "seol7532!?"
	dbname   = "postgres"
)

type Mushroom struct {
	name         string
	index        int
	damage       int
	amount       int
	image_source string
}

type MainData struct {
	//Mushrooms       [12]Mushroom
	MushroomNames   [12]string
	MushroomsEdible [12]bool
	MushroomKills   [12]int

	CurrentIndex int
	//CurrentIndex_Sub1 int
	CurrentStage int

	FriendAlive bool
	MeAlive     bool
	ThrowCounts int
	HealthMe    int
	HealthHe    int
}

var (
	Address      string
	TmpMushrooms = MainData{}

	MainDB *sql.DB
)

var (
	//MrInput Mushroom

	Mr_poisons     []Mushroom
	Mr_general     []Mushroom
	tmp12Mushrooms [12]Mushroom
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

func SettingDB() *sql.DB {

	// postgres DB에서 전부 불러와서 저장해놓음
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, _ := sql.Open("postgres", psqlInfo)
	//defer db.Close()

	rows, _ := db.Query(`SELECT * FROM public ."MushRoom" ORDER BY "Index" ASC`)
	defer rows.Close()
	for rows.Next() {
		var MrInput Mushroom

		_ = rows.Scan(&MrInput.name, &MrInput.index, &MrInput.damage, &MrInput.amount, &MrInput.image_source)

		//fmt.Fprintf(os.Stdout, "%v\n", MrInput)
		if MrInput.damage == 0 { // 일반 버섯일 때
			Mr_general = append(Mr_general, Mushroom{MrInput.name, MrInput.index, MrInput.damage, MrInput.amount, MrInput.image_source})
		} else {
			Mr_poisons = append(Mr_poisons, Mushroom{MrInput.name, MrInput.index, MrInput.damage, MrInput.amount, MrInput.image_source})
		}
	}

	// MainDB에게 Update를 위해 정보를 넘겨주는 과정
	return db
}

func (data *MainData) RandomlyGet12(db *sql.DB) {

	rndIdxPoisons := rand.Perm(len(Mr_poisons))[:4]
	rndIdxGeneral := rand.Perm(len(Mr_general))[:8]
	for i, iVal := range rndIdxPoisons {
		tmp12Mushrooms[i] = Mr_poisons[iVal]
	}
	for i, iVal := range rndIdxGeneral {
		tmp12Mushrooms[i+4] = Mr_general[iVal]
	}
	rand.Shuffle(len(tmp12Mushrooms), func(i, j int) { tmp12Mushrooms[i], tmp12Mushrooms[j] = tmp12Mushrooms[j], tmp12Mushrooms[i] })
	rand.Shuffle(len(tmp12Mushrooms), func(i, j int) { tmp12Mushrooms[i], tmp12Mushrooms[j] = tmp12Mushrooms[j], tmp12Mushrooms[i] })

	data.UpdateOther12(db)
	data.UpdateKills12(db)

	data.FriendAlive = true
	data.MeAlive = true
	data.ThrowCounts = 0
	data.HealthMe = 100
	data.HealthHe = 100

	data.CurrentIndex = 0
}

func (data *MainData) UpdateOther12(db *sql.DB) {

	for i := 0; i < 12; i++ {
		data.MushroomNames[i] = tmp12Mushrooms[i].name
		if dmg := tmp12Mushrooms[i].damage; dmg == 0 {
			data.MushroomsEdible[i] = true
		} else {
			data.MushroomsEdible[i] = false
		}
	}

}

func (data *MainData) UpdateKills12(db *sql.DB) {
	for i := 0; i < 12; i++ {
		name := tmp12Mushrooms[i].name
		var amount_check int
		row := db.QueryRow(fmt.Sprintf("SELECT \"Amount\" FROM public .\"MushRoom\" WHERE \"Name\" = '%v'", name))
		_ = row.Scan(&amount_check)
		//fmt.Fprintf(os.Stdout, "%v kills %v\n", name, amount_check)
		data.MushroomKills[i] = amount_check
	}
}

func SettingCookie(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != http.ErrNoCookie { // 쿠키가 있으면
		//fmt.Fprintf(os.Stdout, "err: %v ---- intro cookie [O]\n", err)
	} else {
		//fmt.Fprintf(os.Stdout, "err: %v ---- intro cookie [X], now a cookie wiil be setted\n", err)
		c := &http.Cookie{
			Name:  "session",
			Value: "",
		}
		http.SetCookie(w, c)
	}

}

func CheckIfCookie(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("session")
	if err != http.ErrNoCookie { // 쿠키가 있으면
		return true
	} else {
		//fmt.Fprintf(os.Stdout, "err: %v\n", err)
		return false
	}

}

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		log.Fatal(err)
	}
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	//http.Redirect(w, r, "/사이트주소", http.StatusFound)
}

//userID := r.PostFormValue("username")
//<input name="username" class="form-control" placeholder="Username" required autofocus>
