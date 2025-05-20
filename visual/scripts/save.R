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

  data_all_analytics$windows_amount <-
    ceiling(data_all_analytics$file_size /
              data_all_analytics$window_length_bytes)

  data_all_analytics$minimum_dictionary_key_length <-
    pmax(ceiling(log2(data_all_analytics$dictionary_length)), 1)

  data_all_analytics$data_size_bytes <-
    ceiling(data_all_analytics$windows_amount *
              data_all_analytics$minimum_dictionary_key_length / 8)

  data_all_analytics$dictionary_size_bytes <-
    ceiling(data_all_analytics$dictionary_length *
              data_all_analytics$window_length_bytes)

  data_all_analytics$compressed_size_bytes <-
    data_all_analytics$data_size_bytes +
    data_all_analytics$dictionary_size_bytes

  data_all_analytics$compression_ratio <-
    data_all_analytics$file_size /
    data_all_analytics$compressed_size_bytes

  data_all_analytics$compressed_size_ratio <-
    data_all_analytics$compressed_size_bytes /
    data_all_analytics$file_size * 100 *
    (1 - data_all_analytics$dictionary_limit_reached) +
    data_all_analytics$dictionary_limit_reached * 100

  write.csv(
            data_all_analytics,
            file = "visual/data/data_all_analytics.csv",
            row.names = TRUE)
} else {
  print("analytics file already exists!")
}

if (!file.exists("visual/data/data_all_cloud.csv")) {
  query_all_cloud <- "SELECT * FROM cloud"
  data_all_cloud <- dbGetQuery(con, query_all_cloud)

  write.csv(
            data_all_cloud,
            file = "visual/data/data_all_cloud.csv",
            row.names = TRUE)
} else {
  print("cloud file already exists!")
}