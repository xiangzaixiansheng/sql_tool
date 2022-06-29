Importing the package

Dumping a query to a file

```go
// we're assuming you've setup your sql.DB etc elsewhere
rows, _ := db.Query("SELECT * FROM users WHERE something=72")

err := sqltocsv.WriteFile("./test.csv", rows)
if err != nil {
    panic(err)
}
```

Return a query as a CSV download on the world wide web

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM users WHERE something=72")
    if err != nil {
        http.Error(w, err, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    w.Header().Set("Content-type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment; filename=\"report.csv\"")

    sqltocsv.Write(w, rows)
})
http.ListenAndServe(":8080", nil)
```

`Write` and `WriteFile` should do cover the common cases by the power of _Sensible Defaults™_ but if you need more flexibility you can get an instance of a `Converter` and fiddle with a few settings.

```go
rows, _ := db.Query("SELECT * FROM users WHERE something=72")

csvConverter := sqltocsv.New(rows)

csvConverter.TimeFormat = time.RFC822
csvConverter.Headers = append(rows.Columns(), "extra_column_one", "extra_column_two")

csvConverter.SetRowPreProcessor(func (columns []string) (bool, []string) {
    // exclude admins from report
    // NOTE: this is a dumb example because "where role != 'admin'" is better
    // but every now and then you need to exclude stuff because of calculations
    // that are a pain in sql and this is a contrived README
    if columns[3] == "admin" {
      return false, []string{}
    }

    extra_column_one = generateSomethingHypotheticalFromColumn(columns[2])
    extra_column_two = lookupSomeApiThingForColumn(columns[4])

    return append(columns, extra_column_one, extra_column_two)
})

csvConverter.WriteFile("~/important_user_report.csv")
```

