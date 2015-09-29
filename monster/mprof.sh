go build;
./monster -memprof monster.mprof -bagdir ../bags -count 100000 -o out ../prods/fast1K.prod
go tool pprof --svg --inuse_space monster monster.mprof > monster.inuse.svg
go tool pprof --svg --alloc_space monster monster.mprof > monster.alloc.svg
