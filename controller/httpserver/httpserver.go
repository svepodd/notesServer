package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"io"
	"net/http"
	"Notes/gates/storage"
	"Notes/models/dto"
	"Notes/pkg"
)

type HttpServer struct {
	srv http.Server
	st  storage.Storage
}

func NewHttpServer(addr string, st storage.Storage) (hs *HttpServer) {
	hs = new(HttpServer)
	hs.srv = http.Server{}
	mux := http.NewServeMux()
	mux.Handle("/create", http.HandlerFunc(hs.notesCreateHandler))
	mux.Handle("/get", http.HandlerFunc(hs.notesGetHandler))
	mux.Handle("/update", http.HandlerFunc(hs.notesUpdateHandler))
	mux.Handle("/delete", http.HandlerFunc(hs.notesDeleteByPhone))
	mux.Handle("/get-all", http.HandlerFunc(hs.notesGetAll))
	hs.srv.Handler = mux
	hs.srv.Addr = addr
	hs.st = st
	return hs
}

func (hs *HttpServer) Start() (err error) {
	eW := pkg.NewEWrapper("(hs *HttpServer) Start()")

	if err != nil {
		err = eW.WrapError(err, "pkg.NewEWrapper()")
		return
	}

	err = hs.srv.ListenAndServe()
	if err != nil {
		err = eW.WrapError(err, "hs.srv.ListenAndServe()")
		return
	}
	return
}

func (hs *HttpServer) notesCreateHandler(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	eW, err := pkg.NewEWrapperWithFile("(hs *HttpServer) notesCreateHandler()")
	if err != nil {
		log.Println("(hs *HttpServer) notesCreateHandler: NewEWrapperWithFile()", err)
	}

	resp := &dto.Response{}
	defer responseReturn(w, eW, resp)

	note := dto.NewNote()
	byteReq, err := io.ReadAll(req.Body)
	if err != nil {
		resp.Wrap("Error reading request", nil, err.Error())
		eW.LogError(err, "io.ReadAll(req.Body)")
		return
	}
	err = json.Unmarshal(byteReq, &note)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Unmarshal(req)")
		return
	}

	if note.Name == "" || note.LastName == "" || note.Note == "" {
		err = errors.New("required data is missing")
		resp.Wrap("Required data is missing", nil, err.Error())
		eW.LogError(err, "json.Unmarshal")
		return
	}

	note.ID = hs.st.NextIndex()
	idx, err := hs.st.Add(note)

	if err != nil {
		resp.Wrap("Error in saving record", nil, err.Error())
		eW.LogError(err, "hs.db.RecordSave(note)")
		return
	}

	idxMap := map[string]interface{}{
		"id": idx,
	}

	idxJson, err := json.Marshal(idxMap)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Marshal(idx)")
		return
	}

	resp.Wrap("Successfully added", idxJson, "")
}

func (hs *HttpServer) notesGetHandler(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	eW, err := pkg.NewEWrapperWithFile("(hs *HttpServer) notesGetHandler()")
	if err != nil {
		log.Println("(hs *HttpServer) notesCreateHandler: NewEWrapperWithFile()", err)
	}
	resp := &dto.Response{}
	defer responseReturn(w, eW, resp)

	note := dto.NewNote()
	byteReq, err := io.ReadAll(req.Body)
	if err != nil {
		resp.Wrap("Error reading request", nil, err.Error())
		eW.LogError(err, "io.ReadAll(req.Body)")
		return
	}
	err = json.Unmarshal(byteReq, &note)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Unmarshal(req)")
		return
	}

	if note.ID == -1 {
		err = errors.New("no ID provided")
		resp.Wrap("No ID provided", nil, err.Error())
		eW.LogError(err, "No ID provided")
		return
	}

	records, status := hs.st.GetByIndex(note.ID)
	if !status {
		resp.Wrap("Error in finding notes", nil, errors.New("no notes found").Error())
		eW.LogError(err, "hs.db.RecordsGet(notes)")
		return
	}

	recordsJSON, err := json.Marshal(records)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Marshal(records)")
		return
	}

	resp.Wrap("Success", recordsJSON, "")
}

func (hs *HttpServer) notesUpdateHandler(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	eW, err := pkg.NewEWrapperWithFile("(hs *HttpServer) notesUpdateHandler()")
	if err != nil {
		log.Println("(hs *HttpServer) notesCreateHandler: NewEWrapperWithFile()", err)
	}

	resp := &dto.Response{}
	defer responseReturn(w, eW, resp)

	note := dto.NewNote()
	byteReq, err := io.ReadAll(req.Body)
	if err != nil {
		resp.Wrap("Error reading request", nil, err.Error())
		eW.LogError(err, "io.ReadAll(req.Body)")
		return
	}

	err = json.Unmarshal(byteReq, &note)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Unmarshal(byteReq, &note)")
		return
	}

	if (note.Name == "" && note.LastName == "" && note.Note == "") || note.ID < 1 {
		err = errors.New("required data is missing or invalid")
		resp.Wrap("Required data is missing", nil, err.Error())
		eW.LogError(err, "json.Unmarshal")
		return
	}

	hs.st.RemoveByIndex(note.ID)
	err = hs.st.AddToIndex(note, note.ID)
	if err != nil {
		resp.Wrap("Error in updating note", nil, err.Error())
		eW.LogError(err, "hs.db.RecordUpdate(note)")
		return
	}
	resp.Wrap("Success", nil, "")
}

func (hs *HttpServer) notesDeleteByPhone(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	eW, err := pkg.NewEWrapperWithFile("(hs *HttpServer) notesDeleteByPhone()")
	if err != nil {
		log.Println("(hs *HttpServer) notesCreateHandler: NewEWrapperWithFile()", err)
	}

	resp := &dto.Response{}
	defer responseReturn(w, eW, resp)

	note := dto.NewNote()
	byteReq, err := io.ReadAll(req.Body)
	if err != nil {
		resp.Wrap("Error reading request", nil, err.Error())
		eW.LogError(err, "io.ReadAll(r.Body)")
		return
	}
	err = json.Unmarshal(byteReq, &note)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Unmarshal(byteReq, &note)")
		return
	}

	if note.ID == -1 {
		err = errors.New("id is missing")
		resp.Wrap("ID is missing", nil, err.Error())
		eW.LogError(err, "json.Unmarshal")
		return
	}

	hs.st.RemoveByIndex(note.ID)
	resp.Wrap("Success", nil, "")
}

func (hs *HttpServer) notesGetAll(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	eW, err := pkg.NewEWrapperWithFile("(hs *HttpServer) notesGetAll()")
	if err != nil {
		log.Println("(hs *HttpServer) notesCreateHandler: NewEWrapperWithFile()", err)
	}

	resp := &dto.Response{}
	defer responseReturn(w, eW, resp)

	notes, status := hs.st.GetAll()
	if !status {
		resp.Wrap("Error in finding notes", nil, errors.New("no notes found").Error())
		eW.LogError(err, "hs.db.RecordsGetAll()")
		return
	}

	recordsJSON, err := json.Marshal(notes)
	if err != nil {
		resp.Wrap("Error JSON", nil, err.Error())
		eW.LogError(err, "json.Marshal(records)")
		return
	}

	resp.Wrap("Success", recordsJSON, "")
}

func responseReturn(w http.ResponseWriter, eW *pkg.EWrapper, resp *dto.Response) {
	errEncode := json.NewEncoder(w).Encode(resp)
	if errEncode != nil {
		eW.LogError(errEncode, "json.NewEncoder(w).Encode(resp)")
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	eW.Close()
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
