package mainSub

import (
	"MushRoomDistribution/src/share"
	"html/template"
	"net/http"
)

type fileURLs struct {
	URLs []string
}

func Question0(w http.ResponseWriter, r *http.Request) {

	var str string
	if share.TakeCookieStr(r) != "" {
		str = share.TakeCookieStr(r)
	}

	//fmt.Fprintf(w, "%v 가나다라 \n", str)

	if r.Method == "POST" {

		if r.Header.Get("Referer") == share.Address+"/Question0" {

			if entry, okay := share.UserMushrooms[str]; okay {

				entry.ThrowCounts += 1

				entry.HealthMe -= 25
				entry.HealthHe -= 25
				if entry.HealthMe <= 0 {
					entry.HealthMe = 0
					entry.MeAlive = false
				}
				if entry.HealthHe <= 0 {
					entry.HealthHe = 0
					entry.FriendAlive = false
				}
				share.UserMushrooms[str] = entry
			}

		}
		if entry, okay := share.UserMushrooms[str]; okay {

			entry.CurrentIndex += 1
			entry.CurrentStage = entry.CurrentIndex + 1

			share.UserMushrooms[str] = entry
		}

		http.Redirect(w, r, "/Question0", http.StatusFound)

	}

	if entry, okay := share.UserMushrooms[str]; okay {
		if entry.MeAlive == false {
			http.Redirect(w, r, "/result", http.StatusFound)
		}
	}

	files := &fileURLs{URLs: []string{"src/mainSub/Question0.html"}}
	files.URLs = append(files.URLs, "src/main/homepage.html")

	t, _ := template.ParseFiles(files.URLs...)

	if entry, okay := share.UserMushrooms[str]; okay {
		t.Execute(w, entry)
	}
}
