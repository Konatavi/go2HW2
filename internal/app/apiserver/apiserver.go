package apiserver

import (
	"net/http"

	"github.com/Konatavi/go2HW2/internal/app/middleware"
	"github.com/Konatavi/go2HW2/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	prefix string = "/api/v1"
)

// type for APIServer object for instancing server
type APIServer struct {
	//Unexported field
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

//APIServer constructor
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start http server and connection to db and logger confs
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info("starting api server at port :", s.config.BindAddr)
	s.configureRouter()
	if err := s.configureStore(); err != nil {
		return err
	}
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

//func for configureate logger, should be unexported
func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return nil
	}
	s.logger.SetLevel(level)

	return nil
}

//func for configure Router
func (s *APIServer) configureRouter() {

	/// Необходимые роутеры

	// 1) POST /register - позволяет зарегестрировать нового пользователя для API. Завершается кодом
	//201 и сообщением {"Message" : "User created. Try to auth"} в случае, если такого
	//пользователя еще не было в БД. В противном случае завершаемся кодом 400 и сообщением
	//{"Error" : "User already exists"}.
	s.router.HandleFunc(prefix+"/register", s.PostUserRegister).Methods("POST")

	// 2) POST /auth - возвращает JWT метку для зарегестрированных пользователей.
	s.router.HandleFunc(prefix+"/auth", s.PostToAuth).Methods("POST")

	// 3) GET /auto/<string:mark> - возвращает информацию про автомобиль с именем mark и код 200.
	//В случае, если автомобиля нет в БД в текущий момент возвращаем {"Error" : "Auto with
	//that mark not found"} и код 404.
	s.router.Handle(prefix+"/auto"+"/{mark}", middleware.JwtMiddleware.Handler(
		http.HandlerFunc(s.GetAutoByMark),
	)).Methods("GET")

	// 4) POST /auto/<string:mark> - добавляет автомобиль с именем mark в БД. В случае успеха - 201 и
	//сообщение {"Message" : "Auto created"}. В случае, если автомобиль с таким именем уже
	//существует - 400 и {"Error" : "Auto with that mark exists"}.
	s.router.Handle(prefix+"/auto"+"/{mark}", middleware.JwtMiddleware.Handler(
		http.HandlerFunc(s.PostAuto),
	)).Methods("POST")

	// 5) PUT /auto/<string:mark> - обновляет информацию про автомобиль с именем mark в БД. В
	//случае успеха - 202 и сообщение {"Message" : "Auto updated"}. В случае, если автомобиля нет
	//в БД в текущий момент возвращаем {"Error" : "Auto with that mark not found"} и код 404.
	s.router.Handle(prefix+"/auto"+"/{mark}", middleware.JwtMiddleware.Handler(
		http.HandlerFunc(s.PutAuto),
	)).Methods("PUT")

	// 6) DELETE /auto/<string:mark> - удаляет информацию про автомобиль с именем mark из БД. В
	//случае успеха - 202 и сообщение {"Message" : "Auto deleted"}. В случае, если автомобиля нет
	//в БД в текущий момент возвращаем {"Error" : "Auto with that mark not found"} и код 404.
	s.router.Handle(prefix+"/auto"+"/{mark}", middleware.JwtMiddleware.Handler(
		http.HandlerFunc(s.DeleteAuto),
	)).Methods("DELETE")

	// 7) GET /stock - возвращает информацию про все имеющиеся на данный момент в БД автомобили
	// и код 200 в случае, если имеется хотя бы один автомобиль в наличии. В противном случае - 400 и
	// сообщение {"Error" : "No one autos found in DataBase"}.
	s.router.Handle(prefix+"/stock", middleware.JwtMiddleware.Handler(
		http.HandlerFunc(s.GetAllAutos),
	)).Methods("GET")

}

//configureStore method
func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}
