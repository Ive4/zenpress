package handler

import (
	"net/http"

	"github.com/Unknwon/com"
	"github.com/insionng/vodka"
	"github.com/insionng/zenpress/models"
)

func ViewHandler(self vodka.Context) error {

	data := make(map[string]interface{})
	tid := com.StrTo(self.Param("tid")).MustInt64()
	tid_handler := models.GetTopic(tid)

	if tid_handler.Id > 0 {

		tid_handler.Views = tid_handler.Views + 1
		models.UpdateTopic(tid, tid_handler)

		data["article"] = tid_handler
		data["replys"] = models.GetReplyByPid(tid, 0, 0, "id")

		tps := models.GetAllTopicByCid(tid_handler.Cid, 0, 0, 0, "asc")

		if tps != nil && tid != 0 {

			for i, v := range tps {

				if v.Id == tid {
					prev := i - 1
					next := i + 1

					for i, v := range tps {
						if prev == i {
							data["previd"] = v.Id
							data["prev"] = v.Title
						}
						if next == i {
							data["nextid"] = v.Id
							data["next"] = v.Title
						}
					}
				}
			}
		}

		self.SetStore(data)
		return self.Render(http.StatusOK, "view.html")
	}

	return self.Redirect(302, "/")

}
