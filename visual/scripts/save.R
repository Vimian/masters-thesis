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

if (!file.exists("visual/data/data_all_analytics.csv")) {
  query_all_analytics <- "SELECT * FROM analytics"
  data_all_analytics <- dbGetQuery(con, query_all_analytics)

  write.csv(
            data_all_analytics,
            file = "visual/data/data_all_analytics.csv",
            row.names = TRUE)
} else {
  print("analytics file already exists!")
}