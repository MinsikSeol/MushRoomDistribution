// 아래와 같은 부분이 필요할 것이다 수정 필요!!
// delete(share.UserMushrooms, str)

package main

//
import (
	"MushRoomDistribution/src/mainSub"
	"MushRoomDistribution/src/share"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// pgAdmin4 information
// username = "postgres"
// password = "{facebook password}"

type mushroom struct {
	name   string
	index  int
	damage int
	amount int
}

type fileURLs struct {
	URLs []string
}

func main() {
	share.MainDB = share.SettingDB()
	defer share.MainDB.Close()

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	println(dir)

	http.Handle("/src/",
		http.StripPrefix("/src/",
			http.FileServer(http.Dir("src"))))

	// Essential!!
	share.UserMushrooms = map[string]*share.MainData{}

	//share.Address = "http://localhost:8080"

	share.Address = "https://mushroomsms.herokuapp.com"

	http.HandleFunc("/", mainpage)
	//--------------------
	//http.HandleFunc("/prepareJS", prepareJS)
	//http.HandleFunc("/test", test)
	//--------------------
	http.HandleFunc("/about", about)
	http.HandleFunc("/intro", intro)
	http.HandleFunc("/Question0", mainSub.Question0)
	//--------------------
	http.HandleFunc("/you_alive", you_alive)
	http.HandleFunc("/you_dead", you_dead)
	http.HandleFunc("/friend_alive", friend_alive)
	http.HandleFunc("/friend_dead", friend_dead)
	//--------------------
	http.HandleFunc("/result", result)

	//http.ListenAndServe(":8080", nil)
	//--------------------

	// for heroku -> Address 변경필요
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}

func test(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/main/transitionTest.html"}}
	//files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	t.Execute(w, nil)
}

func prepareJS(w http.ResponseWriter, r *http.Request) {
	// 404 not found
	//http.ServeFile(w, r, "/src/static/js/jquery.min.js")
}

func mainpage(w http.ResponseWriter, r *http.Request) {
	var t *template.Template

	// 아래와 같은 Delete함수가 필요하다!!
	//delete(share.UserMushrooms, str)

	if r.Header.Get("Referer") == share.Address+"/" {
		_, err := r.Cookie("session")
		if err != http.ErrNoCookie {
			share.DeleteCookie(w, r)
		}
	}
	share.SettingCookie(w, r)

	//if r.Header.Get("Referer") == share.Address+"/intro" {
	//	http.Redirect(w, r, "/intro", http.StatusFound)
	//}

	/*// VERY IMPORTANT!!
	share.UserMushrooms[str] = &share.MainData{}

	if r.Header.Get("Referer") == share.Address+"/intro" {
		share.UserMushrooms[str].RandomlyGet12(share.MainDB)
		if entry, okay := share.UserMushrooms[str]; okay {
			//entry.CurrentIndex = 0
			//entry.CurrentStage = 0

			c, _ := r.Cookie("session")
			entry.LocalCookie = c.Value
			share.UserMushrooms[str] = entry
		}
		// 위 코드는 아래의 두 줄과 같다
		//share.UserMushrooms[str].CurrentIndex = 0
		//share.UserMushrooms[str].CurrentStage = 0
	}*/

	files := &fileURLs{URLs: []string{"src/main/startpage.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ = (template.ParseFiles(files.URLs...))

	t.Execute(w, nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/main/aboutpage.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, err := (template.ParseFiles(files.URLs...))
	if err != nil {
		log.Println(err.Error())
		return
	}

	t.Execute(w, nil)
}

func intro(w http.ResponseWriter, r *http.Request) {

	var str string
	if share.TakeCookieStr(r) != "" {
		str = share.TakeCookieStr(r)
	}

	if r.Header.Get("Referer") == share.Address+"/result" {
		share.UserMushrooms[str] = &share.MainData{}
		share.UserMushrooms[str].RandomlyGet12(share.MainDB)
	}

	//fmt.Fprintf(w, "%v 가나다라 \n", r.Header.Get("Referer"))

	if r.Header.Get("Referer") == share.Address+"/" {

		// VERY IMPORTANT!!
		share.UserMushrooms[str] = &share.MainData{}
		share.UserMushrooms[str].RandomlyGet12(share.MainDB)
		if entry, okay := share.UserMushrooms[str]; okay {

			c, _ := r.Cookie("session")
			entry.LocalCookie = c.Value
			share.UserMushrooms[str] = entry
		}
	}

	files := &fileURLs{URLs: []string{"src/main/intropage.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	t.Execute(w, nil)
	//giver := share.UserMushrooms[str]
	//t.Execute(w, giver)
}

func you_alive(w http.ResponseWriter, r *http.Request) {

	var str string

	if r.Header.Get("Referer") == share.Address+"/Question0" {
		if share.TakeCookieStr(r) != "" {
			str = share.TakeCookieStr(r)
		}

		if r.Method == "GET" {
			if entry, okay := share.UserMushrooms[str]; okay {
				entry.HealthMe = 0
				entry.MeAlive = false
				share.UserMushrooms[str] = entry
			}
		} else {
			if entry, okay := share.UserMushrooms[str]; okay {
				entry.HealthMe += 50
				entry.HealthHe -= 25
				if entry.HealthMe >= 100 {
					entry.HealthMe = 100
				}
				if entry.HealthHe <= 0 {
					entry.HealthHe = 0
					entry.FriendAlive = false
				}
				share.UserMushrooms[str] = entry
			}
		}

	}

	files := &fileURLs{URLs: []string{"src/mainSub/choose_eat_alive.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if entry, okay := share.UserMushrooms[str]; okay {
		t.Execute(w, entry)
	}
}

func you_dead(w http.ResponseWriter, r *http.Request) {

	var str string

	if r.Header.Get("Referer") == share.Address+"/Question0" && r.Method == "POST" {
		if share.TakeCookieStr(r) != "" {
			str = share.TakeCookieStr(r)
		}

		/*idx := share.UserMushrooms[str].CurrentIndex
		you_eaten := share.UserMushrooms[str].MushroomNames[idx]
		Increase1(you_eaten)*/
		if entry, okay := share.UserMushrooms[str]; okay {
			idx := entry.CurrentIndex
			you_eaten := entry.MushroomNames[idx]
			Increase1(you_eaten)
		}

		// me_alive 관련
		if entry, okay := share.UserMushrooms[str]; okay {
			entry.HealthMe = 0
			entry.MeAlive = false
			share.UserMushrooms[str] = entry
		}
	}

	files := &fileURLs{URLs: []string{"src/mainSub/choose_eat_dead.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if entry, okay := share.UserMushrooms[str]; okay {
		t.Execute(w, entry)
	}
}

func friend_alive(w http.ResponseWriter, r *http.Request) {

	var str string

	if r.Header.Get("Referer") == share.Address+"/Question0" {
		if share.TakeCookieStr(r) != "" {
			str = share.TakeCookieStr(r)
		}

		if entry, okay := share.UserMushrooms[str]; okay {
			entry.HealthHe += 50
			entry.HealthMe -= 25
			if entry.HealthHe >= 100 {
				entry.HealthHe = 100
			}
			if entry.HealthMe <= 0 {
				entry.HealthMe = 0
				entry.MeAlive = false
			}
			share.UserMushrooms[str] = entry
		}
	}

	files := &fileURLs{URLs: []string{"src/mainSub/choose_friend_alive.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if entry, okay := share.UserMushrooms[str]; okay {
		t.Execute(w, entry)
	}
}

func friend_dead(w http.ResponseWriter, r *http.Request) {

	var str string

	//fmt.Fprintf(os.Stdout, "%v", r.Header.Get("Referer"))
	if r.Header.Get("Referer") == share.Address+"/Question0" {
		if share.TakeCookieStr(r) != "" {
			str = share.TakeCookieStr(r)
		}

		/*idx := share.UserMushrooms[str].CurrentIndex
		friend_eaten := share.UserMushrooms[str].MushroomNames[idx]
		Increase1(friend_eaten)*/
		if entry, okay := share.UserMushrooms[str]; okay {
			idx := entry.CurrentIndex
			friend_eaten := entry.MushroomNames[idx]
			Increase1(friend_eaten)
		}

		// friend_alive 관련
		if entry, okay := share.UserMushrooms[str]; okay {
			entry.HealthHe = 0
			entry.FriendAlive = false
			share.UserMushrooms[str] = entry
		}

		// me_alive 관련
		if entry, okay := share.UserMushrooms[str]; okay {
			entry.HealthMe -= 25
			if entry.HealthMe <= 0 {
				entry.HealthMe = 0
				entry.MeAlive = false
			}
			share.UserMushrooms[str] = entry
		}

	}

	files := &fileURLs{URLs: []string{"src/mainSub/choose_friend_dead.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if entry, okay := share.UserMushrooms[str]; okay {
		t.Execute(w, entry)
	}
}

func Increase1(str string) {
	//_, err := share.MainDB.Query(`SELECT * FROM public ."MushRoom" ORDER BY "Index" ASC`)
	db0, err := share.MainDB.Query(fmt.Sprintf("UPDATE public .\"MushRoom\" SET \"Amount\" = (\"Amount\" + 1) WHERE \"Name\" = '%v'", str))
	if err != nil {
		panic(err)
	}
	db0.Close()
}

func result(w http.ResponseWriter, r *http.Request) {

	var str string
	if share.TakeCookieStr(r) != "" {
		str = share.TakeCookieStr(r)
	}

	files := &fileURLs{URLs: []string{"src/mainSub/result.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if entry, okay := share.UserMushrooms[str]; okay {
		t.Execute(w, entry)
	}
	//giver := share.UserMushrooms[str]
	//t.Execute(w, giver)

	/*if r.Header.Get("Referer") == share.Address+"/Question0" {
		delete(share.UserMushrooms, str)
	}*/
}
