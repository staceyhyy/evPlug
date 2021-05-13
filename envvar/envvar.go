package envvar

import (
	"golive/logger"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Address        string
	Lat            string
	Lng            string
	DBHost         string
	DBUser         string
	DBPwd          string
	DBName         string
	DBPointColl    string
	DBUserColl     string
	DBVehicleColl  string
	DBChargerColl  string
	DBProviderColl string
	APIKey         string
	LogFile        string
}

//Load all parameter from .env file
func Load() (Env, error) {
	err := godotenv.Load(".env")
	if err != nil {

		logger.Fatal.Println("error in loading env file: ", err)
		return Env{}, err
	}
	var envVar = Env{}
	envVar.Address = HostAddress()
	envVar.Lat = os.Getenv("LAT")
	envVar.Lng = os.Getenv("LNG")
	envVar.DBHost = os.Getenv("DB_HOST")
	envVar.DBUser = os.Getenv("DB_USER")
	envVar.DBPwd = os.Getenv("DB_PWD")
	envVar.DBName = os.Getenv("DB_NAME")
	envVar.DBPointColl = os.Getenv("DB_POINTS_COLL")
	envVar.DBUserColl = os.Getenv("DB_USERS_COLL")
	envVar.DBVehicleColl = os.Getenv("DB_VEHICLES_COLL")
	envVar.DBChargerColl = os.Getenv("DB_CHARGERS_COLL")
	envVar.DBProviderColl = os.Getenv("DB_PROVIDERS_COLL")
	envVar.APIKey = os.Getenv("GOOGLE_API_KEY")
	envVar.LogFile = os.Getenv("LOGFILE")
	return envVar, nil
}

func HostAddress() string {
	return os.Getenv("HOST") + ":" + os.Getenv("PORT")
}

func CookieName() string {
	return os.Getenv("COOKIE_NAME")
}
