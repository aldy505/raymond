package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Deps struct {
	DB *sql.DB
}

func main() {
	log.Println("Server is starting up")

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "80"
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "0.0.0.0"
	}

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dbUrl = "./db.sqlite"
	}

	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	deps := &Deps{DB: db}

	log.Println("Migrating database in progress")

	prepareCtx, prepareCancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer prepareCancel()

	err = deps.Migrate(prepareCtx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Migrating database completed")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/list", deps.List)
	mux.HandleFunc("/api/add", deps.Add)
	mux.HandleFunc("/", deps.Index)

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)

	go func() {
		log.Printf("Server running on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-sig

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}
}

func (d *Deps) Index(w http.ResponseWriter, r *http.Request) {
	sakuraCss := `/* Sakura.css v1.3.1
	* ================
	* Minimal css theme.
	* Project: https://github.com/oxalorg/sakura/
	*/
   /* Body */
   html {
	 font-size: 62.5%;
	 font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif; }
   
   body {
	 font-size: 1.8rem;
	 line-height: 1.618;
	 max-width: 38em;
	 margin: auto;
	 color: #4a4a4a;
	 background-color: #f9f9f9;
	 padding: 13px; }
   
   @media (max-width: 684px) {
	 body {
	   font-size: 1.53rem; } }
   
   @media (max-width: 382px) {
	 body {
	   font-size: 1.35rem; } }
   
   h1, h2, h3, h4, h5, h6 {
	 line-height: 1.1;
	 font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif;
	 font-weight: 700;
	 margin-top: 3rem;
	 margin-bottom: 1.5rem;
	 overflow-wrap: break-word;
	 word-wrap: break-word;
	 -ms-word-break: break-all;
	 word-break: break-word; }
   
   h1 {
	 font-size: 2.35em; }
   
   h2 {
	 font-size: 2.00em; }
   
   h3 {
	 font-size: 1.75em; }
   
   h4 {
	 font-size: 1.5em; }
   
   h5 {
	 font-size: 1.25em; }
   
   h6 {
	 font-size: 1em; }
   
   p {
	 margin-top: 0px;
	 margin-bottom: 2.5rem; }
   
   small, sub, sup {
	 font-size: 75%; }
   
   hr {
	 border-color: #1d7484; }
   
   a {
	 text-decoration: none;
	 color: #1d7484; }
	 a:hover {
	   color: #982c61;
	   border-bottom: 2px solid #4a4a4a; }
	 a:visited {
	   color: #144f5a; }
   
   ul {
	 padding-left: 1.4em;
	 margin-top: 0px;
	 margin-bottom: 2.5rem; }
   
   li {
	 margin-bottom: 0.4em; }
   
   blockquote {
	 margin-left: 0px;
	 margin-right: 0px;
	 padding-left: 1em;
	 padding-top: 0.8em;
	 padding-bottom: 0.8em;
	 padding-right: 0.8em;
	 border-left: 5px solid #1d7484;
	 margin-bottom: 2.5rem;
	 background-color: #f1f1f1; }
   
   blockquote p {
	 margin-bottom: 0; }
   
   img, video {
	 height: auto;
	 max-width: 100%;
	 margin-top: 0px;
	 margin-bottom: 2.5rem; }
   
   /* Pre and Code */
   pre {
	 background-color: #f1f1f1;
	 display: block;
	 padding: 1em;
	 overflow-x: auto;
	 margin-top: 0px;
	 margin-bottom: 2.5rem;
	 font-size: 0.9em; }
   
   code, kbd, samp {
	 font-size: 0.9em;
	 padding: 0 0.5em;
	 background-color: #f1f1f1;
	 white-space: pre-wrap; }
   
   pre > code {
	 padding: 0;
	 background-color: transparent;
	 white-space: pre;
	 font-size: 1em; }
   
   /* Tables */
   table {
	 text-align: justify;
	 width: 100%;
	 border-collapse: collapse; }
   
   td, th {
	 padding: 0.5em;
	 border-bottom: 1px solid #f1f1f1; }
   
   /* Buttons, forms and input */
   input, textarea {
	 border: 1px solid #4a4a4a; }
	 input:focus, textarea:focus {
	   border: 1px solid #1d7484; }
   
   textarea {
	 width: 100%; }
   
   .button, button, input[type="submit"], input[type="reset"], input[type="button"] {
	 display: inline-block;
	 padding: 5px 10px;
	 text-align: center;
	 text-decoration: none;
	 white-space: nowrap;
	 background-color: #1d7484;
	 color: #f9f9f9;
	 border-radius: 1px;
	 border: 1px solid #1d7484;
	 cursor: pointer;
	 box-sizing: border-box; }
	 .button[disabled], button[disabled], input[type="submit"][disabled], input[type="reset"][disabled], input[type="button"][disabled] {
	   cursor: default;
	   opacity: .5; }
	 .button:focus:enabled, .button:hover:enabled, button:focus:enabled, button:hover:enabled, input[type="submit"]:focus:enabled, input[type="submit"]:hover:enabled, input[type="reset"]:focus:enabled, input[type="reset"]:hover:enabled, input[type="button"]:focus:enabled, input[type="button"]:hover:enabled {
	   background-color: #982c61;
	   border-color: #982c61;
	   color: #f9f9f9;
	   outline: 0; }
   
   textarea, select, input {
	 color: #4a4a4a;
	 padding: 6px 10px;
	 /* The 6px vertically centers text on FF, ignored by Webkit */
	 margin-bottom: 10px;
	 background-color: #f1f1f1;
	 border: 1px solid #f1f1f1;
	 border-radius: 4px;
	 box-shadow: none;
	 box-sizing: border-box; }
	 textarea:focus, select:focus, input:focus {
	   border: 1px solid #1d7484;
	   outline: 0; }
   
   input[type="checkbox"]:focus {
	 outline: 1px dotted #1d7484; }
   
   label, legend, fieldset {
	 display: block;
	 margin-bottom: .5rem;
	 font-weight: 600; }`

	htmlResponse := `
	<!DOCTYPE html>
	<html>
	<head>
	<title>How many times Raymond said sorry so far</title>
	<style>` + sakuraCss + `</style>
	<style>
		.pointer:hover {
			cursor: pointer;
		}
	</style>
	<script>
	async function listCounter() {
		const response = await fetch("/api/list", { method: "GET" });
		const respBody = await response.json();

		const counterElement = document.getElementById("counter-content");
		counterElement.innerHTML = respBody.counter;

		const lastTimeElement = document.getElementById("lasttime-content");
		if (new Date(respBody.lastDate).getUTCFullYear() == 1970) {
			lastTimeElement.innerHTML = "never";
		} else {
			lastTimeElement.innerHTML = new Date(respBody.lastDate).toLocaleString("id-ID");
		};
	};
	
	async function addCounter() {
		const response = await fetch("/api/add", { method: "POST" });
		
		await listCounter();
	};

	setInterval(async () => {
		await listCounter();
	}, 5000);
	</script>
	</head>
	<body>
	<h4 style="margin-top: 3rem; text-align: center;">
		How many times Raymond said sorry, so far
	</h4>

	<h1 style="font-size: 8rem; margin-top: 2rem; text-align: center; margin-left: auto; margin-right: auto;">
	  <span id="counter-content">0</span>
	</h1>

	<p style="text-align: center;">Last time he said it, it was at <span id="lasttime-content">never</span></p>
	<div onclick="addCounter()" class="pointer">
		<h3 style="margin-top: 0.5rem; text-align: center;">He said it again!</h3>
	</div>
	</body>
	</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlResponse))
}

func (d *Deps) Migrate(ctx context.Context) error {
	c, err := d.DB.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Println(err)
		}
	}()

	tx, err := c.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS counter (
			count INTEGER NOT NULL,
			created_at DATETIME NOT NULL
		)`,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS counter_aggregate (
			counts INTEGER NOT NULL,
			created_at DATETIME NOT NULL
		)`,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return e
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Deps) Add(w http.ResponseWriter, r *http.Request) {
	conn, err := d.DB.Conn(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}
	defer func() {
		if err := conn.Close(); err != nil && !errors.Is(err, sql.ErrConnDone) {
			log.Println(err)
		}
	}()

	tx, err := conn.BeginTx(r.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}

	_, err = tx.ExecContext(
		r.Context(),
		`INSERT INTO counter (count, created_at) VALUES (?, ?)`,
		1,
		time.Now(),
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}

	if err := tx.Commit(); err != nil {
		if e := tx.Rollback(); e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}

	go d.CreateAggregate()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"success"}`))
}

func (d *Deps) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*15)
	defer cancel()

	c, err := d.DB.Conn(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Println(err)
		}
	}()

	var counts int
	var lastDate time.Time
	err = c.QueryRowContext(
		ctx,
		`SELECT counts, created_at FROM counter_aggregate ORDER BY created_at DESC LIMIT 1`,
	).Scan(
		&counts,
		&lastDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responseBody, err := json.Marshal(map[string]interface{}{
				"counter":  0,
				"lastDate": time.Unix(0, 0).Format(time.RFC3339),
			})
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}

	responseBody, err := json.Marshal(map[string]interface{}{
		"counter":  counts,
		"lastDate": lastDate.Format(time.RFC3339),
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func (d *Deps) CreateAggregate() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	c, err := d.DB.Conn(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	tx, err := c.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		log.Println(err)
		return
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT count FROM counter`,
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var counts int
	for rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			log.Println(err)
			return
		}

		counts += count
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO
			counter_aggregate
			(counts, created_at)
			VALUES
			(?, ?)`,
		counts,
		time.Now(),
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			log.Println(err)
			return
		}

		log.Println(err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Aggregate created, with counts: %d", counts)
}
