package share

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/*const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "채울 부분"
	dbname   = "postgres"
)*/

type Mushroom struct {
	name   string
	index  int
	damage int
	amount int
}

type MainData struct {
	//Mushrooms       [12]Mushroom
	MushroomNames   [12]string
	MushroomsEdible [12]bool
	MushroomKills   [12]int

	CurrentIndex int
	CurrentStage int

	FriendAlive bool
	MeAlive     bool
	ThrowCounts int
	HealthMe    int
	HealthHe    int

	LocalCookie string
}

var (
	UserMushrooms   = map[string]*MainData{}
	staticCookieNum int
	Address         string
	//TmpMushrooms  = MainData{}

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
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, _ := sql.Open("postgres", psqlInfo)*/
	//defer db.Close()

	// using Heroku Postgres
	//db, _ := sql.Open("postgres", "postgres://rsdivmtfbvning:f2fd9f0d65e7269c89af3bef5fb96dab3906de09dd366400e7ab24e8d4a4a129@ec2-3-223-242-224.compute-1.amazonaws.com:5432/dfs6qrqgepkpen")
	db, _ := sql.Open("postgres", "postgres://gmwhzaxuurlxbw:5af79b3371b1046a6c41d2a0e659ee1fd0e92d13f20b91d9a66bd6e63bd5440d@ec2-34-231-42-166.compute-1.amazonaws.com:5432/d5m9tj9eruhjic")

	rows, _ := db.Query(`SELECT * FROM public ."MushRoom" ORDER BY "Index" ASC`)
	defer rows.Close()
	for rows.Next() {
		var MrInput Mushroom

		_ = rows.Scan(&MrInput.name, &MrInput.index, &MrInput.damage, &MrInput.amount)

		//fmt.Fprintf(os.Stdout, "%v\n", MrInput)
		if MrInput.damage == 0 { // 일반 버섯일 때
			Mr_general = append(Mr_general, Mushroom{MrInput.name, MrInput.index, MrInput.damage, MrInput.amount})
		} else {
			Mr_poisons = append(Mr_poisons, Mushroom{MrInput.name, MrInput.index, MrInput.damage, MrInput.amount})
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

	data.UpdateOther12()
	data.UpdateKills12(db)

	data.Initialization()
}

func (data *MainData) Initialization() {
	data.FriendAlive = true
	data.MeAlive = true
	data.ThrowCounts = 0
	data.HealthMe = 100
	data.HealthHe = 100

	data.CurrentIndex = 0
	data.CurrentStage = 0

	data.LocalCookie = "cookie0"
}

func (data *MainData) UpdateOther12() {

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
		tmpCookieStr := MakeCookieStr()
		c := &http.Cookie{
			Name:  "session",
			Value: tmpCookieStr,
		}
		http.SetCookie(w, c)

	}

}

/*func (data *MainData) GettingCookie(r *http.Request) (string, error) {
	c, err := r.Cookie("session")
	if err != http.ErrNoCookie { // 쿠키가 있으면
		return c.Value, nil
	} else {
		//fmt.Fprintf(os.Stdout, "err: %v\n", err)
		return "", err
	}
}*/

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session")
	c = &http.Cookie{
		Name:   "session",
		Value:  "cookie0",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusFound)
}

/*func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}*/

func MakeCookieStr() string {
	staticCookieNum++
	return "cookie" + strconv.Itoa(staticCookieNum)
}

func TakeCookieStr(r *http.Request) string {
	c, err := r.Cookie("session")
	if err != http.ErrNoCookie { // 쿠키가 있으면
		return c.Value
	} else {
		return ""
	}
}
