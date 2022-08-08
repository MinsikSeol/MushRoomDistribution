package mainSub

import (
	"MushRoomDistribution/src/share"
	"html/template"
	"log"
	"net/http"
)

type fileURLs struct {
	URLs []string
}

func Question0(w http.ResponseWriter, r *http.Request) {
	files := &fileURLs{URLs: []string{"src/mainSub/Question0.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, err := (template.ParseFiles(files.URLs...))
	if err != nil {
		log.Println(err.Error())
		return
	}

	//fmt.Fprintf(os.Stdout, "%v\n", r.Method)

	if r.Method == "POST" {
		if r.Header.Get("Referer") == share.Address+"/Question0" {
			share.TmpMushrooms.ThrowCounts += 1

			share.TmpMushrooms.HealthMe -= 25
			if share.TmpMushrooms.HealthMe <= 0 {
				share.TmpMushrooms.MeAlive = false
			}

			share.TmpMushrooms.HealthHe -= 25
			if share.TmpMushrooms.HealthHe <= 0 {
				share.TmpMushrooms.FriendAlive = false
			}
		}

		share.TmpMushrooms.CurrentIndex += 1
		//share.TmpMushrooms.CurrentIndex_Sub1 = share.TmpMushrooms.CurrentIndex - 1
		share.TmpMushrooms.CurrentStage = share.TmpMushrooms.CurrentIndex + 1

		http.Redirect(w, r, "/Question0", http.StatusFound)
	}

	t.Execute(w, &share.TmpMushrooms)

}
