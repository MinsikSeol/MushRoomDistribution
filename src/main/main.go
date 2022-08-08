package main

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

	share.TmpMushrooms.CurrentIndex = 0

	//share.Address = "http://localhost:8080"

	port := os.Getenv("PORT")
	share.Address = ":" + port

	http.HandleFunc("/", mainpage)
	//--------------------
	//http.HandleFunc("/prepareJS", prepareJS)
	http.HandleFunc("/test", test)
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
	http.ListenAndServe(share.Address, nil)

}

func test(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/main/transitionTest.html"}}
	//files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	t.Execute(w, nil)
}

/*
func prepareJS(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("/src/static/js/transitionTest.js")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}
*/

func prepareJS(w http.ResponseWriter, r *http.Request) {
	// 404 not found
	//http.ServeFile(w, r, "/src/static/js/jquery.min.js")
}

func mainpage(w http.ResponseWriter, r *http.Request) {

	files := &fileURLs{URLs: []string{"src/main/startpage.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, err := (template.ParseFiles(files.URLs...))
	if err != nil {
		log.Println(err.Error())
		return
	}

	//fmt.Fprintf(os.Stdout, "%v %v %v\n", r.Method, r.Body, r.GetBody)
	//fmt.Fprintf(os.Stdout, "%v %v %v\n", r.Form, r.PostForm, r.TLS)
	//fmt.Fprintf(os.Stdout, "%v %v %v\n", r.URL, r.RequestURI, r.RemoteAddr)
	//fmt.Fprintf(os.Stdout, "%v %v %v\n", r.MultipartForm, r.ContentLength, r.Host)
	//fmt.Fprintln(os.Stdout, "")

	//fmt.Fprintf(os.Stdout, "%v", r.Header.Get("Referer"))
	if r.Header.Get("Referer") == share.Address+"/intro" {
		share.TmpMushrooms.RandomlyGet12(share.MainDB)
		share.TmpMushrooms.CurrentIndex = 0
		//share.TmpMushrooms.CurrentIndex_Sub1 = 0
		share.TmpMushrooms.CurrentStage = 0
	}

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
	files := &fileURLs{URLs: []string{"src/main/intropage.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, err := (template.ParseFiles(files.URLs...))
	if err != nil {
		log.Println(err.Error())
		return
	}

	t.Execute(w, nil)
}

func you_alive(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/mainSub/choose_eat_alive.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))
	if r.Method == "GET" {
		share.TmpMushrooms.HealthMe = 0
		if share.TmpMushrooms.HealthMe <= 0 {
			share.TmpMushrooms.MeAlive = false
		}

	} else if r.Header.Get("Referer") == share.Address+"/Question0" {
		share.TmpMushrooms.HealthMe += 50
		if share.TmpMushrooms.HealthMe >= 100 {
			share.TmpMushrooms.HealthMe = 100
		}

		share.TmpMushrooms.HealthHe -= 25
		if share.TmpMushrooms.HealthHe <= 0 {
			share.TmpMushrooms.FriendAlive = false
		}
	}

	t.Execute(w, &share.TmpMushrooms)
}

func you_dead(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/mainSub/choose_eat_dead.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if r.Header.Get("Referer") == share.Address+"/Question0" {
		idx := share.TmpMushrooms.CurrentIndex
		you_eaten := share.TmpMushrooms.MushroomNames[idx]
		Increase1(you_eaten)

		// me_alive 관련
		share.TmpMushrooms.HealthMe = 0
		if share.TmpMushrooms.HealthMe <= 0 {
			share.TmpMushrooms.MeAlive = false
		}
	}

	t.Execute(w, &share.TmpMushrooms)
}

func friend_alive(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/mainSub/choose_friend_alive.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	if r.Header.Get("Referer") == share.Address+"/Question0" {
		share.TmpMushrooms.HealthHe += 50
		if share.TmpMushrooms.HealthHe >= 100 {
			share.TmpMushrooms.HealthHe = 100
		}

		share.TmpMushrooms.HealthMe -= 25
		if share.TmpMushrooms.HealthMe <= 0 {
			share.TmpMushrooms.MeAlive = false
		}
	}

	t.Execute(w, &share.TmpMushrooms)
}

func friend_dead(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/mainSub/choose_friend_dead.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	//fmt.Fprintf(os.Stdout, "%v", r.Header.Get("Referer"))
	if r.Header.Get("Referer") == share.Address+"/Question0" {
		idx := share.TmpMushrooms.CurrentIndex
		//fmt.Println(share.TmpMushrooms.MushroomNames[idx])
		friend_eaten := share.TmpMushrooms.MushroomNames[idx]
		//fmt.Fprintf(os.Stdout, "%v 을 먹었음", you_eaten)
		Increase1(friend_eaten)

		// friend_alive 관련
		share.TmpMushrooms.HealthHe = 0
		if share.TmpMushrooms.HealthHe <= 0 {
			share.TmpMushrooms.FriendAlive = false
		}

		// me_alive 관련
		share.TmpMushrooms.HealthMe -= 25
		if share.TmpMushrooms.HealthMe <= 0 {
			share.TmpMushrooms.MeAlive = false
		}

	}

	t.Execute(w, &share.TmpMushrooms)
}

func Increase1(str string) {
	//_, err := share.MainDB.Query(`SELECT * FROM public ."MushRoom" ORDER BY "Index" ASC`)
	_, err := share.MainDB.Query(fmt.Sprintf("UPDATE public .\"MushRoom\" SET \"Amount\" = (\"Amount\" + 1) WHERE \"Name\" = '%v'", str))
	if err != nil {
		panic(err)
	}
}

func result(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/mainSub/result.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := (template.ParseFiles(files.URLs...))

	t.Execute(w, &share.TmpMushrooms)
}
