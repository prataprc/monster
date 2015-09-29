go build;
./monster -pprof monster.pprof -bagdir ../bags -count 100000 -o out ../prods/fast1K.prod
go tool pprof --svg monster monster.pprof > monster.pprof.svg
