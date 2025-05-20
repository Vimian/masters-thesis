if (!require(DBI)) {install.packages("DBI")}
if (!require(RPostgres)) {install.packages("RPostgres")}

library(DBI)
library(RPostgres)

# TODO: change to dinamic values
con <- dbConnect(RPostgres::Postgres(),
                 host = "postgres",
                 port = 5432,
                 dbname = "benchmark",
                 user = "postgres",
                 password = "password")

if (dbIsValid(con)) {
  print("database connection successful!")
} else {
  print("database connection failed.")
  stop("database connection error")
}

# if csv file already exists do not overwrite it
if (!file.exists("visual/data/data_all_measurements.csv")) {
  query_all_measurements <- "SELECT * FROM measurements"
  data_all_measurements <- dbGetQuery(con, query_all_measurements)

  write.csv(
            data_all_measurements,
            file = "visual/data/data_all_measurements.csv",
            row.names = TRUE)
} else {
  print("benchmark file already exists!")
}