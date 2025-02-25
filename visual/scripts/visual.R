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
  print("Database connection successful!")
} else {
  print("Database connection failed.")
  stop("Database connection error")
}

# load the data
query_all <- "SELECT * FROM measurements"
data_all <- dbGetQuery(con, query_all)

# scale from ns to ms
data_all$compress_duration_algorithm <-
  data_all$compress_duration_algorithm / 1000000

data_all$decompress_duration_algorithm <-
  data_all$decompress_duration_algorithm / 1000000

#head(data_all)
#data_all$algorithm

if (!require(ggplot2)) {install.packages("ggplot2")}

library(ggplot2)

algorithms <- unique(data_all$algorithm)

# box plot
matrics <- list(
                c(
                  "compress_duration_algorithm",
                  "compress_duration",
                  "Compression Duration (ms)"),
                c(
                  "decompress_duration_algorithm",
                  "decompress_duration",
                  "Decompression Duration (ms)"),
                c(
                  "compression_ratio",
                  "compression_ratio",
                  "Compression Ratio [Higher equals small resulting file]"))

for (matric in matrics) {
  p <- ggplot(data_all, aes(x = algorithm, y = !!sym(matric[1]))) +
    geom_boxplot() +
    labs(
         title = paste(
                       "Min: ", min(data_all[[matric[1]]]),
                       "Max: ", max(data_all[[matric[1]]])),
         x = "Algorithm",
         y = matric[3])

  ggsave(
         filename = paste("boxplot_all_", matric[2], ".png"),
         plot = p,
         path = "visual/out/",
         create.dir = TRUE)

  for (algorithm in algorithms) {
    data <- data_all[data_all$algorithm == algorithm, ]

    min <- min(data[[matric[1]]])
    max <- max(data[[matric[1]]])
    mean <- mean(data[[matric[1]]])
    median <- median(data[[matric[1]]])

    p <- ggplot(data, aes(x = algorithm, y = !!sym(matric[1]))) +
      geom_boxplot() +
      labs(
           title = paste(
                         "Min: ", min,
                         "Max: ", max,
                         "Mean: ", mean,
                         "Median: ", median),
           x = "Algorithm",
           y = matric[3])

    ggsave(
           filename = paste("boxplot_", algorithm, "_", matric[2], ".png"),
           plot = p,
           path = "visual/out/")
  }
}

# scatter plot
data_scatter <- data.frame(
                           algorithm = character(),
                           compress_duration_algorithm = numeric(),
                           decompress_duration_algorithm = numeric(),
                           compression_ratio = numeric())

for (algorithm in algorithms) {
  data <- data_all[data_all$algorithm == algorithm, ]

  compress_duration <- median(data$compress_duration_algorithm)
  decompress_duration <- median(data$decompress_duration_algorithm)
  compression_ratio <- median(data$compression_ratio)

  data_scatter <- rbind(
                        data_scatter,
                        data.frame(
                                   algorithm = algorithm,
                                   compress_duration_algorithm =
                                     compress_duration,
                                   decompress_duration_algorithm =
                                     decompress_duration,
                                   compression_ratio = compression_ratio))
}

p <- ggplot(
            data_scatter,
            aes(
                x = compress_duration_algorithm,
                y = decompress_duration_algorithm,
                size = compression_ratio)) +
  geom_point(aes(color = algorithm), alpha = 0.7) +
  scale_size_continuous(range = c(10, 25)) +
  labs(
       title = "Compression vs Decompression Duration",
       x = "Compression Duration (ms)",
       y = "Decompression Duration (ms)")

ggsave(
       filename = "scatter_compress_decompress.png",
       plot = p,
       path = "visual/out/")